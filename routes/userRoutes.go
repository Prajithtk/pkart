package routes

import (
	"pkart/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	r.POST("/signup", controller.UserSignUp)
	r.POST("/login", controller.UserLogin)
}
