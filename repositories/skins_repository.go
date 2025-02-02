package repositories

import (
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
