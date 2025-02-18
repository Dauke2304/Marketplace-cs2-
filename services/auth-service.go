package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/models"
	"Marketplace-cs2-/repositories"
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
	return base64.RawURLEncoding.EncodeToString(bytes)

}

func ValidateAuthorization(r *http.Request) error {
	database.InitDB()
	rep := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))
	fmt.Println("from validate")
	cookie, err := r.Cookie("sessiontoken")
	if err != nil {
		return fmt.Errorf("auth error: sessiontoken not found")
	}

	user, err := rep.GetUserBySessionToken(cookie.Value)
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

	database.InitDB()
	rep := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))

	existingUser, err := rep.GetUserByUsername(username)
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		fmt.Println("user already exist: ", existingUser.Username)
		return
	}

	hashedPassword, _ := hashPassword(password)
	user := models.User{Username: username, Password: hashedPassword}
	rep.CreateUser(user)
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

	database.InitDB()
	rep := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))

	user, err := rep.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		fmt.Println("user not found")
		return
	}

	if !checkPasswordHash(password, user.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	sessiontoken := generateToken(32)
	csrftoken := generateToken(32)

	err = rep.UpdateUser(
		user.ID,
		bson.M{"sessiontoken": sessiontoken, "csrftoken": csrftoken},
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

	if user.IsAdmin {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "invalid method", err)
		return
	}
	database.InitDB()
	rep := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))

	cookie, err := r.Cookie("sessiontoken")
	if err != nil {
		http.Error(w, "Session token not found", http.StatusUnauthorized)
		return
	}

	user, err := rep.GetUserBySessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}

	username := user.Username
	if username == "" {
		http.Error(w, "Username not found", http.StatusUnauthorized)
		return
	}
	fmt.Println(user.Username)
	fmt.Println(user.SessionToken)
	fmt.Println(user.ID)

	er := rep.UpdateUser(
		user.ID,
		bson.M{
			"sessiontoken": "",
			"csrftoken":    "",
		},
	)
	if er != nil {
		http.Error(w, "Failed to clear tokens from database", http.StatusInternalServerError)
		fmt.Println(err)
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
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	database.InitDB()
	rep := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))
	cookie, err := r.Cookie("sessiontoken")
	if err != nil {
		http.Error(w, "Session token not found", http.StatusUnauthorized)
		return
	}

	user, err := rep.GetUserBySessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}
	username := user.Username
	fmt.Println(username)
}
