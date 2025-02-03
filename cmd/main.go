package main

import (
	"Marketplace-cs2-/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	routes.InitRoutes()
	log.Fatal(http.ListenAndServe(":9000", nil))
	fmt.Println("Server running at :9000")
}
