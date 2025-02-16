package model

import "time"

type Game struct {
	ID          string    `json:"id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Genre       []string  `json:"genre" bson:"genre"`
	ReleaseDate time.Time `json:"release_date" bson:"release_date"`
	Developer   string    `json:"developer" bson:"developer"`
	Publisher   string    `json:"publisher" bson:"publisher"`
	Platforms   []string  `json:"platforms" bson:"platforms"`
	ImageURL    string    `json:"image_url" bson:"image_url"` // Хранит относительный путь
}
