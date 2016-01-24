package hub

import "github.com/desertbit/glue"

type Registry struct {
	socks map[string]map[*glue.Socket]bool
	rooms map[*glue.Socket]map[string]bool
}

func NewRegistry() *Registry {
	return &Registry{
		socks: map[string]map[*glue.Socket]bool{},
		rooms: map[*glue.Socket]map[string]bool{},
	}
}

func (r *Registry) register(sock *glue.Socket, room string) {
	if r.socks[room] == nil {
		r.socks[room] = map[*glue.Socket]bool{}
	}
	r.socks[room][sock] = true

	if r.rooms[sock] == nil {
		r.rooms[sock] = map[string]bool{}
	}
	r.rooms[sock][room] = true
}

func (r *Registry) deregister(sock *glue.Socket, room string) {
	_, ok := r.socks[room]
	if !ok {
		return
	}
	_, ok = r.rooms[sock]
	if !ok {
		return
	}

	delete(r.socks[room], sock)
	delete(r.rooms[sock], room)

	if len(r.socks[room]) == 0 {
		delete(r.socks, room)
	}
	if len(r.rooms[sock]) == 0 {
		delete(r.rooms, sock)
	}
}

func (r *Registry) deregisterAll(sock *glue.Socket) {
	_, ok := r.rooms[sock]
	if !ok {
		return
	}

	for room := range r.rooms[sock] {
		delete(r.socks[room], sock)
		if len(r.socks[room]) == 0 {
			delete(r.socks, room)
		}
	}
	delete(r.rooms, sock)
}

var registry = NewRegistry()
