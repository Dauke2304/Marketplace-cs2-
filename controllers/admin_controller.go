package controllers

import (
	"Marketplace-cs2-/services"
	"net/http"
)

func AdminPanel(w http.ResponseWriter, r *http.Request) {

	services.HandleAdminPanel(w, r)
}
