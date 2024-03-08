package controller

import (
	"net/http"
	"pkart/database"
	"pkart/helper"
	"pkart/model"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)
var UserInfo model.Users
func UserSignUp(c *gin.Context) {
	// var userInfo model.Users
	err := c.ShouldBindJSON(&UserInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, "**********failed to bind json**********")
		return
	}
	var ifUser model.Users
	uerr := database.DB.Where("email=?",UserInfo.Email).Find(&ifUser)
	if uerr !=nil{
		c.JSON(http.StatusBadRequest,"User Already exist")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(UserInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(401, "=-=-=-=-=-Failed to hash password-=-=-=-=-=")
		return
	}
	UserInfo.Password= string(hashedPassword)

	UserInfo.Status = "Active"

	otp:= helper.GenerateOtp()
	newOtp :=model.Otp{
		Otp: otp,
		Email: UserInfo.Email,
		Expires: time.Now().Add(1*time.Minute),
		
	}
	err:= database.DB.Create(&newOtp)
	if err!=nil {
		c.JSON(401,"failed to send otp")
	}
	
	// if createUser.Error != nil {
	// 	c.JSON(401, "Email already exists")
	// } else {
	// 	c.JSON(http.StatusSeeOther, "Account created successfully")
	// }
}

func UserLogin(c *gin.Context) {
	var userDetails model.Users
	err := c.ShouldBindJSON(&userDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, "**********failed to bind json**********")
		return
	}
	
	var checkUser model.Users
	perror := database.DB.Where("email=?", userDetails.Email).First(&checkUser)
	if checkUser.Status == "blocked"{
		c.JSON(404,"Your account is blocked ")
	}
	if perror.Error != nil {
		c.JSON(http.StatusInternalServerError, "User not found")
		return
	}
	hashpassword := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(userDetails.Password))
	if hashpassword != nil {
		c.JSON(http.StatusInternalServerError, "Incorrect Password")
		return
	}
	c.JSON(200, "successfully LoggedIn")
}

