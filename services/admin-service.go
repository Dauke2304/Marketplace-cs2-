package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/models"
	"Marketplace-cs2-/repositories"
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func HandleAdminPanel(w http.ResponseWriter, r *http.Request) {
	db := database.Client.Database("cs2_skins_marketplace")

	userRepo := repositories.NewUserRepository(db)
	skinRepo := repositories.NewSkinRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	users, _ := userRepo.GetAllUsers()
	skins, _ := skinRepo.GetAllSkins()
	transactions, _ := transactionRepo.GetAllTransactions()

	data := struct {
		Users        []models.User
		Skins        []models.Skin
		Transactions []models.Transaction
	}{
		Users:        users,
		Skins:        skins,
		Transactions: transactions,
	}

	tmpl := template.Must(template.ParseFiles("./frontend/templates/admin.html"))
	tmpl.Execute(w, data)
}

func HandleAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	db := database.Client.Database("cs2_skins_marketplace")
	userRepo := repositories.NewUserRepository(db)

	balance, err := strconv.ParseFloat(r.FormValue("balance"), 64)
	if err != nil {
		http.Error(w, "Invalid balance format", http.StatusBadRequest)
		return
	}
	password := r.FormValue("password")
	hashedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user := models.User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: hashedPassword,
		Balance:  balance,
	}

	_, err = userRepo.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	db := database.Client.Database("cs2_skins_marketplace")
	userRepo := repositories.NewUserRepository(db)

	userID := r.URL.Path[len("/admin/users/"):]
	objID, _ := primitive.ObjectIDFromHex(userID)

	if err := userRepo.DeleteUser(objID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleAddSkin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	db := database.Client.Database("cs2_skins_marketplace")
	skinRepo := repositories.NewSkinRepository(db)

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	skin := models.Skin{
		Name:   r.FormValue("name"),
		Price:  price,
		Rarity: r.FormValue("rarity"),
		Image:  r.FormValue("image"),
	}

	_, err = skinRepo.CreateSkin(skin)
	if err != nil {
		http.Error(w, "Failed to create skin", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func HandleDeleteSkin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	db := database.Client.Database("cs2_skins_marketplace")
	skinRepo := repositories.NewSkinRepository(db)

	skinID := r.URL.Path[len("/admin/skins/"):]
	objID, _ := primitive.ObjectIDFromHex(skinID)

	if err := skinRepo.DeleteSkin(objID); err != nil {
		http.Error(w, "Failed to delete skin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Add these handlers

func parseFloat(value string) float64 {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return f
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Path[len("/admin/users/"):]
	objID, _ := primitive.ObjectIDFromHex(userID)

	db := database.Client.Database("cs2_skins_marketplace")
	userRepo := repositories.NewUserRepository(db)

	user, _ := userRepo.GetUserByID(objID)
	json.NewEncoder(w).Encode(user)
}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Path[len("/admin/users/"):]
	objID, _ := primitive.ObjectIDFromHex(userID)

	db := database.Client.Database("cs2_skins_marketplace")
	userRepo := repositories.NewUserRepository(db)

	update := bson.M{
		"username": r.FormValue("username"),
		"email":    r.FormValue("email"),
		"balance":  parseFloat(r.FormValue("balance")),
	}

	userRepo.UpdateUser(objID, update)
	w.WriteHeader(http.StatusOK)
}

func HandleGetSkin(w http.ResponseWriter, r *http.Request) {
	skinID := r.URL.Path[len("/admin/skins/"):]
	objID, _ := primitive.ObjectIDFromHex(skinID)

	db := database.Client.Database("cs2_skins_marketplace")
	skinRepo := repositories.NewSkinRepository(db)

	skin, _ := skinRepo.GetSkinByID(objID)
	json.NewEncoder(w).Encode(skin)
}

func HandleUpdateSkin(w http.ResponseWriter, r *http.Request) {
	skinID := r.URL.Path[len("/admin/skins/"):]
	objID, _ := primitive.ObjectIDFromHex(skinID)

	db := database.Client.Database("cs2_skins_marketplace")
	skinRepo := repositories.NewSkinRepository(db)

	update := bson.M{
		"name":   r.FormValue("name"),
		"price":  parseFloat(r.FormValue("price")),
		"rarity": r.FormValue("rarity"),
		"image":  r.FormValue("image"),
	}

	skinRepo.UpdateSkin(objID, update)
	w.WriteHeader(http.StatusOK)
}
