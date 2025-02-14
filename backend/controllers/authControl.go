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
	"golang.org/x/crypto/bcrypt"
)

// Регистрация пользователя
func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка на существование пользователя
	collection := database.GetCollection("users")
	var existingUser model.User
	err := collection.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хеширования пароля"})
		return
	}

	// Создаем нового пользователя
	newUser := model.User{
		Username:   input.Username,
		Email:      input.Email,
		Password:   string(hashedPassword),
		OwnedGames: []primitive.ObjectID{},
		CreatedAt:  time.Now(),
	}

	// Сохранение пользователя в MongoDB
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

// Логин пользователя
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Проверка на корректность входных данных
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск пользователя по username
	collection := database.GetCollection("users")
	var user model.User
	err := collection.FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	// Генерация JWT токена
	token, err := jwtServices.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	// Возвращаем успешный ответ с токеном и сообщением + email
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"token":    token,
		"email":    user.Email, // <---- Добавили email
		"redirect": "/FrontEnd/public/index.html",
	})
}

// UpdateUser обновляет имя пользователя и email
func UpdateUser(c *gin.Context) {
	userID := c.GetString("userID") // Получаем ID пользователя из middleware
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

	// Проверяем, не занят ли email другим пользователем (опционально)
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
