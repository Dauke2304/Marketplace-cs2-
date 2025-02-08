package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/repositories"
	"encoding/json"
	"fmt"
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
	//csrfTokenHeader, _ := url.QueryUnescape(r.Header.Get("X-CSRF-Token"))
	//fmt.Println("HEADER")
	//fmt.Println(csrfTokenHeader)
	//fmt.Println("skinid")
	//skinID1 := r.FormValue("skinID::::::::")
	//fmt.Println(skinID1)
	ValidateAuthorization(r)
	//fmt.Println("after validate")
	database.InitDB()
	rep := repositories.NewSkinRepository(database.Client.Database("cs2_skins_marketplace"))
	repUser := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))

	cookie, err := r.Cookie("sessiontoken")
	if err != nil {
		http.Error(w, "Session token not found", http.StatusUnauthorized)
		return
	}
	//fmt.Println("after validate 1")
	// Get user from session token
	user, err := repUser.GetUserBySessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}
	//fmt.Println("after validate 2")

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

	fmt.Println("ObjectID:", objID)
	fmt.Println("after validate 3")
	skin, err := rep.GetSkinByID(objID) // Now pass the converted ObjectID
	if err != nil || skin == nil {
		http.Error(w, "Skin not found", http.StatusNotFound)
		return
	}
	fmt.Println("after validate 4")

	if user.Balance < skin.Price {
		http.Error(w, "Insufficient balance", http.StatusForbidden)
		return
	}
	fmt.Println("after validate 5")

	// 6️Process the purchase
	err = rep.TransferSkinOwnership(skinID, user.ID.Hex(), skin.Price)
	if err != nil {
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}
	fmt.Println("SUCCESS")
	// 7️Redirect to refresh the page
	http.Redirect(w, r, "/main", http.StatusSeeOther)

}
