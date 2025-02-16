package routes

import (
	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/controllers"
	"github.com/NurNurzhanuly/NoSQL_project_GAMELOG/backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/auth/register", controllers.Register)
	r.POST("/api/auth/login", controllers.Login)
	r.GET("/api/games/getall", controllers.GetAllGames)
	r.GET("/api/games/search", controllers.SearchGames)
	r.GET("/api/games/:game_id", controllers.GetGameByID)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "You are authorized"})
		})

		// User routes
		protected.POST("/user/update", controllers.UpdateUser)

		// Cart routes
		protected.POST("/cart", controllers.AddToCart)
		protected.DELETE("/cart/:game_id", controllers.RemoveFromCart)
		protected.GET("/cart", controllers.GetCart)
		protected.POST("/cart/purchase", controllers.PurchaseAllGames)
	}
}
