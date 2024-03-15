package controller

import (
	"fmt"
	"net/http"
	"pkart/database"
	"pkart/middleware"
	"pkart/model"
	onetp "pkart/onetp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var UserInfo model.Users

const RoleUser = "user"

func UserSignUp(c *gin.Context) {
	// var userInfo model.Users
	err := c.ShouldBindJSON(&UserInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}
	var ifUser model.Users
	uerr := database.DB.Where("email=?", UserInfo.Email).First(&ifUser)
	if uerr.Error == nil {
		c.JSON(http.StatusBadRequest, "User Already exist")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(UserInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(401, gin.H{"error": "Failed to hash password"})
		return
	}
	UserInfo.Password = string(hashedPassword)

	UserInfo.Status = "Active"

	otp := onetp.GenerateOTP(6)
	newOtp := model.Otp{
		Otp:     otp,
		Email:   UserInfo.Email,
		Expires: time.Now().Add(1 * time.Minute),
	}

	database.DB.Create(&newOtp)

	onetp.SendOtp(newOtp.Email, newOtp.Otp)
	c.JSON(200, gin.H{"message": "OTP sent succesfully",
		"otp": otp})
}

func OtpSignUp(c *gin.Context) {
	var otp model.Otp
	err := c.BindJSON(&otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}
	var notp model.Otp
	oerr := database.DB.Where("Email=?", UserInfo.Email).First(&notp)
	if oerr.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch OTP"})
		return
	}

	currentTime := time.Now()
	if currentTime.After(notp.Expires) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP expired"})
		return
	}

	if otp.Otp != notp.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid otp"})
		return
	}

	create := database.DB.Create(&UserInfo)
	fmt.Println(UserInfo)
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
	}
}

func ResendOtp(c *gin.Context) {
	var resend model.Otp
	err := c.ShouldBindJSON(&resend)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch otp"})
		return
	}
	var exotp model.Otp
	perr := database.DB.Where("email=?", resend.Email).First(&exotp)
	if perr.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email not found"})
		return
	}

	newOtp := onetp.GenerateOTP(6)

	res := database.DB.Model(&model.Otp{}).Where("email=?", resend.Email).Updates(model.Otp{Otp: newOtp, Expires: time.Now().Add(1 * time.Minute)})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resend otp"})
		return
	}
	onetp.SendOtp(resend.Email, newOtp)
	c.JSON(http.StatusOK, gin.H{"message": "otp resent successfully"})
}

func UserLogin(c *gin.Context) {
	var userDetails model.Users
	err := c.ShouldBindJSON(&userDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}

	var checkUser model.Users
	perror := database.DB.Where("email=?", userDetails.Email).First(&checkUser)
	if checkUser.Status == "blocked" {
		c.JSON(404, gin.H{"error": "Your account is blocked"})
	}
	if perror.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	hashpassword := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(userDetails.Password))
	if hashpassword != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Incorrect Password"})
		return
	}
	middleware.JwtToken(c, userDetails.ID, userDetails.Email, RoleUser)
	c.JSON(200, gin.H{"message": "successfully LoggedIn"})
}

func UserViewProducts(c *gin.Context) {
	var productList []model.Products
	database.DB.Preload("Category").Order("ID asc").Find(&productList)

	for _, val := range productList {
		c.JSON(200, gin.H{
			"id":          val.ID,
			"name":        val.Name,
			"color":       val.Color,
			"quantity":    val.Quantity,
			"description": val.Description,
			"category":    val.Category.Name,
			"status":      val.Status,
		})
	}
}
