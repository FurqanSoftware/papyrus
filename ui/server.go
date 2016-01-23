package ui

import (
	"net/http"

	"github.com/gorilla/mux"
)

var Router = mux.NewRouter()

type Server struct {
	Router *mux.Router
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s.Router.ServeHTTP(w, r)
}

func NewServer() *Server {
	return &Server{
		Router: Router,
	}
}
