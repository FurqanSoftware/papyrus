package ui

import (
	"math"
	"mime"
	"net/http"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func ServeProject(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if !bson.IsObjectIdHex(idStr) {
		ServeNotFound(w, r)
		return
	}
	id := bson.ObjectIdHex(idStr)
	prj, err := data.GetProject(id)
	catch(r, err)
	if prj == nil {
		ServeNotFound(w, r)
		return
	}

	mems, err := prj.Members()
	catch(r, err)

	found := false

	for _, mem := range mems {
		if mem.AccountID == ctx.Account.ID {
			found = true
			break
		}
	}

	if !found {
		ServeForbidden(w, r)
		return
	}

	docs, err := data.ListDocumentsProject(prj.ID, 0, math.MaxInt32)
	catch(r, err)

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplServeProject, struct {
		Context   *Context
		Project   *data.Project
		Members   []data.Member
		Documents []data.Document
	}{
		Context:   ctx,
		Project:   prj,
		Members:   mems,
		Documents: docs,
	})
}

func init() {
	Router.NewRoute().Methods("GET").Path("/projects/{id}").HandlerFunc(ServeProject)
}
