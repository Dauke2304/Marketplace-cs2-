package repositories

import (
	"Marketplace-cs2-/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(db *mongo.Database) *TransactionRepository {
	return &TransactionRepository{
		collection: db.Collection("transactions"),
	}
}

// CreateTransaction inserts a new transaction into the database
func (r *TransactionRepository) CreateTransaction(transaction models.Transaction) (primitive.ObjectID, error) {
	transaction.Date = time.Now()
	result, err := r.collection.InsertOne(context.Background(), transaction)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *TransactionRepository) GetTransactionByID(transactionID primitive.ObjectID) (*models.Transaction, error) {
	var transaction models.Transaction
	filter := bson.M{"_id": transactionID}
	err := r.collection.FindOne(context.Background(), filter).Decode(&transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionsByUserID fetches all transactions for a specific user (buyer or seller)
func (r *TransactionRepository) GetTransactionsByUserID(userID primitive.ObjectID) ([]models.Transaction, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"buyer_id": userID},
			{"seller_id": userID},
		},
	}
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var transactions []models.Transaction
	if err := cursor.All(context.Background(), &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) UpdateTransaction(transactionID primitive.ObjectID, update bson.M) error {
	filter := bson.M{"_id": transactionID}
	_, err := r.collection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	return err
}

func (r *TransactionRepository) DeleteTransaction(transactionID primitive.ObjectID) error {
	filter := bson.M{"_id": transactionID}
	_, err := r.collection.DeleteOne(context.Background(), filter)
	return err
}

func (r *TransactionRepository) GetAllTransactions() ([]models.Transaction, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var transactions []models.Transaction
	if err := cursor.All(context.Background(), &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}
