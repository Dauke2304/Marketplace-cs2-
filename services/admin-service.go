package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/repositories"
	"net/http"
	"text/template"
)

func ValidateAdmin(r *http.Request) bool {
	if err := ValidateAuthorization(r); err != nil {
		return false
	}

	cookie, _ := r.Cookie("sessiontoken")
	userRepo := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))
	user, _ := userRepo.GetUserBySessionToken(cookie.Value)

	return user.Username == "admin" // Special admin username
}

func HandleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	if !ValidateAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tmpl := template.Must(template.ParseFiles("../frontend/templates/admin-dashboard.html"))
	tmpl.Execute(w, nil)
}

func HandleAdminUsers(w http.ResponseWriter, r *http.Request) {
	if !ValidateAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userRepo := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))
	users, _ := userRepo.GetAllUsers()

	tmpl := template.Must(template.ParseFiles("../frontend/templates/admin-users.html"))
	tmpl.Execute(w, users)
}

func HandleAdminSkins(w http.ResponseWriter, r *http.Request) {
	if !ValidateAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	skinRepo := repositories.NewSkinRepository(database.Client.Database("cs2_skins_marketplace"))
	skins, _ := skinRepo.GetAllSkins()

	tmpl := template.Must(template.ParseFiles("../frontend/templates/admin-skins.html"))
	tmpl.Execute(w, skins)
}

func HandleAdminTransactions(w http.ResponseWriter, r *http.Request) {
	if !ValidateAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	transactionRepo := repositories.NewTransactionRepository(database.Client.Database("cs2_skins_marketplace"))
	transactions, _ := transactionRepo.GetAllTransactions()

	tmpl := template.Must(template.ParseFiles("../frontend/templates/admin-transactions.html"))
	tmpl.Execute(w, transactions)
}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if !ValidateAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Implementation using userRepo.DeleteUser()
}

func HandleDeleteSkin(w http.ResponseWriter, r *http.Request) {
	if !ValidateAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Implementation using skinRepo.DeleteSkin()
}
