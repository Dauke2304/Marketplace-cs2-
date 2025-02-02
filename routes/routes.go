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
}
