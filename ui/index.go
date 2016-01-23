package ui

import (
	"mime"
	"net/http"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplIndex, struct {
		Context *Context
	}{
		Context: ctx,
	})
}

func init() {
	Router.NewRoute().
		Methods("GET").
		Path("/").
		HandlerFunc(ServeIndex)
}
