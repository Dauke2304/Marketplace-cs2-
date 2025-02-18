package controllers

import (
	"Marketplace-cs2-/services"
	"net/http"
)

func AboutPage(w http.ResponseWriter, r *http.Request) {
	services.HandleAboutPage(w, r)
}

func ContactPage(w http.ResponseWriter, r *http.Request) {
	services.HandleContactPage(w, r)
}
