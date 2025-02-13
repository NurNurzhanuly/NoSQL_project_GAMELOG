package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderItem struct {
	GameID   primitive.ObjectID `bson:"game_id" json:"game_id"`
	Quantity int                `bson:"quantity" json:"quantity"`
}

type Order struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items  []OrderItem        `bson:"items" json:"items"`
	Total  float64            `bson:"total" json:"total"`
}
