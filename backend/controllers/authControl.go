package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/database"
	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/model"
	jwtServices "github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	collection := database.GetCollection("users")
	var existingUser model.User

	// Check if the user already exists by email or username
	err := collection.FindOne(context.TODO(), bson.M{"$or": []bson.M{
		{"email": input.Email},
		{"username": input.Username},
	}}).Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A user with this email already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing the password"})
		return
	}

	// Create a new user
	newUser := model.User{
		ID:         primitive.NewObjectID(),
		Username:   input.Username,
		Email:      input.Email,
		Password:   string(hashedPassword),
		OwnedGames: []model.OwnedGame{}, // 🔥 Исправлено! Теперь это массив объектов
		CreatedAt:  time.Now(),
	}

	// Save the user to the database
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// search for the user in the database
	collection := database.GetCollection("users")
	var user model.User

	log.Println("Searching for user:", input.Username)

	err := collection.FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя"})
		return
	}

	// password check
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	// generate jwt token
	token, err := jwtServices.GenerateToken(user.ID.Hex(), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	// send response
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"token":    token,
		"email":    user.Email,
		"redirect": "/FrontEnd/public/index.html",
	})
}

func UpdateUser(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		log.Println("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Invalid userID format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, не занят ли email другим пользователем
	collection := database.GetCollection("users")
	var existingUser model.User
	err = collection.FindOne(context.TODO(), bson.M{"email": input.Email, "_id": bson.M{"$ne": userObjectID}}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	// Обновляем пользователя в MongoDB
	update := bson.M{"$set": bson.M{"username": input.Username, "email": input.Email}}
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": userObjectID}, update)
	if err != nil {
		log.Println("Error updating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления пользователя"})
		return
	}

	if result.ModifiedCount == 0 {
		log.Println("User not found or no changes applied")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or no changes applied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func GetUserLibrary(c *gin.Context) {
	// Получаем userID из контекста
	userID := c.GetString("userID")
	log.Println("Extracted userID from middleware:", userID)

	if userID == "" {
		log.Println("Unauthorized access: userID is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var user model.User
	userCollection := database.GetCollection("users")

	// Преобразуем userID в ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Invalid user ID format:", userID, "Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	log.Println("Extracted userObjectID:", userObjectID.Hex())

	// Поиск пользователя в базе
	err = userCollection.FindOne(c, bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		log.Println("User not found with ID:", userObjectID.Hex(), "Error:", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	log.Println("User found:", user.Username, "OwnedGames count:", len(user.OwnedGames))

	// Проверяем, есть ли вообще игры
	if len(user.OwnedGames) == 0 {
		log.Println("User has no owned games:", userObjectID.Hex())
		c.JSON(http.StatusOK, gin.H{"message": "No games in library", "games": []model.Game{}})
		return
	}

	// 🔥 Исправление: Теперь owned_games - массив объектов, достаем GameID
	var gameIDs []primitive.ObjectID
	for _, item := range user.OwnedGames {
		gameIDs = append(gameIDs, item.GameID)
	}

	log.Println("Converted gameIDs:", gameIDs)

	// Получаем список игр
	gameCollection := database.GetCollection("games")
	var games []model.Game

	filter := bson.M{"_id": bson.M{"$in": gameIDs}}
	cursor, err := gameCollection.Find(c, filter)
	if err != nil {
		log.Println("Error retrieving games for user:", userObjectID.Hex(), "Filter:", filter, "Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving games"})
		return
	}
	defer cursor.Close(c)

	// Декодируем игры
	if err = cursor.All(c, &games); err != nil {
		log.Println("Error decoding games for user:", userObjectID.Hex(), "Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding games"})
		return
	}

	if len(games) == 0 {
		log.Println("No games found in database for user:", userObjectID.Hex())
		c.JSON(http.StatusOK, gin.H{"message": "No games found", "games": []model.Game{}})
		return
	}

	log.Println("Library retrieved successfully for user:", userObjectID.Hex(), "Total games:", len(games))
	c.JSON(http.StatusOK, games)
}
