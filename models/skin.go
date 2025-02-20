package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Skin struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Image       string             `bson:"image" json:"image"`
	Rarity      string             `bson:"rarity" json:"rarity"`               // e.g., Common, Rare, Epic, Legendary
	Condition   string             `bson:"condition" json:"condition"`         // e.g., Factory New, Minimal Wear, etc.
	OwnerID     primitive.ObjectID `bson:"owner_id,omitempty" json:"owner_id"` // null if available for sale
	IsListed    bool               `bson:"is_listed" json:"is_listed"`         // true if the skin is listed for sale
}
