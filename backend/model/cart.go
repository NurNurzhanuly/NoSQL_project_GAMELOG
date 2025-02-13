package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartItem struct {
	GameID   primitive.ObjectID `bson:"game_id" json:"game_id"`
	Quantity int                `bson:"quantity" json:"quantity"`
}

type Cart struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items  []CartItem         `bson:"items" json:"items"`
}
