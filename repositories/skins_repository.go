package repositories

import (
	"Marketplace-cs2-/models"
	"context"

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
