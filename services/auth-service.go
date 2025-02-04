package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/models"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"time"

	"go.mongodb.org/mongo-driver/bson"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Printf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func ValidateAuthorization(r *http.Request) error {
	username := r.FormValue("username")

	database.ConnectDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := database.GetCollection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		fmt.Println("not found")
		return fmt.Errorf("auth error: user not found")
	}

	sessionTokenCookie, err := r.Cookie("sessiontoken")
	if err != nil || sessionTokenCookie.Value == "" || sessionTokenCookie.Value != user.SessionToken {
		fmt.Println("session bug")
		fmt.Println(sessionTokenCookie)
		fmt.Println(user.SessionToken)
		return fmt.Errorf("auth error: invalid session token")
	}

	csrfTokenHeader, _ := url.QueryUnescape(r.Header.Get("X-CSRF-Token"))
	fmt.Println(csrfTokenHeader)
	if csrfTokenHeader != user.CSRFToken || csrfTokenHeader == "" {
		fmt.Println("csrf bug")
		fmt.Println(csrfTokenHeader)
		fmt.Println(user.CSRFToken)
		fmt.Println("auth error: invalid CSRF token")

		return fmt.Errorf("auth error: invalid CSRF token")

	}

	// If both tokens are valid, return nil (authentication is successful)
	return nil
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
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

	existingUser := collection.FindOne(ctx, bson.M{"username": username})
	if existingUser.Err() == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	hashedPassword, _ := hashPassword(password)
	user := models.User{Username: username, Password: hashedPassword}
	collection.InsertOne(ctx, user)
	fmt.Println("User Registered!")
	fmt.Fprintf(w, "User registered")
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "invalid method", err)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	database.ConnectDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := database.GetCollection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !checkPasswordHash(password, user.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	sessiontoken := generateToken(32)
	csrftoken := generateToken(32)

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{
			"sessiontoken": sessiontoken,
			"csrftoken":    csrftoken,
		}},
	)
	if err != nil {
		http.Error(w, "Failed to store session or CSRF token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessiontoken",
		Value:    sessiontoken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrftoken",
		Value:    csrftoken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: false,
	})
	http.Redirect(w, r, "/main", http.StatusSeeOther)
	fmt.Fprintf(w, "Logged in, %s", username)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "invalid method", err)
		return
	}

	// Get the username from the request or session (e.g., from a cookie)
	username := r.FormValue("username")
	if username == "" {
		http.Error(w, "Username not found", http.StatusUnauthorized)
		return
	}

	// Connect to the database and clear the session and CSRF tokens
	database.ConnectDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := database.GetCollection("users")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.M{
			"$set": bson.M{
				"sessiontoken": "",
				"csrftoken":    "",
			},
		},
	)
	if err != nil {
		http.Error(w, "Failed to clear tokens from database", http.StatusInternalServerError)
		return
	}

	// Clear the session token and CSRF token from the cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "sessiontoken",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour), // Set expiration to the past to effectively delete it
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrftoken",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour), // Set expiration to the past to effectively delete it
	})
	fmt.Fprintf(w, "Logged out, %s", username)
}

func HandleProtected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "invalid method", err)
		return
	}
	err := ValidateAuthorization(r)
	if err != nil {
		err := http.StatusUnauthorized
		http.Error(w, "StatusUnauthorized", err)
		fmt.Println("From handle protected 1")
		return
	}
	username := r.FormValue("username")
	fmt.Fprintf(w, "Welcome, %s", username)
}
