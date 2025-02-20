package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/repositories"
	"encoding/json"
	"net/http"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mu sync.Mutex

func HandleBuySkin(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ValidateAuthorization(r)
	database.InitDB()
	rep := repositories.NewSkinRepository(database.Client.Database("cs2_skins_marketplace"))
	repUser := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))

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

	var skinID string
	if r.Header.Get("Content-Type") == "application/json" {
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err == nil {
			skinID = data["skinID"]
		}
	}

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

	if user.Balance < skin.Price {
		http.Error(w, "Insufficient balance", http.StatusForbidden)
		return
	}

	// 6️Process the purchase
	err = rep.TransferSkinOwnership(skinID, user.ID.Hex(), skin.Price)
	if err != nil {
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}
	// 7️Redirect to refresh the page
	http.Redirect(w, r, "/main", http.StatusSeeOther)

}
