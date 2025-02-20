package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/repositories"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleSellSkin(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Debug 0....................")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println("Debug 1....................")
	ValidateAuthorization(r)
	fmt.Println("Debug 2....................")
	database.InitDB()
	rep := repositories.NewSkinRepository(database.Client.Database("cs2_skins_marketplace"))
	var skinID string
	if r.Header.Get("Content-Type") == "application/json" {
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err == nil {
			skinID = data["skinID"]
		}
	}
	fmt.Println("Debug 4....................")
	if skinID == "" {
		if err := r.ParseForm(); err == nil {
			skinID = r.FormValue("skinID")
		}
	}

	if skinID == "" {
		http.Error(w, "Missing skin ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(skinID)
	if err != nil {
		http.Error(w, "Invalid skin ID", http.StatusBadRequest)
		return
	}

	skin, err := rep.GetSkinByID(objID) // Now pass the converted ObjectID
	if err != nil || skin == nil {
		http.Error(w, "Skin not found", http.StatusNotFound)
		return
	}
	err1 := rep.ToggleIsListed(objID, true) // Set is_listed to false
	if err1 != nil {
		fmt.Println("Failed to update is_listed:", err)
	}
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
