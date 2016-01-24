package hub

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/desertbit/glue"
	"github.com/gophergala2016/papyrus/ot"
)

type Hub struct {
	attachInChan  chan AttachIn
	attachOutChan chan AttachOut
	changeInChan  chan ChangeIn
	changeOutChan chan ChangeOut
	errorOutChan  chan ErrorOut

	nexus *nexus
}

func New(repo Repo) *Hub {
	hub := &Hub{
		attachInChan:  make(chan AttachIn),
		attachOutChan: make(chan AttachOut),
		changeInChan:  make(chan ChangeIn),
		changeOutChan: make(chan ChangeOut),
		errorOutChan:  make(chan ErrorOut),
		nexus:         newNexus(repo),
	}
	go hub.processAttachIn()
	go hub.processAttachOut()
	go hub.processChangeIn()
	go hub.processChangeOut()
	go hub.processErrorOut()
	return hub
}

func (h *Hub) processAttachIn() {
	for v := range h.attachInChan {
		err := h.nexus.attach(v.Socket, v.DocumentID)
		if err != nil {
			h.sendError(v.Socket, "internal server error")
			return
		}

		_, ok := h.nexus.sockDoc[v.Socket]
		if !ok {
			h.sendError(v.Socket, "not found")
			return
		}

		h.attachOutChan <- AttachOut{
			Socket: v.Socket,
		}
	}
}

func (h *Hub) processAttachOut() {
	for v := range h.attachOutChan {
		doc, ok := h.nexus.sockDoc[v.Socket]
		if !ok {
			h.sendError(v.Socket, "not attached")
			return
		}

		b, err := json.Marshal(NewChangeData(doc.Head()))
		if err != nil {
			h.sendError(v.Socket, "internal server error")
			return
		}
		v.Socket.Write("change " + string(b))
	}
}

func (h *Hub) processChangeIn() {
	for v := range h.changeInChan {
		doc, ok := h.nexus.sockDoc[v.Socket]
		if !ok {
			h.sendError(v.Socket, "not found")
			return
		}

		ch, err := doc.Apply(v.Change)
		if err != nil {
			h.sendError(v.Socket, "invalid change")
			return
		}

		h.changeOutChan <- ChangeOut{
			Socket: v.Socket,
			Change: ch,
		}
	}
}

func (h *Hub) processChangeOut() {
	for v := range h.changeOutChan {
		doc, ok := h.nexus.sockDoc[v.Socket]
		if !ok {
			h.sendError(v.Socket, "not attached")
			return
		}

		h.nexus.broadcast(doc.ID, NewChangeData(v.Change))
	}
}

func (h *Hub) processErrorOut() {
	for v := range h.errorOutChan {
		h.sendError(v.Socket, v.Message)
	}
}

func (h *Hub) sendError(sock *glue.Socket, msg string) {
	sock.Write("error " + strconv.Quote(msg))
}

func (h *Hub) HandleSocket(sock *glue.Socket) {
	sock.OnClose(func() {
		h.nexus.detach(sock)
	})

	sock.OnRead(func(data string) {
		fields := strings.SplitN(data, " ", 2)
		if len(fields) != 2 {
			return
		}
		switch fields[0] {
		case "attach":
			h.attachInChan <- AttachIn{
				Socket:     sock,
				DocumentID: fields[1],
			}

		case "change":
			data := ChangeData{}
			err := json.Unmarshal([]byte(fields[1]), &data)
			if err != nil {
				h.sendError(sock, "bad request")
				return
			}

			h.changeInChan <- ChangeIn{
				Socket: sock,
				Change: data.Change(),
			}
		}
	})
}

type ChangeData struct {
	ID   string        `json:"id"`
	Root int           `json:"root"`
	Ops  []interface{} `json:"ops"`
}

func NewChangeData(ch Change) ChangeData {
	data := ChangeData{
		ID:   ch.ID,
		Root: ch.Root,
		Ops:  []interface{}{},
	}
	for _, o := range ch.Ops {
		switch o := o.(type) {
		case ot.RetainOp:
			data.Ops = append(data.Ops, int(o))

		case ot.InsertOp:
			data.Ops = append(data.Ops, string(o))

		case ot.DeleteOp:
			data.Ops = append(data.Ops, -int(o))
		}
	}
	return data
}

func (d ChangeData) Change() Change {
	ch := Change{
		ID:   d.ID,
		Root: d.Root,
	}
	for _, u := range d.Ops {
		switch u := u.(type) {
		case string:
			ch.Ops = append(ch.Ops, ot.InsertOp(u))

		case float64:
			if u >= 0 {
				ch.Ops = append(ch.Ops, ot.RetainOp(u))
			} else {
				ch.Ops = append(ch.Ops, ot.DeleteOp(-u))
			}
		}
	}
	return ch
}

type AttachIn struct {
	Socket *glue.Socket

	DocumentID string
}

type AttachOut struct {
	Socket *glue.Socket
}

type ChangeIn struct {
	Socket *glue.Socket

	Change Change
}

type ChangeOut struct {
	Socket *glue.Socket

	Change Change
}

type ErrorOut struct {
	Socket *glue.Socket

	Error   error
	Message string
}
