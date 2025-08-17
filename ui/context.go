package ui

import (
	"net/http"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gorilla/sessions"
)

type Context struct {
	Request *http.Request

	Session *sessions.Session
	Account *data.Account
}

func GetContext(r *http.Request) *Context {
	ctx, _ := r.Context().Value("context").(*Context)
	return ctx
}
