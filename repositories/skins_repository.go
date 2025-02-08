package repositories

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SkinRepository struct {
	collection *mongo.Collection
}

func NewSkinRepository(db *mongo.Database) *SkinRepository {
	return &SkinRepository{
		collection: db.Collection("skins"),
	}
}

// CreateSkin inserts a new skin into the database
func (r *SkinRepository) CreateSkin(skin models.Skin) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(context.Background(), skin)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *SkinRepository) GetSkinByID(skinID primitive.ObjectID) (*models.Skin, error) {
	var skin models.Skin
	filter := bson.M{"_id": skinID}
	err := r.collection.FindOne(context.Background(), filter).Decode(&skin)
	if err != nil {
		return nil, err
	}
	return &skin, nil
}

func (r *SkinRepository) GetSkinsByOwnerID(ownerID primitive.ObjectID) ([]models.Skin, error) {
	filter := bson.M{"owner_id": ownerID}
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var skins []models.Skin
	if err := cursor.All(context.Background(), &skins); err != nil {
		return nil, err
	}
	return skins, nil
}

func (r *SkinRepository) UpdateSkin(skinID primitive.ObjectID, update bson.M) error {
	filter := bson.M{"_id": skinID}
	_, err := r.collection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	return err
}

func (r *SkinRepository) DeleteSkin(skinID primitive.ObjectID) error {
	filter := bson.M{"_id": skinID}
	_, err := r.collection.DeleteOne(context.Background(), filter)
	return err
}

func (r *SkinRepository) GetAllSkins() ([]models.Skin, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var skins []models.Skin
	if err := cursor.All(context.Background(), &skins); err != nil {
		return nil, err
	}
	return skins, nil
}

func (r *SkinRepository) GetListedSkins() ([]models.Skin, error) {
	filter := bson.M{"is_listed": true}
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var skins []models.Skin
	if err := cursor.All(context.Background(), &skins); err != nil {
		return nil, err
	}
	return skins, nil
}

func (r *SkinRepository) GetSkinByIDctx(sessCtx context.Context, skinID primitive.ObjectID) (*models.Skin, error) {
	var skin models.Skin
	filter := bson.M{"_id": skinID}
	err := r.collection.FindOne(sessCtx, filter).Decode(&skin)
	if err != nil {
		return nil, err
	}
	return &skin, nil
}

func (r *SkinRepository) UpdateSkinOwner(sessCtx context.Context, skinID primitive.ObjectID, newOwnerID primitive.ObjectID) error {
	filter := bson.M{"_id": skinID}
	update := bson.M{"$set": bson.M{"owner_id": newOwnerID}}

	_, err := r.collection.UpdateOne(sessCtx, filter, update)
	return err
}

func (r *SkinRepository) TransferSkinOwnership(skinID string, buyerID string, price float64) error {
	// Convert string IDs to ObjectIDs
	objSkinID, err := primitive.ObjectIDFromHex(skinID)
	if err != nil {
		fmt.Println("invalid skin ID format")
		return fmt.Errorf("invalid skin ID format: %w", err)
	}
	objBuyerID, err := primitive.ObjectIDFromHex(buyerID)
	if err != nil {
		fmt.Println("invalid buyer ID format")
		return fmt.Errorf("invalid buyer ID format: %w", err)
	}

	// Initialize database and repositories
	db := database.Client.Database("cs2_skins_marketplace")
	userRepo := NewUserRepository(db)

	// Start a MongoDB session for transaction
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := database.Client.StartSession()
	if err != nil {
		fmt.Println("failed to start MongoDB sessio")
		return fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	// Transaction function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Get buyer details
		buyer, err := userRepo.GetUserByIDCtx(sessCtx, objBuyerID)
		if err != nil {
			fmt.Println("error retrieving buyer")
			return nil, fmt.Errorf("error retrieving buyer: %w", err)
		}
		if buyer == nil {
			fmt.Println("buyer not found")
			return nil, errors.New("buyer not found")
		}

		// Get skin details
		skin, err := r.GetSkinByIDctx(sessCtx, objSkinID)
		if err != nil {
			fmt.Println("error retrieving skin")
			return nil, fmt.Errorf("error retrieving skin: %w", err)
		}
		if skin == nil {
			fmt.Println("skin not found")
			return nil, errors.New("skin not found")
		}

		// Get current owner's details
		owner, err := userRepo.GetUserByIDCtx(sessCtx, skin.OwnerID)
		if err != nil {
			fmt.Println("error retrieving skin owner")
			return nil, fmt.Errorf("error retrieving skin owner: %w", err)
		}
		if owner == nil {
			fmt.Println("skin owner not found")
			return nil, errors.New("skin owner not found")
		}

		// Check buyer's balance
		if buyer.Balance < price {
			fmt.Println("insufficient funds")
			return nil, errors.New("insufficient funds")
		}

		// Update balances
		err = userRepo.UpdateUserBalance(sessCtx, objBuyerID, buyer.Balance-price)
		if err != nil {
			fmt.Println("error updating buyer balance: %w")
			return nil, fmt.Errorf("error updating buyer balance: %w", err)
		}

		err = userRepo.UpdateUserBalance(sessCtx, owner.ID, owner.Balance+price)
		if err != nil {
			fmt.Println("error updating owner balance: %w")
			return nil, fmt.Errorf("error updating owner balance: %w", err)
		}

		// Transfer skin ownership
		err = r.UpdateSkinOwner(sessCtx, objSkinID, objBuyerID)
		if err != nil {
			fmt.Println("error updating skin owner: %w")
			return nil, fmt.Errorf("error updating skin owner: %w", err)
		}

		return nil, nil
	}

	// Execute the transaction
	_, err = callback(mongo.NewSessionContext(ctx, session))

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}
