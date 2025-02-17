package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OwnedGame struct {
	GameID primitive.ObjectID `bson:"game_id"`
}

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Username   string             `bson:"username"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	OwnedGames []OwnedGame        `bson:"owned_games"` // Теперь массив объектов
	CreatedAt  time.Time          `bson:"created_at"`
}
