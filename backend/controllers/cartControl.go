package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/database"
	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add item to cart
func AddToCart(c *gin.Context) {
	// Получение userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	// Преобразование userID в строку
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Преобразование userID в ObjectID
	userIDObj, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Обработка входных данных
	var input struct {
		GameID string `json:"game_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Конвертация GameID в ObjectID
	gameID, err := primitive.ObjectIDFromHex(input.GameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Game ID format"})
		return
	}

	collection := database.GetCollection("cart")
	var cart model.Cart

	// Проверка существования корзины
	err = collection.FindOne(context.TODO(), bson.M{"user_id": userIDObj}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		// Если корзина не найдена, создаем новую
		newCart := model.Cart{
			UserID:    userIDObj,
			Items:     []model.CartItem{{GameID: gameID}},
			UpdatedAt: time.Now(),
		}

		_, err := collection.InsertOne(context.TODO(), newCart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding to cart"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Game added to cart"})
		return
	}

	// Добавление игры в существующую корзину
	cart.Items = append(cart.Items, model.CartItem{GameID: gameID})
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"user_id": userIDObj},
		bson.M{"$set": bson.M{"items": cart.Items, "updated_at": time.Now()}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game added to cart"})
}

// Get cart items
func GetCart(c *gin.Context) {
	log.Println("GetCart endpoint hit") // Логирование для проверки
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string)) // Преобразование userID в ObjectID
	if err != nil {
		log.Println("Invalid user ID format:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	collection := database.GetCollection("cart")
	var cart model.Cart

	// Фильтр по user_id
	err = collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Cart not found for user:", userIDStr)
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		} else {
			log.Println("Error fetching cart:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching cart"})
		}
		return
	}

	// Возвращаем только товары пользователя
	c.JSON(http.StatusOK, cart)
}

// Remove item from cart
func RemoveFromCart(c *gin.Context) {
	// Извлечение userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	// Преобразование userID в строку
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Преобразование userID в ObjectID
	userIDObj, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Извлечение game_id из параметров URL
	gameID := c.Param("game_id")
	objectID, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Game ID format"})
		return
	}

	// Получение коллекции cart
	collection := database.GetCollection("cart")

	// Фильтр по user_id и обновление items
	filter := bson.M{"user_id": userIDObj}
	update := bson.M{"$pull": bson.M{"items": bson.M{"game_id": objectID}}}

	// Обновление корзины
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error removing from cart:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing from cart"})
		return
	}

	log.Println("Item removed from cart:", gameID)
	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}
