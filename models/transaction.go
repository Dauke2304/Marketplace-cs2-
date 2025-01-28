package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	BuyerID  primitive.ObjectID `bson:"buyer_id"`
	SellerID primitive.ObjectID `bson:"seller_id"`
	SkinID   primitive.ObjectID `bson:"skin_id"`
	Amount   float64            `bson:"amount"`
	Date     time.Time          `bson:"date"`
}
