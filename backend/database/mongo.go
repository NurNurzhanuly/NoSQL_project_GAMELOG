package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectDB() {
	uri := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ Can not connect to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Can not ping MongoDB:", err)
	}

	fmt.Println("✅ Successfully connected to MongoDB")
}

func GetCollection(collectionName string) *mongo.Collection {
	if client == nil {
		log.Fatal("❌ Database is not connected")
	}
	return client.Database("game_log").Collection(collectionName)
}

// CloseDB closes the MongoDB connection
func CloseDB() {
	if client != nil {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal("❌ Can not disconnect from MongoDB:", err)
		} else {
			fmt.Println("✅ Connection to MongoDB is closed")
		}
	}
}
