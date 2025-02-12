package services

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/repositories"
)

type AdminService struct {
	userRepo        *repositories.UserRepository
	skinRepo        *repositories.SkinRepository
	transactionRepo *repositories.TransactionRepository
}

func NewAdminService() *AdminService {
	db := database.Client.Database("cs2_skins_marketplace")
	return &AdminService{
		userRepo:        repositories.NewUserRepository(db),
		skinRepo:        repositories.NewSkinRepository(db),
		transactionRepo: repositories.NewTransactionRepository(db),
	}
}

func (s *AdminService) GetAllUsers() ([]interface{}, error) {
	return s.userRepo.GetAllUsers()
}

func (s *AdminService) GetAllSkins() ([]interface{}, error) {
	return s.skinRepo.GetAllSkins()
}

func (s *AdminService) GetAllTransactions() ([]interface{}, error) {
	return s.transactionRepo.GetAllTransactions()
}
