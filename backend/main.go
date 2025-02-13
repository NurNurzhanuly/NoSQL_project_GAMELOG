package main

import (
	"fmt"
	"log"

	"github.com/OrazKhairulla/Online-Gaming-Marketplace/backend/database"
	"github.com/OrazKhairulla/Online-Gaming-Marketplace/backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	r := gin.Default()

	// used to serve static files of front
	r.Static("/FrontEnd", "../FrontEnd")
	r.StaticFile("/", "../FrontEnd/public/index.html")

	routes.SetupRoutes(r)
	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
