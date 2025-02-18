package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/models"
	"Marketplace-cs2-/repositories"
	"net/http"
	"text/template"
)

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

	tmpl := template.Must(template.ParseFiles("../frontend/templates/admin.html"))
	tmpl.Execute(w, data)
}
func GetCurrentUser(r *http.Request) *models.User {
	// Get session token from cookies
	cookie, err := r.Cookie("sessiontoken")
	if err != nil {
		return nil
	}

	// Initialize database and repository
	db := database.Client.Database("cs2_skins_marketplace")
	userRepo := repositories.NewUserRepository(db)

	// Get user by session token
	user, err := userRepo.GetUserBySessionToken(cookie.Value)
	if err != nil {
		return nil
	}

	return user
}
