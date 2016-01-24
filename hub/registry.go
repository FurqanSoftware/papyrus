package hub

import (
	"encoding/json"

	"github.com/desertbit/glue"
)

type Registry struct {
	socks map[string]map[*glue.Socket]bool
	docs  map[*glue.Socket]*Document
}

func NewRegistry() *Registry {
	return &Registry{
		socks: map[string]map[*glue.Socket]bool{},
		docs:  map[*glue.Socket]*Document{},
	}
}

func (r *Registry) attach(sock *glue.Socket, docID string) error {
	doc, err := DefaultRepository.Get(docID)
	if err != nil {
		return err
	}

	if r.socks[docID] == nil {
		r.socks[docID] = map[*glue.Socket]bool{}
	}
	r.socks[docID][sock] = true

	if r.docs[sock] != nil {
		delete(r.socks[docID], sock)
	}
	r.docs[sock] = doc

	return nil
}

func (r *Registry) detach(sock *glue.Socket) {
	doc, ok := r.docs[sock]
	if !ok {
		return
	}
	_, ok = r.socks[doc.ID]
	if !ok {
		return
	}

	delete(r.socks[doc.ID], sock)
	delete(r.docs, sock)

	if len(r.socks[doc.ID]) == 0 {
		delete(r.socks, doc.ID)
	}
}

func (r *Registry) broadcast(docID string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	for sock := range r.socks[docID] {
		sock.Write("change " + string(b))
	}
	return nil
}

func (r *Registry) document(sock *glue.Socket) *Document {
	return r.docs[sock]
}

var registry = NewRegistry()
