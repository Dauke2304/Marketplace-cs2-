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

	// For CSS
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("C:\\Users\\Ernar\\Desktop\\Marketplace-cs2-\\frontend\\style"))))
}
