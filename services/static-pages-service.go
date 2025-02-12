package services

import (
	"net/http"
	"text/template"
)

func HandleContactPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../frontend/templates/contact.html"))
	tmpl.Execute(w, nil)
}

func HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../frontend/templates/about.html"))
	tmpl.Execute(w, nil)
}
