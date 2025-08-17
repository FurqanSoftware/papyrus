package hub

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/desertbit/glue"
	"github.com/dgrijalva/jwt-go"
	"github.com/gophergala2016/papyrus/auth"
	"github.com/gophergala2016/papyrus/ot"
)

type Hub struct {
	attachInCh  chan AttachIn
	attachOutCh chan AttachOut
	changeInCh  chan ChangeIn
	changeOutCh chan ChangeOut
	errorOutCh  chan ErrorOut

	nexus *nexus
}

func New(repo Repo) *Hub {
	hub := &Hub{
		attachInCh:  make(chan AttachIn),
		attachOutCh: make(chan AttachOut),
		changeInCh:  make(chan ChangeIn),
		changeOutCh: make(chan ChangeOut),
		errorOutCh:  make(chan ErrorOut),
		nexus:       newNexus(repo),
	}
	go hub.processAttachIn()
	go hub.processAttachOut()
	go hub.processChangeIn()
	go hub.processChangeOut()
	go hub.processErrorOut()
	return hub
}

func (h *Hub) processAttachIn() {
	for v := range h.attachInCh {
		claims := auth.Claims{}
		_, err := jwt.ParseWithClaims(v.Token, &claims, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			h.sendError(v.Socket, "invalid token")
			continue
		}

		err = h.nexus.attach(v.Socket, claims.DocumentID)
		if err != nil {
			h.sendError(v.Socket, "internal server error")
			continue
		}

		_, ok := h.nexus.sockDoc[v.Socket]
		if !ok {
			h.sendError(v.Socket, "not found")
			continue
		}

		h.attachOutCh <- AttachOut{
			Socket: v.Socket,
		}
	}
}

func (h *Hub) processAttachOut() {
	for v := range h.attachOutCh {
		doc, ok := h.nexus.sockDoc[v.Socket]
		if !ok {
			h.sendError(v.Socket, "not attached")
			continue
		}

		b, err := json.Marshal(NewChangeData(doc.Head()))
		if err != nil {
			h.sendError(v.Socket, "internal server error")
			continue
		}
		v.Socket.Write("change " + string(b))
	}
}

func (h *Hub) processChangeIn() {
	for v := range h.changeInCh {
		doc, ok := h.nexus.sockDoc[v.Socket]
		if !ok {
			h.sendError(v.Socket, "not found")
			continue
		}

		ch, err := doc.Apply(v.Change)
		if err != nil {
			h.sendError(v.Socket, "invalid change")
			continue
		}

		h.changeOutCh <- ChangeOut{
			Socket: v.Socket,
			Change: ch,
		}
	}
}

func (h *Hub) processChangeOut() {
	for v := range h.changeOutCh {
		doc, ok := h.nexus.sockDoc[v.Socket]
		if !ok {
			h.sendError(v.Socket, "not attached")
			continue
		}

		h.nexus.broadcast(doc.ID, NewChangeData(v.Change))
	}
}

func (h *Hub) processErrorOut() {
	for v := range h.errorOutCh {
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
			h.attachInCh <- AttachIn{
				Socket: sock,
				Token:  fields[1],
			}

		case "change":
			data := ChangeData{}
			err := json.Unmarshal([]byte(fields[1]), &data)
			if err != nil {
				h.sendError(sock, "bad request")
				return
			}

			h.changeInCh <- ChangeIn{
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

	Token string
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
