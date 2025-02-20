package controllers

import (
	"Marketplace-cs2-/services"
	"net/http"
	"path/filepath"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, filepath.Join("./frontend", "register.html"))
		services.HandleRegister(w, r)
		return
	}
	if r.Method == http.MethodPost {
		services.HandleRegister(w, r)
		return
	}
	http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, filepath.Join("./frontend", "login.html"))
		services.HandleLogin(w, r)
		return
	}
	if r.Method == http.MethodPost {
		services.HandleLogin(w, r)
		return
	}
	http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	services.HandleLogout(w, r)
}
