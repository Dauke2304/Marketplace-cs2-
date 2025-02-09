package repositories

import (
	"Marketplace-cs2-/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(user models.User) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *UserRepository) GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	var user models.User
	filter := bson.M{"_id": userID}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	filter := bson.M{"username": username}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserBySessionToken(sessiontoken string) (*models.User, error) {
	var user models.User
	filter := bson.M{"sessiontoken": sessiontoken}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserBySteamID(steamID string) (*models.User, error) {
	var user models.User
	filter := bson.M{"steam_id": steamID}
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(userID primitive.ObjectID, update bson.M) error {
	filter := bson.M{"_id": userID}
	_, err := r.collection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	return err
}

func (r *UserRepository) DeleteUser(userID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	_, err := r.collection.DeleteOne(context.Background(), filter)
	return err
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUserByIDCtx(sessCtx context.Context, userID primitive.ObjectID) (*models.User, error) {
	// Ensure context is not expired
	if err := sessCtx.Err(); err != nil {
		fmt.Println("context error")
		return nil, fmt.Errorf("context error: %w", err)
	}

	var user models.User
	err := r.collection.FindOne(sessCtx, bson.M{"_id": userID}).Decode(&user)

	// 2 Handle "no user found" case separately
	if err == mongo.ErrNoDocuments {
		fmt.Println("no user")
		return nil, nil // No error, just no user
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err) // Wrap other errors
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserBalance(sessCtx context.Context, userID primitive.ObjectID, newBalance float64) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"balance": newBalance}}

	_, err := r.collection.UpdateOne(sessCtx, filter, update)
	return err
}
