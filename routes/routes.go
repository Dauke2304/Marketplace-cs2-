package routes

import (
	"Marketplace-cs2-/controllers"
	"Marketplace-cs2-/services"
	"net/http"
)

func InitRoutes() {
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/protected", controllers.Protected)
	http.HandleFunc("/main", controllers.MainPage)
	http.HandleFunc("/profile", controllers.ProfilePage)
	http.HandleFunc("/buy-skin", controllers.BuySkin)
	http.HandleFunc("/sell-skin", controllers.SellSkin)

	http.HandleFunc("/admin", services.HandleAdminDashboard)
	http.HandleFunc("/admin/users", services.HandleAdminUsers)
	http.HandleFunc("/admin/skins", services.HandleAdminSkins)
	http.HandleFunc("/admin/transactions", services.HandleAdminTransactions)
	http.HandleFunc("/admin/delete-user", services.HandleDeleteUser)
	http.HandleFunc("/admin/delete-skin", services.HandleDeleteSkin)

	// For CSS an JS
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("../frontend/style"))))
	http.Handle("/styleTemplates/", http.StripPrefix("/styleTemplates/", http.FileServer(http.Dir("../frontend//templates//styleTemplates"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../frontend/templates/js"))))
}
