package hub

type Repo interface {
	Get(string) (*Doc, error)
	Put(string, *Doc) error
}
