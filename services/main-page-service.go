package services

import (
	"net/http"
	"path/filepath"
	"text/template"
)

type PageData struct {
	Skins []string
}

func HandleMainPage(w http.ResponseWriter, r *http.Request) {

	data := PageData{
		Skins: []string{"AK-47 | Redline", "M4A4 | Asiimov", "AWP | Dragon Lore"},
	}
	templatePath := filepath.Join("C:\\Users\\Ernar\\Desktop\\Marketplace-cs2-\\frontend\\templates\\index.html")
	// Parse the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
