package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Skin struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Price   float64            `bson:"price"`
	Image   string             `bson:"image"`
	OwnerID primitive.ObjectID `bson:"owner_id,omitempty"` // null if available for sale
}
