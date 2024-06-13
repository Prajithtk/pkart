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
		c.JSON(500, gin.H{
			"Status":  "failed",
			"Code":    500,
			"Message": "failed to fetch user profile",
			"Data":    gin.H{},
		})
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
		"Status": "success",
		"Code" : 200,
		"Message": "profile details are:",
		"Data":  response})
}

func EditProfile(c *gin.Context) {
	var editdata model.Users
	if err := c.ShouldBindJSON(&editdata); err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
		return
	}
	fmt.Println(editdata)
	fmt.Println("editdata")
	userId := c.GetUint("userid")
	fmt.Println(userId)
	if err := database.DB.Where("id=?", userId).Updates(&editdata); err.Error != nil {
		c.JSON(404, gin.H{
			"Status":  "failed",
			"Code":    404,
			"Message": "failed to update address",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "profile edited successfully",
		"Data":    gin.H{},
	})
}

var Useremail model.Users

func ForgetPassword(c *gin.Context) {
	var newOtp model.Otp
	var CheckOtp model.Otp
	var checkMail model.Users
	if err := c.ShouldBindJSON(&Useremail); err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
		return
	}
	if Useremail.Email == "" {
		c.JSON(303, gin.H{
			"Status":  "error",
			"Code":    303,
			"Message": "enter the email",
			"Data":    gin.H{},
		})
		return
	}
	if err := database.DB.Where("email=?", Useremail.Email).First(&checkMail).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "no account found, SignUp",
			"Data":    gin.H{},
		})
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
		c.JSON(500, gin.H{
			"Status":  "failed",
			"Code":    500,
			"Message": "failed to sent OTP",
			"Data":    gin.H{},
		})
		return
	}
	if res := database.DB.Where("email=?", Useremail.Email).First(&CheckOtp).Error; res != nil {
		if err := database.DB.Create(&newOtp).Error; err != nil {
			c.JSON(500, gin.H{
				"Status":  "failed",
				"Code":    500,
				"Message": "failed to save OTP",
				"Data":    gin.H{},
			})
			return
		}
	} else {
		if err := database.DB.Model(&CheckOtp).Where("email=?", Useremail.Email).Updates(&newOtp).Error; err != nil {
			c.JSON(500, gin.H{
				"Status":  "failed",
				"Code":    500,
				"Message": "failed to update OTP table",
				"Data":    gin.H{},
			})
			return
		}
	}
	c.JSON(303, gin.H{
		"Status":  "success",
		"Code":    303,
		"Message": "OTP sent successfully",
		"Data":    gin.H{
			"OTP": otp,
		},
	})
}

func CheckOtp(c *gin.Context) {
	var otp model.Otp
	if err := c.ShouldBindJSON(&otp); err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
		return
	}
	var notp model.Otp
	oerr := database.DB.Where("email=?", Useremail.Email).First(&notp)
	fmt.Println(Useremail.Email)
	if oerr.Error != nil {
		c.JSON(401, gin.H{
			"Status":  "failed",
			"Code":    401,
			"Message": "failed to fetch OTP",
			"Data":    gin.H{},
		})
		return
	}
	currentTime := time.Now()
	if currentTime.After(notp.Expires) {
		c.JSON(401, gin.H{
			"Status":  "failed",
			"Code":    401,
			"Message": "OTP expired",
			"Data":    gin.H{},
		})
		return
	}
	if otp.Otp != notp.Otp {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "invalid otp",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "enter new password",
		"Data":    gin.H{},
	})
}

func NewPassword(c *gin.Context) {
	var newPass model.Users
	if err := c.ShouldBindJSON(&newPass); err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(401, gin.H{
			"Status":  "failed",
			"Code":    401,
			"Message": "failed to hash password",
			"Data":    gin.H{},
		})
		return
	}
	newPass.Password = string(hashedPassword)
	if err := database.DB.Model(&newPass).Where("email=?", Useremail.Email).Updates(&newPass).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "failed",
			"Code":    404,
			"Message": "failed to update password",
			"Data":    gin.H{},
		})
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "successfully updated password",
		"Data":    gin.H{},
	})
}
