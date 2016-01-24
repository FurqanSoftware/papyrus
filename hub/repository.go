package hub

type Repository struct {
	Documents map[string]*Document

	GetBlob func(string) ([]byte, error)
	PutBlob func(string, []byte) error
}

func NewRepository() *Repository {
	return &Repository{
		Documents: map[string]*Document{},
	}
}

func (r *Repository) Get(id string) (*Document, error) {
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
			Blob: blob,
		}
		r.Documents[id] = doc
	}
	return doc, nil
}

var DefaultRepository = NewRepository()