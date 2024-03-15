package routes

import (
	"pkart/controller"
	"pkart/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	r.POST("/signup", controller.UserSignUp)
	r.POST("/login", controller.UserLogin)
	r.POST("/otp", controller.OtpSignUp)
	r.POST("/resend-otp", controller.ResendOtp)
	r.GET("/products", middleware.AuthMiddleware("user"), controller.UserViewProducts)
}
