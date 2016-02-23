package ui

import (
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gophergala2016/papyrus/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2"
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

	org, err := prj.Organization()
	catch(r, err)

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplDocumentNew, struct {
		Context      *Context
		Organization *data.Organization
		Project      *data.Project
	}{
		Context:      ctx,
		Organization: org,
		Project:      prj,
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
	if doc == nil || doc.Deleted {
		ServeNotFound(w, r)
		return
	}

	mem, err := data.GetMemberProjectAccount(doc.ProjectID, acc.ID)
	catch(r, err)
	if mem == nil {
		ServeForbidden(w, r)
		return
	}

	prj, err := doc.Project()
	catch(r, err)

	org, err := prj.Organization()
	catch(r, err)

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["accountID"] = ctx.Account.ID.Hex()
	token.Claims["documentID"] = doc.ID.Hex()
	token.Claims["expires"] = time.Now().Add(time.Minute * 15).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	catch(r, err)

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplDocumentView, struct {
		Context      *Context
		Organization *data.Organization
		Project      *data.Project
		Document     *data.Document
		Token        string
	}{
		Context:      ctx,
		Organization: org,
		Project:      prj,
		Document:     doc,
		Token:        tokenString,
	})
}

func HandleDocumentPublish(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if !bson.IsObjectIdHex(idStr) {
		ServeNotFound(w, r)
		return
	}
	id := bson.ObjectIdHex(idStr)

	doc, err := data.GetDocument(id)
	catch(r, err)
	if doc == nil || doc.Deleted {
		ServeNotFound(w, r)
		return
	}

	prj, err := doc.Project()
	catch(r, err)

	if prj.OwnerID != ctx.Account.ID {
		ServeForbidden(w, r)
		return
	}
	if doc.Published {
		http.Redirect(w, r, "/documents/"+doc.ID.Hex(), http.StatusSeeOther)
		return
	}

	doc.Published = true
	doc.PublishedAt = time.Now()
	doc.ShortID, err = data.GenerateShortID()
	catch(r, err)
	err = doc.Put()
	for mgo.IsDup(err) {
		doc.ShortID, err = data.GenerateShortID()
		catch(r, err)
		err = doc.Put()
	}

	http.Redirect(w, r, "/documents/"+doc.ID.Hex(), http.StatusSeeOther)
}

func HandleDocumentUnpublish(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if !bson.IsObjectIdHex(idStr) {
		ServeNotFound(w, r)
		return
	}
	id := bson.ObjectIdHex(idStr)

	doc, err := data.GetDocument(id)
	catch(r, err)
	if doc == nil || doc.Deleted {
		ServeNotFound(w, r)
		return
	}

	prj, err := doc.Project()
	catch(r, err)

	if prj.OwnerID != ctx.Account.ID {
		ServeForbidden(w, r)
		return
	}
	if !doc.Published {
		http.Redirect(w, r, "/documents/"+doc.ID.Hex(), http.StatusSeeOther)
		return
	}

	doc.Published = false
	doc.PublishedAt = time.Time{}
	err = doc.Put()
	catch(r, err)

	http.Redirect(w, r, "/documents/"+doc.ID.Hex(), http.StatusSeeOther)
}

func HandleDocumentDelete(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if !bson.IsObjectIdHex(idStr) {
		ServeNotFound(w, r)
		return
	}
	id := bson.ObjectIdHex(idStr)

	doc, err := data.GetDocument(id)
	catch(r, err)
	if doc == nil || doc.Deleted {
		ServeNotFound(w, r)
		return
	}

	prj, err := doc.Project()
	catch(r, err)

	if prj.OwnerID != ctx.Account.ID {
		ServeForbidden(w, r)
		return
	}

	doc.Deleted = true
	doc.DeletedAt = time.Now()
	doc.Published = false
	doc.PublishedAt = time.Time{}
	err = doc.Put()
	catch(r, err)

	http.Redirect(w, r, "/projects/"+prj.ID.Hex(), http.StatusSeeOther)
}

func HandleDocumentUndelete(w http.ResponseWriter, r *http.Request) {

	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

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

	prj, err := doc.Project()
	catch(r, err)

	if prj.OwnerID != ctx.Account.ID {
		ServeForbidden(w, r)
		return
	}

	doc.Deleted = false
	doc.DeletedAt = time.Time{}
	err = doc.Put()
	catch(r, err)

	http.Redirect(w, r, "/projects/"+prj.ID.Hex(), http.StatusSeeOther)
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
	Router.NewRoute().
		Methods("POST").
		Path("/documents/{id}/publish").
		HandlerFunc(HandleDocumentPublish)
	Router.NewRoute().
		Methods("POST").
		Path("/documents/{id}/unpublish").
		HandlerFunc(HandleDocumentUnpublish)
	Router.NewRoute().
		Methods("POST").
		Path("/documents/{id}/delete").
		HandlerFunc(HandleDocumentDelete)
	Router.NewRoute().
		Methods("POST").
		Path("/documents/{id}/undelete").
		HandlerFunc(HandleDocumentUndelete)
}
