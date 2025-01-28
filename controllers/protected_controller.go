package controllers

import (
	"Marketplace-cs2-/services"
	"net/http"
)

func Protected(w http.ResponseWriter, r *http.Request) {
	services.HandleProtected(w, r)
}
