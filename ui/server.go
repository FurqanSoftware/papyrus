package ui

import (
	"net/http"
	"os"

	"gopkg.in/mgo.v2/bson"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var Router = mux.NewRouter()

var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET")))

type Server struct {
	Router *mux.Router
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, err := store.Get(r, "s")
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}

	accID, ok := sess.Values["accountID"].(string)
	if ok {
		if !bson.IsObjectIdHex(accID) {
			ServeBadRequest(w, r)
			return
		}
		acc, err := data.GetAccount(bson.ObjectIdHex(accID))
		if err != nil {
			ServeInternalServerError(w, r)
			return
		}
		context.Set(r, "account", acc)
	}

	s.Router.ServeHTTP(w, r)
}

func NewServer() *Server {
	return &Server{
		Router: Router,
	}
}
