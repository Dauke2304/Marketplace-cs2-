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

	tmpl := template.Must(template.ParseFiles("./frontend/templates/profile.html"))

	data := ProfileData{
		Username: "USERNAME",
		Avatar:   "avatar.png",
		Skins: []Skin{
			{Name: "AK-47 | Redline", Price: 100, Image: "redline.png"},
			{Name: "AWP | Dragon Lore", Price: 200, Image: "dragon.png"},
			{Name: "USP-S | Kill Confirmed", Price: 130, Image: "kill_confirmed.png"},
			{Name: "Karambit | Fade", Price: 350, Image: "fade.png"},
			{Name: "Butterfly Knife | Slaughter", Price: 450, Image: "slaughter.png"},
		},
	}

	tmpl.Execute(w, data)
}
