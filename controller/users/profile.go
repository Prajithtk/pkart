package controller

import (
	"fmt"
	"pkart/database"
	"pkart/model"
	onetp "pkart/onetp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowProfile(c *gin.Context) {
	userId := c.GetUint("userid")
	var userProfile model.Users
	// var userInfo []gin.H
	if err := database.DB.Preload("Address").Where("id=?", userId).First(&userProfile).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch user profile"})
		return
	}
	response := gin.H{
		"Name":   userProfile.Name,
		"Email":  userProfile.Email,
		"Phone":  userProfile.Phone,
		"Gender": userProfile.Gender,
	}
	// if len(userProfile.Address) > 0 {
	// 	response["Address"] = userProfile.Address
	// }
	c.JSON(200, gin.H{
		"success": "true",
		"message": "profile details are:",
		"values":  response})
}

func EditProfile(c *gin.Context) {
	var editdata model.Users
	fmt.Println("hello world")
	if err := c.ShouldBindJSON(&editdata); err != nil {
		c.JSON(400, gin.H{"error": "failed to bind json"})
		return
	}
	fmt.Println(editdata)
	fmt.Println("editdata")
	userId := c.GetUint("userid")
	fmt.Println(userId)
	if err := database.DB.Where("id=?", userId).Updates(&editdata); err.Error != nil {
		c.JSON(404, gin.H{"error": "failed to update address"})
		return
	}
	c.JSON(200, gin.H{"message": "profile editted successfully"})
}

var Useremail model.Users

func ForgetPassword(c *gin.Context) {
	var newOtp model.Otp
	var CheckOtp model.Otp
	var checkMail model.Users
	if err := c.ShouldBindJSON(&Useremail); err != nil {
		c.JSON(400, gin.H{"error": "failed to bind json"})
		return
	}
	if Useremail.Email == "" {
		c.JSON(303, gin.H{"error": "enter the email"})
		return
	}
	if err := database.DB.Where("email=?", Useremail.Email).First(&checkMail).Error; err != nil {
		c.JSON(400, gin.H{"error": "no account found,SignUp"})
		return
	}
	otp := onetp.GenerateOTP(6)
	newOtp = model.Otp{
		Otp:     otp,
		Email:   Useremail.Email,
		Expires: time.Now().Add(3 * time.Minute),
	}
	fmt.Println(newOtp)
	if err := onetp.SendOtp(Useremail.Email, otp); err != nil {
		c.JSON(500, gin.H{"error": "failed to send otp"})
		return
	}
	if res := database.DB.Where("email=?", Useremail.Email).First(&CheckOtp).Error; res != nil {
		if err := database.DB.Create(&newOtp).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to save otp"})
			return
		}
	} else {
		if err := database.DB.Model(&CheckOtp).Where("email=?", Useremail.Email).Updates(&newOtp).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to update otp table"})
			return
		}
	}
	c.JSON(303, gin.H{"message": "OTP sent to successfully", "otp": otp})
}

func CheckOtp(c *gin.Context) {
	var otp model.Otp
	if err := c.ShouldBindJSON(&otp); err != nil {
		c.JSON(404, gin.H{"error": "failed to bind json"})
		return
	}
	var notp model.Otp
	oerr := database.DB.Where("email=?", Useremail.Email).First(&notp)
	fmt.Println(Useremail.Email)
	if oerr.Error != nil {
		c.JSON(401, gin.H{"error": "failed to fetch OTP"})
		return
	}

	currentTime := time.Now()
	if currentTime.After(notp.Expires) {
		c.JSON(401, gin.H{"error": "OTP expired"})
		return
	}

	if otp.Otp != notp.Otp {
		c.JSON(400, gin.H{"error": "invalid otp"})
		return
	}
	c.JSON(200, gin.H{"message": "Enter new password"})
}

func NewPassword(c *gin.Context) {
	var newPass model.Users
	if err := c.ShouldBindJSON(&newPass); err != nil {
		c.JSON(404, gin.H{"error": "failed to bind json"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(401, gin.H{"error": "Failed to hash password"})
		return
	}
	newPass.Password = string(hashedPassword)
	if err := database.DB.Model(&newPass).Where("email=?", Useremail.Email).Updates(&newPass).Error; err != nil {
		c.JSON(404, gin.H{"error": "failed to update password"})
	}
	c.JSON(200, gin.H{"message": "successfully updated password"})
}
