package main

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/routes"
	"Marketplace-cs2-/setup"
	"fmt"
	"log"
	"net/http"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	setup.SeedAdmin()
	routes.InitRoutes()
	log.Fatal(http.ListenAndServe(":9000", nil))
	fmt.Println("Server running at :9000")
}
