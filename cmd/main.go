package main

import (
	"Marketplace-cs2-/models"
	"Marketplace-cs2-/repositories"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	// Initialize MongoDB connection
	uri := "mongodb://localhost:27017" // Local MongoDB without authentication
	clientOptions := options.Client().ApplyURI(uri)

	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")
}

func main() {
	// Use the "cs2_marketplace" database
	db := client.Database("cs2_marketplace")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	skinRepo := repositories.NewSkinRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Test User Repository
	fmt.Println("Testing User Repository...")
	userID, err := testUserRepository(userRepo)
	if err != nil {
		log.Fatalf("User Repository test failed: %v", err)
	}

	// Test Skin Repository
	fmt.Println("\nTesting Skin Repository...")
	skinID, err := testSkinRepository(skinRepo, userID)
	if err != nil {
		log.Fatalf("Skin Repository test failed: %v", err)
	}

	// Test Transaction Repository
	fmt.Println("\nTesting Transaction Repository...")
	err = testTransactionRepository(transactionRepo, userID, skinID)
	if err != nil {
		log.Fatalf("Transaction Repository test failed: %v", err)
	}

	fmt.Println("\nAll repository tests completed successfully!")
}

// Test User Repository
func testUserRepository(userRepo *repositories.UserRepository) (primitive.ObjectID, error) {
	// Create a new user
	user := models.User{
		SteamID:  "12345",
		Username: "test_user",
		Email:    "test@example.com",
		Password: "password123",
		Balance:  100.0,
	}
	userID, err := userRepo.CreateUser(user)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %s\n", userID.Hex())

	// Fetch the user by ID
	fetchedUser, err := userRepo.GetUserByID(userID)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to fetch user: %v", err)
	}
	fmt.Printf("Fetched user: %+v\n", fetchedUser)

	// Update the user's balance
	update := bson.M{"balance": 200.0}
	err = userRepo.UpdateUser(userID, update)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to update user: %v", err)
	}
	fmt.Println("Updated user's balance to 200.0")

	// Fetch all users
	users, err := userRepo.GetAllUsers()
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to fetch all users: %v", err)
	}
	fmt.Printf("All users: %+v\n", users)

	return userID, nil
}

// Test Skin Repository
func testSkinRepository(skinRepo *repositories.SkinRepository, userID primitive.ObjectID) (primitive.ObjectID, error) {
	// Create a new skin
	skin := models.Skin{
		Name:    "Dragon Lore",
		Price:   500.0,
		Image:   "https://example.com/dragon_lore.jpg",
		OwnerID: userID,
	}
	skinID, err := skinRepo.CreateSkin(skin)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create skin: %v", err)
	}
	fmt.Printf("Created skin with ID: %s\n", skinID.Hex())

	// Fetch the skin by ID
	fetchedSkin, err := skinRepo.GetSkinByID(skinID)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to fetch skin: %v", err)
	}
	fmt.Printf("Fetched skin: %+v\n", fetchedSkin)

	// Fetch skins by owner ID
	skins, err := skinRepo.GetSkinsByOwnerID(userID)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to fetch skins by owner ID: %v", err)
	}
	fmt.Printf("Skins owned by user: %+v\n", skins)

	// Fetch all skins
	allSkins, err := skinRepo.GetAllSkins()
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to fetch all skins: %v", err)
	}
	fmt.Printf("All skins: %+v\n", allSkins)

	return skinID, nil
}

// Test Transaction Repository
func testTransactionRepository(transactionRepo *repositories.TransactionRepository, userID, skinID primitive.ObjectID) error {
	// Create a new transaction
	transaction := models.Transaction{
		BuyerID:  userID,
		SellerID: primitive.NilObjectID, // Assume no seller for this test
		SkinID:   skinID,
		Amount:   500.0,
	}
	transactionID, err := transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}
	fmt.Printf("Created transaction with ID: %s\n", transactionID.Hex())

	// Fetch the transaction by ID
	fetchedTransaction, err := transactionRepo.GetTransactionByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to fetch transaction: %v", err)
	}
	fmt.Printf("Fetched transaction: %+v\n", fetchedTransaction)

	// Fetch transactions by user ID
	transactions, err := transactionRepo.GetTransactionsByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch transactions by user ID: %v", err)
	}
	fmt.Printf("Transactions involving user: %+v\n", transactions)

	// Fetch all transactions
	allTransactions, err := transactionRepo.GetAllTransactions()
	if err != nil {
		return fmt.Errorf("failed to fetch all transactions: %v", err)
	}
	fmt.Printf("All transactions: %+v\n", allTransactions)

	return nil
}
