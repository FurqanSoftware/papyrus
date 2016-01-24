package ui

import (
	"math"
	"mime"
	"net/http"
	"time"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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

func HandleMemberAdd(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	catch(r, err)

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

	if prj.OwnerID != ctx.Account.ID {
		ServeForbidden(w, r)
		return
	}

	body := struct {
		Email string `schema:"email"`
	}{}

	err = schema.NewDecoder().Decode(&body, r.PostForm)
	catch(r, err)

	acc, err := data.GetAccountEmail(body.Email)
	catch(r, err)

	if acc == nil {
		RedirectBack(w, r)
		return
	}

	mem, err := data.GetMemberProjectAccount(prj.ID, acc.ID)
	catch(r, err)

	if mem != nil {
		RedirectBack(w, r)
		return
	}

	nM := data.Member{
		OrganizationID: prj.OrganizationID,
		ProjectID:      prj.ID,
		AccountID:      acc.ID,
		InviterID:      ctx.Account.ID,
		InvitedAt:      time.Now(),
	}
	err = nM.Put()
	catch(r, err)

	mems, err := data.ListMembersProject(prj.ID, 0, math.MaxInt32)
	catch(r, err)

	prj.MemberIDs = []bson.ObjectId{}
	for _, mem := range mems {
		prj.MemberIDs = append(prj.MemberIDs, mem.ID)
	}
	err = prj.Put()
	catch(r, err)

	http.Redirect(w, r, "/projects/"+prj.ID.Hex(), http.StatusSeeOther)
}

func init() {
	Router.NewRoute().Methods("GET").Path("/projects/{id}").HandlerFunc(ServeProject)
	Router.NewRoute().Methods("POST").Path("/projects/{id}/members/add").HandlerFunc(HandleMemberAdd)
}
