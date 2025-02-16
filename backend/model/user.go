package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Username   string               `bson:"username"`
	Email      string               `bson:"email"`
	Password   string               `bson:"password"`
	OwnedGames []primitive.ObjectID `bson:"owned_games"`
	CreatedAt  time.Time            `bson:"created_at"`
}
