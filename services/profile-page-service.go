package services

import (
	"net/http"
	"text/template"
)

type Skin struct {
	Name  string
	Price int
	Image string
}

type ProfileData struct {
	Username string
	Avatar   string
	Skins    []Skin
}

func HandleProfilePage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("C:\\Users\\Ernar\\Desktop\\Marketplace-cs2-\\frontend\\templates\\profile.html"))

	data := ProfileData{
		Username: "WITOX",
		Avatar:   "avatar.png",
		Skins: []Skin{
			{Name: "AK-47 | Redline", Price: 100, Image: "redline.png"},
			{Name: "AWP | Dragon Lore", Price: 200, Image: "dragon.png"},
		},
	}

	tmpl.Execute(w, data)
}
