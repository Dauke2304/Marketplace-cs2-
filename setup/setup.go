package setup

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/models"
	"Marketplace-cs2-/repositories"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	db := database.Client.Database("cs2_skins_marketplace")
	userRepo := repositories.NewUserRepository(db)

	adminUser, _ := userRepo.GetUserByUsername("admin")
	if adminUser == nil {
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		userRepo.CreateUser(models.User{
			Username: "admin",
			Password: string(hashedPass),
			IsAdmin:  true,
		})
	}
}
