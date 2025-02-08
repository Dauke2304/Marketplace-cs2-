package controllers

import (
	"Marketplace-cs2-/services"
	"net/http"
)

func BuySkin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		services.HandleBuySkin(w, r)
		return
	}
	// Return error for invalid method
	http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
}
