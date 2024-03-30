package routes

import (
	"pkart/controller"
	"pkart/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	r.POST("/signup", controller.UserSignUp)
	r.POST("/login", controller.UserLogin)
	r.GET("/otp", controller.OtpSignUp)
	r.GET("/resendotp", controller.ResendOtp)
	
	///products
	r.GET("/products", middleware.AuthMiddleware("user"), controller.UserViewProducts)
	r.GET("/sort",middleware.AuthMiddleware("user"), controller.SortProduct)
	
	///address
	r.POST("/address", middleware.AuthMiddleware("user"), controller.AddAddress)
	r.PUT("/address/:ID", middleware.AuthMiddleware("user"), controller.EditAddress)
	r.DELETE("/address/:ID", middleware.AuthMiddleware("user"), controller.DeleteAddress)
	r.GET("/listaddress", middleware.AuthMiddleware("user"), controller.ListAddress)
	
	///profile
	r.GET("/profile", middleware.AuthMiddleware("user"), controller.ShowProfile)
	r.PATCH("/profile", middleware.AuthMiddleware("user"), controller.EditProfile)

	///forget password
	r.POST("/forgetpassword", controller.ForgetPassword)
	r.GET("/checkotp", controller.CheckOtp)
	r.PATCH("/newpassword", controller.NewPassword)

	///cart
	r.GET("/cart", middleware.AuthMiddleware("user"), controller.ViewCart)
	r.POST("/cart/:ID", middleware.AuthMiddleware("user"), controller.AddToCart)
	r.PATCH("/cart/:ID", middleware.AuthMiddleware("user"), controller.RemoveCart)
}
