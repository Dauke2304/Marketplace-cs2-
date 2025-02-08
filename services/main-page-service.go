package services

import (
	"net/http"
	"path/filepath"
	"text/template"
)

type PageData struct {
	Items []Item // Or whatever your data is
}

type Item struct {
	Name  string
	Price float64
	Image string
}

func HandleMainPage(w http.ResponseWriter, r *http.Request) {

	data := PageData{
		Items: []Item{
			{Name: "AK-47 | Redline", Price: 100, Image: "redline.png"},
			{Name: "AWP | Dragon Lore", Price: 200, Image: "dragon.png"},
			{Name: "M4A4 | Howl", Price: 150, Image: "howl.png"},
			{Name: "Desert Eagle | Blaze", Price: 120, Image: "blaze.png"},
			{Name: "Glock-18 | Water Elemental", Price: 80, Image: "water_elemental.png"},
			{Name: "USP-S | Kill Confirmed", Price: 130, Image: "kill_confirmed.png"},
			{Name: "Karambit | Fade", Price: 350, Image: "fade.png"},
			{Name: "Butterfly Knife | Slaughter", Price: 450, Image: "slaughter.png"},
			{Name: "AK-47 | Fire Serpent", Price: 250, Image: "fire_serpent.png"},
			{Name: "AWP | Medusa", Price: 500, Image: "medusa.png"},
		},
	}

	templatePath := filepath.Join("./frontend/templates/index.html")
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
