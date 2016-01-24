package repo

import (
	"errors"

	"gopkg.in/mgo.v2/bson"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gophergala2016/papyrus/hub"
	"github.com/gophergala2016/papyrus/ot"
)

var ErrBadID = errors.New("repo: invalid document ID")

type Repo struct {
	docs   map[string]*hub.Doc
	syncCh chan *hub.Doc
	nextCh chan string
}

func New() *Repo {
	repo := &Repo{
		docs:   map[string]*hub.Doc{},
		syncCh: make(chan *hub.Doc),
		nextCh: make(chan string),
	}
	go repo.loop()
	return repo
}

func (r *Repo) Get(id string) (*hub.Doc, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, ErrBadID
	}
	doc, err := data.GetDocument(bson.ObjectIdHex(id))
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}
	return &hub.Doc{
		ID:   doc.ID.Hex(),
		Blob: ot.Blob(doc.Content),
	}, nil
}

func (r *Repo) Put(id string, doc *hub.Doc) error {
	go func() {
		r.syncCh <- doc
	}()
	return nil
}

func (r *Repo) loop() {
	for {
		select {
		case doc := <-r.syncCh:
			r.docs[doc.ID] = doc
			go func() {
				r.nextCh <- doc.ID
			}()

		case id := <-r.nextCh:
			doc := r.docs[id]
			delete(r.docs, doc.ID)

			doc2, err := data.GetDocument(bson.ObjectIdHex(id))
			if err != nil {
				panic(err)
			}
			doc2.Content = string(doc.Blob)
			err = doc2.Put()
			if err != nil {
				panic(err)
			}
		}
	}
}
