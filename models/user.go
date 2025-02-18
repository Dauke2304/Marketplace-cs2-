package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SteamID      string             `bson:"steam_id" json:"steam_id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"-"`
	Balance      float64            `bson:"balance" json:"balance"`
	SessionToken string             `bson:"sessiontoken"`
	CSRFToken    string             `bson:"csrftoken"`
	IsAdmin      bool               `bson:"is_admin" json:"is_admin"`
}
