package controllers

import (
	"Marketplace-cs2-/services"
	"net/http"
)

func SellSkin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		services.HandleSellSkin(w, r)
		return
	}
	// Return error for invalid method
	http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
}
