package routes

import (
	"pkart/controller"
	"pkart/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	r.POST("/login", controller.AdminLogin)
	r.GET("/users", middleware.AuthMiddleware("admin"), controller.ListUsers)
	r.PATCH("/user/:ID", middleware.AuthMiddleware("admin"), controller.BlockUser)
	r.GET("/categories", middleware.AuthMiddleware("admin"), controller.ViewCategory)
	r.POST("/category", middleware.AuthMiddleware("admin"), controller.AddCategory)
	r.PUT("/category/:ID", middleware.AuthMiddleware("admin"), controller.EditCategory)
	r.PATCH("/category/:ID", middleware.AuthMiddleware("admin"), controller.BlockCategory)
	r.DELETE("/category/:ID", middleware.AuthMiddleware("admin"), controller.DeleteCategory)
	r.GET("/products", middleware.AuthMiddleware("admin"), controller.ViewProducts)
	r.POST("/product", middleware.AuthMiddleware("admin"), controller.AddProducts)
	r.POST("/images", middleware.AuthMiddleware("admin"), controller.ProductImage)
	r.PATCH("/product/:ID", middleware.AuthMiddleware("admin"), controller.EditProducts)
	r.DELETE("/product/:ID", middleware.AuthMiddleware("admin"), controller.DeleteProducts)

}
