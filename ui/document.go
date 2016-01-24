package ui

import (
	"mime"
	"net/http"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2/bson"
)

func ServeNewDocument(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	acc := ctx.Account

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
	if acc.ID != prj.OwnerID {
		ServeForbidden(w, r)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplServeDocumentNew, struct {
		Context *Context
	}{
		Context: ctx,
	})
}

func HandleDocumentCreate(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	acc := ctx.Account

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
	if acc.ID != prj.OwnerID {
		ServeForbidden(w, r)
		return
	}

	err = r.ParseForm()
	catch(r, err)

	body := struct {
		Title string `schema:"title"`
	}{}

	err = schema.NewDecoder().Decode(&body, r.PostForm)
	catch(r, err)

	switch {
	case body.Title == "":
		RedirectBack(w, r)
		return
	}
	doc := data.Document{
		Title:     body.Title,
		ProjectID: prj.ID,
		Published: false,
	}
	err = doc.Put()
	catch(r, err)

	http.Redirect(w, r, "/projects/"+prj.ID.Hex(), http.StatusSeeOther)
}

func ServeDocument(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	acc := ctx.Account

	vars := mux.Vars(r)
	idStr := vars["id"]
	if !bson.IsObjectIdHex(idStr) {
		ServeNotFound(w, r)
		return
	}
	id := bson.ObjectIdHex(idStr)
	doc, err := data.GetDocument(id)
	catch(r, err)
	if doc == nil {
		ServeNotFound(w, r)
		return
	}

	mem, err := data.GetMemberProjectAccount(doc.ProjectID, acc.ID)
	catch(r, err)
	if mem == nil {
		ServeForbidden(w, r)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplServeDocument, struct {
		Context  *Context
		Document *data.Document
	}{
		Context:  ctx,
		Document: doc,
	})
}

func init() {
	Router.NewRoute().
		Methods("GET").
		Path("/projects/{id}/documents/new").
		HandlerFunc(ServeNewDocument)
	Router.NewRoute().
		Methods("POST").
		Path("/projects/{id}/documents/new").
		HandlerFunc(HandleDocumentCreate)
	Router.NewRoute().
		Methods("GET").
		Path("/documents/{id}").
		HandlerFunc(ServeDocument)
}
