package ui

import "text/template"

var (
	tplLayout = template.Must(template.New("layout.html").ParseFiles("ui/templates/layout.html", "ui/templates/common.html"))
)
