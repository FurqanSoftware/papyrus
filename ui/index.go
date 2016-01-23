package ui

import (
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
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

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplLogin, struct {
		Context *Context
	}{
		Context: ctx,
	})
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	acc, err := data.GetAccountEmail(user.Email)
	catch(r, err)
	if acc == nil {
		accEmail, err := data.NewAccountEmail(user.Email)
		catch(r, err)
		accEmail.Primary = true
		accEmail.Verified = true
		accEmail.VerifiedAt = time.Now()

		nAcc := data.Account{}
		nAcc.Emails = append(nAcc.Emails, accEmail)
		err = nAcc.Put()

		acc = &nAcc
	}

	ctx.Session.Values["accountID"] = acc.ID.Hex()
	ctx.Session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	delete(ctx.Session.Values, "accountID")
	ctx.Session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func init() {
	goth.UseProviders(
		gplus.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("BASE")+"/login/gplus/callback", "email"),
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), os.Getenv("BASE")+"/login/github/callback"),
	)

	gothic.Store = store

	gothic.GetState = func(r *http.Request) string {
		return r.URL.Query().Get("state")
	}
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		provider, _ := mux.Vars(r)["provider"]
		return provider, nil
	}

	Router.NewRoute().
		Methods("GET").
		Path("/").
		HandlerFunc(ServeIndex)
	Router.NewRoute().
		Methods("GET").
		Path("/login").
		HandlerFunc(ServeLogin)
	Router.NewRoute().
		Methods("GET").
		Path("/logout").
		HandlerFunc(HandleLogout)
	Router.NewRoute().
		Methods("GET").
		Path("/login/{provider}").
		HandlerFunc(gothic.BeginAuthHandler)
	Router.NewRoute().
		Methods("GET").
		Path("/login/{provider}/callback").
		HandlerFunc(HandleAuthCallback)

}
