package routes

import (
	controller "pkart/controller/admin"
	"pkart/middleware"

	"github.com/gin-gonic/gin"
)

var RoleAdmin = "Admin"

func AdminRoutes(r *gin.RouterGroup) {
	r.POST("/login", controller.AdminLogin)

	/////////users
	r.GET("/users", middleware.AuthMiddleware(RoleAdmin), controller.ListUsers)
	r.PATCH("/user/:ID", middleware.AuthMiddleware(RoleAdmin), controller.BlockUser)

	/////////categories
	r.GET("/categories", middleware.AuthMiddleware(RoleAdmin), controller.ViewCategory)
	r.POST("/category", middleware.AuthMiddleware(RoleAdmin), controller.AddCategory)
	r.PUT("/category/:ID", middleware.AuthMiddleware(RoleAdmin), controller.EditCategory)
	r.PATCH("/category/:ID", middleware.AuthMiddleware(RoleAdmin), controller.BlockCategory)
	r.DELETE("/category/:ID", middleware.AuthMiddleware(RoleAdmin), controller.DeleteCategory)

	/////////coupons
	r.POST("/coupon", middleware.AuthMiddleware(RoleAdmin), controller.AddCoupon)
	r.GET("/coupons", middleware.AuthMiddleware(RoleAdmin), controller.ViewCoupon)
	r.PATCH("/coupon", middleware.AuthMiddleware(RoleAdmin), controller.ViewCoupon)
	r.DELETE("/coupon", middleware.AuthMiddleware(RoleAdmin), controller.ViewCoupon)

	/////////Products
	r.GET("/products", middleware.AuthMiddleware(RoleAdmin), controller.ViewProducts)
	r.POST("/product", middleware.AuthMiddleware(RoleAdmin), controller.AddProducts)
	r.POST("/images", middleware.AuthMiddleware(RoleAdmin), controller.ProductImage)
	r.PATCH("/product/:ID", middleware.AuthMiddleware(RoleAdmin), controller.EditProducts)
	r.DELETE("/product/:ID", middleware.AuthMiddleware(RoleAdmin), controller.DeleteProducts)

	///////Orders
	r.GET("/orders", middleware.AuthMiddleware(RoleAdmin), controller.ShowOrders)
	r.GET("/order/status", middleware.AuthMiddleware(RoleAdmin), controller.OrdersStatusChange)

	//////report
	r.GET("/sales/report", middleware.AuthMiddleware(RoleAdmin), controller.SalesReport)
	r.GET("/report", middleware.AuthMiddleware(RoleAdmin), controller.GetReportData)
}
