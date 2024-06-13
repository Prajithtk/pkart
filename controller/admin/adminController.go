package controller

import (
	"fmt"
	"os"
	"pkart/database"
	"pkart/middleware"
	"pkart/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminPage(c *gin.Context) {
	var totalSales []model.Orders
	var totalAmount float64
	var totalOrder int
	if err := database.DB.Find(&totalSales).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": err.Error(),
		})
	}
	for _, v := range totalSales {
		totalAmount += float64(v.Amount)
		totalOrder += 1
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "welcome admin page",
		"Data":    gin.H{"sales": totalAmount, "orderCount": totalOrder},
	})
}

var Product model.Products
var Ida, _ = strconv.Atoi(os.Getenv("ID"))
var Email = os.Getenv("ADMIN")

const RoleAdmin = "Admin"

func AdminLogin(c *gin.Context) {
	var admin model.Admin
	err := c.ShouldBindJSON(&admin)
	if err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
		return
	}
	email := os.Getenv("ADMIN")
	password := os.Getenv("ADMIN_PASSWORD")

	if email != admin.Name || password != admin.Password {
		c.JSON(404, gin.H{
			"Status":  "failed",
			"Code":    404,
			"Message": "Incorrect username or password",
			"Data":    gin.H{},
		})
		return
	} else {
		token, err := middleware.JwtToken(c, uint(Ida), Email, RoleAdmin)
		if err != nil {
			c.JSON(403, gin.H{
				"Status":  "failed",
				"Code":    403,
				"Message": "failed to create token",
				"Data":    gin.H{},
			})
			return
		}
		c.SetCookie("JwtTokenAdmin", token, int((time.Hour * 100).Seconds()), "/", "pkartz.shop", false, false)
		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": "successfully Logged in!",
			"Data": gin.H{
				"Token": token,
			},
		})
	}
}

func Logout(c *gin.Context) {

	c.SetCookie("JwtTokenAdmin", "", -1, "/", "pkartz.shop", false, true)
	c.JSON(200, gin.H{
		"Status":  "Success",
		"Code":    200,
		"Message": "Logged out successfully!",
		"Data":    gin.H{},
	})
}

//------------------------user management---------------------------//
//------------------------------------------------------------------//

func ListUsers(c *gin.Context) {
	var usersList []model.Users
	var userInfo []gin.H

	database.DB.Order("ID asc").Find(&usersList)

	for _, val := range usersList {
		userDetails := gin.H{
			"id":     val.ID,
			"name":   val.Name,
			"email":  val.Email,
			"phone":  val.Phone,
			"status": val.Status,
		}
		userInfo = append(userInfo, userDetails)
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "the user details are:",
		"Data": gin.H{
			"Values": userInfo,
		},
	})
}
func BlockUser(c *gin.Context) {
	var user model.Users
	id := c.Param("ID")
	database.DB.First(&user, id)
	if user.Status == "Blocked" {
		database.DB.Model(&user).Update("status", "Active")
		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": fmt.Sprintf("User with ID %s has been unblocked", id),
			"Data":    gin.H{},
		})
	} else {
		database.DB.Model(&user).Update("status", "Blocked")
		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": fmt.Sprintf("User with ID %s has been blocked", id),
			"Data":    gin.H{},
		})
	}
}
