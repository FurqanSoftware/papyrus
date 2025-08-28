package ui

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
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

	// tplIndex = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/index.html"))
	tplLogin = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/login.html"))

	tplOrganizationList = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/organizationList.html"))
	tplOrganizationNew  = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/organizationNewForm.html"))

	tplProjectList = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/projectList.html"))
	tplProjectNew  = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/projectNewForm.html"))
	tplProjectView = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/projectView.html"))

	tplDocumentNew        = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/documentNewForm.html"))
	tplDocumentView       = template.Must(template.Must(tplLayout.Clone()).ParseFiles("ui/templates/documentView.html"))
	tplDocumentViewPublic = template.Must(template.ParseFiles("ui/templates/documentViewPublic.html"))
)
