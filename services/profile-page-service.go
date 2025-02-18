package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/repositories"
	"fmt"
	"net/http"
	"text/template"
)

type Skin struct {
	ID        string
	Name      string
	Price     float64
	Image     string
	CSRFToken string
}

type ProfileData struct {
	Username  string
	SteamID   string
	Email     string
	Balance   float64
	Skins     []Skin
	CSRFToken string
}

func HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	// Initialize DB
	database.InitDB()
	repSkin := repositories.NewSkinRepository(database.Client.Database("cs2_skins_marketplace"))
	repUser := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))

	// Get session token
	cookie, err := r.Cookie("sessiontoken")
	if err != nil {
		http.Error(w, "Session token not found", http.StatusUnauthorized)
		return
	}

	// Get user from session token
	user, err := repUser.GetUserBySessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}

	// Get user skins
	skins, err := repSkin.GetSkinsByOwnerID(user.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve skins", http.StatusInternalServerError)
		return
	}

	// Convert to Skin struct
	var userSkins []Skin
	for _, skin := range skins {
		userSkins = append(userSkins, Skin{
			ID:    skin.ID.Hex(),
			Name:  skin.Name,
			Price: skin.Price,
			Image: skin.Image,
		})
	}

	// Prepare template data
	data := ProfileData{
		Username:  user.Username,
		SteamID:   user.SteamID,
		Email:     user.Email,
		Balance:   user.Balance,
		Skins:     userSkins,
		CSRFToken: user.CSRFToken,
	}
	fmt.Println("debug profile service......")
	fmt.Println(data)
	// Load and parse template
	tmpl := template.Must(template.ParseFiles("../frontend/templates/profile.html"))

	// Render template
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
