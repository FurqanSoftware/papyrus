package ui

import (
	"bytes"
	"io"
	"net/http"
	"text/template"
)

func ServeHTMLTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
	buf := bytes.Buffer{}
	err := tpl.Execute(&buf, data)
	catch(r, err)
	w.Header().Set("Content-Type", "text/html")
	_, err = io.Copy(w, &buf)
	catch(r, err)
}

var (
	tplLayout = template.Must(template.New("layout.html").ParseFiles("ui/templates/layout.html", "ui/templates/common.html"))

	tplIndex = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/index.html"))
	tplLogin = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/login.html"))

	tplServeOrganization = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/organizationView.html"))
)
