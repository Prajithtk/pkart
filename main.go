package main

import (
	"fmt"
	"os"
	"pkart/database"
	"pkart/helper"
	"pkart/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	helper.Envload()
	database.DbConnect()
}

func main() {
	fmt.Println("Welcome to Pkart")

	router := gin.Default()

	admin := router.Group("/admin")
	routes.AdminRoutes(admin)

	user := router.Group("/user")
	routes.UserRoutes(user)

	// router.POST("/logine", controller.AdminLogin)
	// router.Run(":9090")
	port := os.Getenv("PORT")
	if port == "" {
		port = "0909"
	}
	router.Run(":" + port)

}
