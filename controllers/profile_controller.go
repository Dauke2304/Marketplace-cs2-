package controllers

import (
	"Marketplace-cs2-/services"
	"net/http"
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Call the service that handles rendering the template
		services.HandleProfilePage(w, r)
		return
	}
	// Return error for invalid method
	http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
}
