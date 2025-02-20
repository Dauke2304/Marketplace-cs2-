package routes

import (
	"Marketplace-cs2-/controllers"
	"Marketplace-cs2-/services"
	"net/http"
)

func InitRoutes() {
	// Authentication routes
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)

	// General pages
	http.HandleFunc("/protected", controllers.Protected)
	http.HandleFunc("/main", controllers.MainPage)
	http.HandleFunc("/profile", controllers.ProfilePage)
	http.HandleFunc("/about", controllers.AboutPage)
	http.HandleFunc("/contact", controllers.ContactPage)

	// Skin-related routes
	http.HandleFunc("/buy-skin", controllers.BuySkin)
	http.HandleFunc("/sell-skin", controllers.SellSkin)

	// Admin panel
	http.HandleFunc("/admin", services.HandleAdminPanel)

	// User management routes
	http.HandleFunc("/admin/users", services.HandleAddUser)
	http.HandleFunc("/admin/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			services.HandleGetUser(w, r)
		case http.MethodPut:
			services.HandleUpdateUser(w, r)
		case http.MethodDelete:
			services.HandleDeleteUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Skin management routes
	http.HandleFunc("/admin/skins", services.HandleAddSkin)
	http.HandleFunc("/admin/skins/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			services.HandleGetSkin(w, r)
		case http.MethodPut:
			services.HandleUpdateSkin(w, r)
		case http.MethodDelete:
			services.HandleDeleteSkin(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Serve static files (CSS & JS)
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("../frontend/style"))))
	http.Handle("/styleTemplates/", http.StripPrefix("/styleTemplates/", http.FileServer(http.Dir("../frontend/templates/styleTemplates"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../frontend/templates/js"))))
}
