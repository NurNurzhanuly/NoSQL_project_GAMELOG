package routes

import (
	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/controllers"
	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/auth/register", controllers.Register)
	r.POST("/api/auth/login", controllers.Login)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "You are authorized"})
		})

		// User routes
		protected.POST("/user/update", controllers.UpdateUser) // Добавлен маршрут для обновления пользователя

		// Cart routes
		protected.POST("/cart", controllers.AddToCart)
		protected.DELETE("/cart", controllers.RemoveFromCart)

		// Order routes
		protected.POST("/orders", controllers.PlaceOrder)
		protected.GET("/orders", controllers.GetOrders)
	}
}
