package routes

import (
	"Marketplace-cs2-/controllers"
	"net/http"
)

func InitRoutes() {
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/protected", controllers.Protected)
	http.HandleFunc("/main", controllers.MainPage)

	// For CSS an JS
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("C:\\Users\\Ernar\\Desktop\\Marketplace-cs2-\\frontend\\style"))))
	http.Handle("/styleTemplates/", http.StripPrefix("/styleTemplates/", http.FileServer(http.Dir("C:\\Users\\Ernar\\Desktop\\Marketplace-cs2-\\frontend\\templates\\styleTemplates"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("C:\\Users\\Ernar\\Desktop\\Marketplace-cs2-\\frontend\\templates\\js"))))
}
