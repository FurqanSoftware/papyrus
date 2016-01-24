package main

import (
	"github.com/gophergala2016/papyrus/hub"
	"github.com/gophergala2016/papyrus/ot"
)

type Repo map[string][]byte

func (r Repo) Get(id string) (*hub.Doc, error) {
	b, ok := r[id]
	if !ok {
		return nil, nil
	}
	return &hub.Doc{
		ID:   id,
		Blob: ot.Blob(b),
	}, nil
}

func (r Repo) Put(id string, doc *hub.Doc) error {
	r[id] = []byte(doc.Blob)
	return nil
}
