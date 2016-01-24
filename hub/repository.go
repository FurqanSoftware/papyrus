package hub

import "sync"

type Repository struct {
	Documents map[string]*Document

	GetBlob func(string) ([]byte, error)
	PutBlob func(string, []byte) error

	mutex sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		Documents: map[string]*Document{},
	}
}

func (r *Repository) Get(id string) (*Document, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	doc, ok := r.Documents[id]
	if !ok {
		blob, err := r.GetBlob(id)
		if err != nil {
			return nil, err
		}
		if blob == nil {
			return nil, nil
		}
		doc = &Document{
			ID:   id,
			Blob: blob,
		}
		r.Documents[id] = doc
	}
	return doc, nil
}

var DefaultRepository = NewRepository()
