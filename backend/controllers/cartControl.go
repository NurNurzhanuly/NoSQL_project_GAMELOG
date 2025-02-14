package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/database"
	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Add item to cart
func AddToCart(c *gin.Context) {
	var cartItem model.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetString("userID") // Assume userID is set in middleware
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	collection := database.GetCollection("carts")
	filter := bson.M{"user_id": userObjectID}
	update := bson.M{"$push": bson.M{"items": cartItem}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error adding to cart:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

// Get cart items
func GetCart(c *gin.Context) {
	userID := c.GetString("userID") // Assume userID is set in middleware
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	collection := database.GetCollection("cart")
	var cart model.Cart
	err := collection.FindOne(context.TODO(), bson.M{"user_id": userObjectID}).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// Remove item from cart
func RemoveFromCart(c *gin.Context) {
	var cartItem model.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetString("userID") // Assume userID is set in middleware
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	collection := database.GetCollection("carts")
	filter := bson.M{"user_id": userObjectID}
	update := bson.M{"$pull": bson.M{"items": bson.M{"game_id": cartItem.GameID}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error removing from cart:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}
