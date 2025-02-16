package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") // Используй тот же ключ

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Проверка формата токена "Bearer <токен>"
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			log.Println("Invalid token format:", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Парсинг токена с использованием jwt/v4
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid token or error parsing token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Println("Invalid claims in token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Извлечение user_id и добавление его в контекст
		userID, ok := claims["user_id"].(string) // Приведение к строке
		if !ok || strings.TrimSpace(userID) == "" {
			log.Println("user_id missing or invalid in token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user_id in token"})
			c.Abort()
			return
		}

		log.Println("Extracted user_id from token:", userID)
		c.Set("userID", userID)
		c.Next()
	}
}
