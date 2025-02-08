package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/repositories"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

type PageData struct {
	Items     []Item
	CSRFToken string
}

type Item struct {
	ID        string
	Name      string
	Price     float64
	Image     string
	CSRFToken string
}

func HandleMainPage(w http.ResponseWriter, r *http.Request) {

	database.InitDB()
	rep := repositories.NewSkinRepository(database.Client.Database("cs2_skins_marketplace"))
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
	fmt.Println("FROM MAIN SERVICE")
	fmt.Println(user.CSRFToken)
	skins, err := rep.GetListedSkins()
	if err != nil {
		http.Error(w, "Failed to retrieve skins", http.StatusInternalServerError)
		return
	}

	var items []Item
	for _, skin := range skins {
		items = append(items, Item{
			ID:        skin.ID.Hex(),
			Name:      skin.Name,
			Price:     skin.Price,
			Image:     skin.Image,
			CSRFToken: user.CSRFToken,
		})
	}

	data := PageData{
		Items:     items,
		CSRFToken: user.CSRFToken, // Pass the CSRF token globally
	}

	templatePath := filepath.Join("C:\\Users\\Ernar\\Desktop\\Marketplace-cs2-\\frontend\\templates\\index.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
