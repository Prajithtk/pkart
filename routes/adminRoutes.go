package routes

import (
	"pkart/controller"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	r.POST("/login", controller.AdminLogin)
	r.GET("/listusers", controller.ListUsers)
	r.PATCH("/blockuser/:ID", controller.BlockUser)
	r.GET("/viewcategory", controller.ViewCategory)
	r.POST("/addcategory", controller.AddCategory)
	r.PATCH("/editcategory/:ID", controller.EditCategory)
	r.PATCH("/blockcategory/:ID", controller.BlockCategory)
	r.DELETE("/deletecategory/:ID", controller.DeleteCategory)
	r.GET("/products", controller.ViewProducts)
	r.POST("/addproducts", controller.AddProducts)
	r.POST("/addimages", controller.ProductImage)
	r.PATCH("/editproducts/:ID", controller.EditProducts)
	r.DELETE("/deleteproducts/:ID", controller.DeleteProducts)

}
