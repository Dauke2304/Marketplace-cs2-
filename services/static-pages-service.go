package services

import (
	"net/http"
	"text/template"
)

func HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../frontend/templates/about.html"))
	tmpl.Execute(w, nil)
}

func HandleContactPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle form submission
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("../frontend/templates/contact.html"))
	tmpl.Execute(w, nil)
}
