package main

import (
	"Marketplace-cs2-/database"
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

type User struct {
	Name     string `bson:"name"`
	Password string `bson:"password"`
}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/protected", protected)
	http.ListenAndServe(":9000", nil)
	fmt.Println("Server running at :9000")
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "invalid method", err)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) < 8 || len(password) < 8 {
		err := http.StatusNotAcceptable
		http.Error(w, "Invalid password, username format", err)
		return
	}

	database.ConnectDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := database.GetCollection("users")

	existingUser := collection.FindOne(ctx, bson.M{"name": username})
	if existingUser.Err() == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	hashedPassword, _ := hashPassword(password)
	user := User{Name: username, Password: hashedPassword}
	collection.InsertOne(ctx, user)
	fmt.Println("User Registered!")
	fmt.Fprintf(w, "User registered")
}
func login(w http.ResponseWriter, r *http.Request)     {}
func logout(w http.ResponseWriter, r *http.Request)    {}
func protected(w http.ResponseWriter, r *http.Request) {}
