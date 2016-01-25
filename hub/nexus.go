package hub

import (
	"encoding/json"

	"github.com/desertbit/glue"
)

type nexus struct {
	repo    Repo
	docs    map[string]*Doc
	socks   map[string]map[*glue.Socket]bool
	sockDoc map[*glue.Socket]*Doc
}

func newNexus(repo Repo) *nexus {
	return &nexus{
		repo:    repo,
		docs:    map[string]*Doc{},
		socks:   map[string]map[*glue.Socket]bool{},
		sockDoc: map[*glue.Socket]*Doc{},
	}
}

func (r *nexus) attach(sock *glue.Socket, docID string) error {
	_, ok := r.docs[docID]
	if !ok {
		doc, err := r.repo.Get(docID)
		if err != nil {
			return err
		}
		if doc == nil {
			return nil
		}
		r.docs[docID] = doc
	}
	doc := r.docs[docID]

	if r.socks[docID] == nil {
		r.socks[docID] = map[*glue.Socket]bool{}
	}
	r.socks[docID][sock] = true

	if r.sockDoc[sock] != nil {
		delete(r.socks[docID], sock)
	}
	r.sockDoc[sock] = doc

	return nil
}

func (r *nexus) detach(sock *glue.Socket) {
	doc, ok := r.sockDoc[sock]
	if !ok {
		return
	}
	_, ok = r.socks[doc.ID]
	if !ok {
		return
	}

	delete(r.socks[doc.ID], sock)
	delete(r.sockDoc, sock)

	if len(r.socks[doc.ID]) == 0 {
		delete(r.socks, doc.ID)
	}
}

func (r *nexus) broadcast(docID string, data interface{}) error {
	doc := r.docs[docID]
	err := r.repo.Put(docID, doc)
	if err != nil {
		return err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	for sock := range r.socks[docID] {
		sock.Write("change " + string(b))
	}
	return nil
}
