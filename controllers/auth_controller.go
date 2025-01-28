package controllers

import (
	"Marketplace-cs2-/services"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	services.HandleRegister(w, r)
}

func Login(w http.ResponseWriter, r *http.Request) {
	services.HandleLogin(w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	services.HandleLogout(w, r)
}
