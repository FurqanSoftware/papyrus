package hub

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/desertbit/glue"
	"github.com/gophergala2016/papyrus/ot"
)

type ChangeData struct {
	ID   string        `json:"id"`
	Root int           `json:"root"`
	Ops  []interface{} `json:"ops"`
}

func HandleSocket(sock *glue.Socket) {
	sock.OnClose(func() {
		registry.deregisterAll(sock)
	})

	var doc *Document

	sock.OnRead(func(data string) {
		fields := strings.SplitN(data, " ", 2)
		if len(fields) != 2 {
			return
		}
		switch fields[0] {
		case "subscribe":
			var err error
			doc, err = DefaultRepository.Get(fields[1])
			if err != nil {
				sock.Write("error \"internal server error\"")
				return
			}
			if doc == nil {
				sock.Write("error \"not found\"")
				return
			}

			registry.register(sock, "document:"+doc.ID)

			data := ChangeData{
				ID:   "",
				Root: len(doc.History),
				Ops:  []interface{}{string(doc.Blob)},
			}
			dataB, err := json.Marshal(data)
			if err != nil {
				sock.Write("error \"internal server error\"")
				return
			}
			sock.Write("change " + string(dataB))

		case "change":
			data := ChangeData{}
			err := json.Unmarshal([]byte(fields[1]), &data)
			if err != nil {
				sock.Write("error \"bad request\"")
				return
			}

			ch := Change{}
			ch.ID = data.ID
			ch.Root = data.Root
			for _, u := range data.Ops {
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

			ch, err = doc.Apply(ch)
			if err != nil {
				log.Print(err)
				sock.Write("error \"internal server error\"")
				return
			}

			data = ChangeData{
				ID:   data.ID,
				Root: ch.Root,
			}
			for _, u := range ch.Ops {
				switch u := u.(type) {
				case ot.RetainOp:
					data.Ops = append(data.Ops, int(u))

				case ot.InsertOp:
					data.Ops = append(data.Ops, string(u))

				case ot.DeleteOp:
					data.Ops = append(data.Ops, -int(u))
				}
			}

			dataB, err := json.Marshal(data)
			if err != nil {
				sock.Write("error \"internal server error\"")
			}

			registry.deliver("document:"+doc.ID, "change "+string(dataB))
		}
	})
}
