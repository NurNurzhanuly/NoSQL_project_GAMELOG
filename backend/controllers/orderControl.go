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

// Place an order
func PlaceOrder(c *gin.Context) {
	userID := c.GetString("userID") // Assume userID is set in middleware
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Error converting userID to ObjectID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get cart items
	cartCollection := database.GetCollection("cart")
	var cart model.Cart
	err = cartCollection.FindOne(context.TODO(), bson.M{"user_id": userObjectID}).Decode(&cart)
	if err != nil {
		log.Println("Error getting cart:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// Calculate total
	var total float64
	for _, item := range cart.Items {
		total += float64(item.Quantity) * 59.99 // Assume each game costs $59.99
	}

	// Convert cart items to order items
	var orderItems []model.OrderItem
	for _, cartItem := range cart.Items {
		orderItems = append(orderItems, model.OrderItem{
			GameID:   cartItem.GameID,
			Quantity: cartItem.Quantity,
		})
	}

	// Create order
	order := model.Order{
		ID:     primitive.NewObjectID(),
		UserID: userObjectID,
		Items:  orderItems,
		Total:  total,
	}

	orderCollection := database.GetCollection("orders")
	_, err = orderCollection.InsertOne(context.TODO(), order)
	if err != nil {
		log.Println("Error placing order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error placing order"})
		return
	}

	// Clear cart
	_, err = cartCollection.DeleteOne(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		log.Println("Error clearing cart:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error clearing cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": order})
}

func GetOrders(c *gin.Context) {
	userID := c.GetString("userID")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Error converting userID to ObjectID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	log.Println("Fetching orders for user ID:", userObjectID)

	collection := database.GetCollection("orders")
	cursor, err := collection.Find(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		log.Println("Error fetching orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching orders"})
		return
	}
	defer cursor.Close(context.TODO())

	var orders []model.Order
	if err = cursor.All(context.TODO(), &orders); err != nil {
		log.Println("Error decoding orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding orders"})
		return
	}

	log.Println("Successfully fetched orders:", orders)

	c.JSON(http.StatusOK, orders)
}
