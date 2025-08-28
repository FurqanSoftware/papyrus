package ui

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"

	"gopkg.in/mgo.v2/bson"

	"github.com/gophergala2016/papyrus/data"
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

	ctx := Context{
		Request: r,
		Session: sess,
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

		ctx.Account = acc
	}

	r = r.WithContext(context.WithValue(r.Context(), contextKeyContext, &ctx))

	func() {
		defer func() {
			err := recover()
			if err != nil {
				switch err := err.(type) {
				case *net.OpError:
					if err.Err == syscall.EPIPE || err.Err == syscall.ECONNRESET {
						break
					}
				case error:
					log.Print(err)
					ServeInternalServerError(w, r)
				default:
					panic(err)
				}
			}
		}()
		s.Router.ServeHTTP(w, r)
	}()
}

func NewServer() *Server {
	return &Server{
		Router: Router,
	}
}
