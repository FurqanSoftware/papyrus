package hub

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/desertbit/glue"
	"github.com/gophergala2016/papyrus/ot"
)

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

var (
	AttachInChan  = make(chan AttachIn)
	AttachOutChan = make(chan AttachOut)
	ChangeInChan  = make(chan ChangeIn)
	ChangeOutChan = make(chan ChangeOut)
	ErrorOutChan  = make(chan ErrorOut)
)

func processAttachIn() {
	for v := range AttachInChan {
		err := registry.attach(v.Socket, v.DocumentID)
		if err != nil {
			sendError(v.Socket, "internal server error")
			return
		}

		doc := registry.document(v.Socket)
		if doc == nil {
			sendError(v.Socket, "not found")
			return
		}

		AttachOutChan <- AttachOut{
			Socket: v.Socket,
		}
	}
}

func processAttachOut() {
	for v := range AttachOutChan {
		doc := registry.document(v.Socket)
		if doc == nil {
			sendError(v.Socket, "not attached")
			return
		}

		b, err := json.Marshal(NewChangeData(doc.Head()))
		if err != nil {
			sendError(v.Socket, "internal server error")
			return
		}
		v.Socket.Write("change " + string(b))
	}
}

func processChangeIn() {
	for v := range ChangeInChan {
		doc := registry.document(v.Socket)
		if doc == nil {
			sendError(v.Socket, "not found")
			return
		}

		ch, err := doc.Apply(v.Change)
		if err != nil {
			sendError(v.Socket, "invalid change")
			return
		}

		ChangeOutChan <- ChangeOut{
			Socket: v.Socket,
			Change: ch,
		}
	}
}

func processChangeOut() {
	for v := range ChangeOutChan {
		doc := registry.document(v.Socket)
		if doc == nil {
			sendError(v.Socket, "not attached")
			return
		}

		registry.broadcast(doc.ID, NewChangeData(v.Change))
	}
}

func processErrorOut() {
	for v := range ErrorOutChan {
		sendError(v.Socket, v.Message)
	}
}

func sendError(sock *glue.Socket, msg string) {
	sock.Write("error " + strconv.Quote(msg))
}

func HandleSocket(sock *glue.Socket) {
	sock.OnClose(func() {
		registry.detach(sock)
	})

	sock.OnRead(func(data string) {
		fields := strings.SplitN(data, " ", 2)
		if len(fields) != 2 {
			return
		}
		switch fields[0] {
		case "attach":
			AttachInChan <- AttachIn{
				Socket:     sock,
				DocumentID: fields[1],
			}

		case "change":
			data := ChangeData{}
			err := json.Unmarshal([]byte(fields[1]), &data)
			if err != nil {
				sendError(sock, "bad request")
				return
			}

			ChangeInChan <- ChangeIn{
				Socket: sock,
				Change: data.Change(),
			}
		}
	})
}

func init() {
	go processAttachIn()
	go processAttachOut()
	go processChangeIn()
	go processChangeOut()
	go processErrorOut()
}
