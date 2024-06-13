package controller

import (
	"fmt"
	"pkart/database"
	"pkart/helper"
	"pkart/middleware"
	"pkart/model"
	onetp "pkart/onetp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var UserInfo model.Users
var CommisionUser model.Users
var RoleUser = "User"

func UserSignUp(c *gin.Context) {
	// var userInfo model.Users
	err := c.ShouldBindJSON(&UserInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to bind JSON",
			"Data":    gin.H{},
		})
		return
	}
	var ifUser model.Users
	uerr := database.DB.Where("email=?", UserInfo.Email).First(&ifUser)
	if uerr.Error == nil {
		c.JSON(401, gin.H{
			"Status":  "failed",
			"Code":    401,
			"Message": "user already exist",
			"Data":    gin.H{},
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(UserInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(401, gin.H{
			"Status":  "failed",
			"Code":    401,
			"Message": "failed to hash password",
			"Data":    gin.H{},
		})
		return
	}
	cerr := database.DB.Where("referal_code", UserInfo.ReferalCode).First(&CommisionUser)
	if cerr != nil {
		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": "found commission user",
			"Data":    gin.H{},
		})
		// return
	}
	UserInfo.Password = string(hashedPassword)
	UserInfo.Status = "Active"
	Code, _ := helper.GenerateRandomAlphanumericCode(6)
	// CommisionCode := UserInfo.ReferalCode
	UserInfo.ReferalCode = Code
	fmt.Println("commission code", CommisionUser)
	otp := onetp.GenerateOTP(6)
	err = onetp.SendOtp(UserInfo.Email, otp)
	if err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to send otp",
			"Error":   err.Error(),
			"Data":    gin.H{},
		})
		return
	}
	// var otpStore model.Otp
	// result := database.DB.First(&otpStore, "email=?", UserInfo.Email)
	// if result.Error != nil {
	// 	otpStore = model.Otp{
	// 		Otp:       otp,
	// 		Email:     UserInfo.Email,
	// 		// CreatedAt: time.Now(),
	// 		Expires:  time.Now().Add(180 * time.Second),
	// 	}
	// 	err := database.DB.Create(&otpStore)
	// 	if err.Error != nil {
	// 		c.JSON(400, gin.H{
	// 			"status": "Fail",
	// 			"error":  "failed to save otp details",
	// 			"code":   400,
	// 		})
	// 		return
	// 	}
	// } else {
	// 	err := database.DB.Model(&otpStore).Where("email=?", UserInfo.Email).Updates(model.Otp{
	// 		Otp:      otp,
	// 		Expires: time.Now().Add(180 * time.Second),
	// 	})
	// 	if err.Error != nil {
	// 		c.JSON(400, gin.H{
	// 			"status": "Fail",
	// 			"error":  "Failed to update OTP Details",
	// 			"code":   400,
	// 		})
	// 		return
	// 	}
	// }

	newOtp := model.Otp{
		Otp:     otp,
		Email:   UserInfo.Email,
		Expires: time.Now().Add(1 * time.Minute),
	}

	database.DB.Create(&newOtp)

	onetp.SendOtp(newOtp.Email, newOtp.Otp)
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "OTP sent successfully",
		"Data": gin.H{
			"OTP": otp,
		},
	})
}

func OtpSignUp(c *gin.Context) {
	var otp model.Otp
	err := c.BindJSON(&otp)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to bind JSON",
			"Data":    gin.H{},
		})
		return
	}
	var notp model.Otp
	oerr := database.DB.Where("Email=?", UserInfo.Email).Find(&notp)
	if oerr.Error != nil {
		c.JSON(401, gin.H{
			"Status":  "failed",
			"Code":    401,
			"Message": "failed to fetch OTP",
			"Data":    gin.H{},
		})
		return
	}

	// currentTime := time.Now()
	// if currentTime.After(notp.Expires) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "OTP expired"})
	// 	return
	// }

	if otp.Otp != notp.Otp {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "invalid OTP",
			"Data":    gin.H{},
		})
		return
	}

	// fmt.Println(UserInfo)
	create := database.DB.Create(&UserInfo)
	if create.Error != nil {
		c.JSON(500, gin.H{
			"Status":  "failed",
			"Code":    500,
			"Message": "failed to create user",
			"Data":    gin.H{},
		})
		return
	} else {
		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": "user created successfully",
			"Data":    gin.H{},
		})
	}
	var userFetchData model.Users
	if err := database.DB.First(&userFetchData, "email=?", UserInfo.Email).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to fetch user details",
			"Data":    gin.H{},
		})
		return
	}
	// fmt.Println(userFetchData)
	var refFetch model.Users
	var wallet model.Wallet
	// fmt.Println("commcode", CommisionUser)
	if UserInfo.ReferalCode != "" {

		if err := database.DB.First(&refFetch, "referal_code=?", CommisionUser.ReferalCode).Error; err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"Code":    400,
				"Message": "failed to fetch user details for referal code",
				"Data":    gin.H{},
			})
			return
		}
		wallet.Amount = 50
		wallet.UserId = userFetchData.ID
		fmt.Println(refFetch)
		var refwallet model.Wallet
		database.DB.First(&refwallet, "user_id=?", refFetch.ID)
		refwallet.Amount += 50
		refwallet.UserId = refFetch.ID
		database.DB.Save(&refwallet)
	}
	wallet.UserId = userFetchData.ID
	database.DB.Create(&wallet)

	c.SetCookie("sessionId", "", -1, "/", "", false, false)
	c.JSON(201, gin.H{
		"Status":  "success",
		"Code":    201,
		"message": "user created successfully",
		"Data":    gin.H{},
	})
	UserInfo = model.Users{}
}

func ResendOtp(c *gin.Context) {
	var resend model.Otp
	err := c.ShouldBindJSON(&resend)
	if err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to fetch OTP",
			"Data":    gin.H{},
		})
		return
	}
	var exotp model.Otp
	perr := database.DB.Where("email=?", resend.Email).First(&exotp)
	if perr.Error != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "email not found",
			"Data":    gin.H{},
		})
		return
	}

	newOtp := onetp.GenerateOTP(6)

	res := database.DB.Model(&model.Otp{}).Where("email=?", resend.Email).Updates(model.Otp{Otp: newOtp, Expires: time.Now().Add(1 * time.Minute)})
	if res.Error != nil {
		c.JSON(500, gin.H{
			"Status":  "failed",
			"Code":    500,
			"Message": "failed to resend OTP",
			"Data":    gin.H{},
		})
		return
	}
	onetp.SendOtp(resend.Email, newOtp)
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "OTP resend successfully",
		"Data":    gin.H{},
	})
}



func UserLogin(c *gin.Context) {
	var userDetails model.Users
	err := c.ShouldBindJSON(&userDetails)
	if err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to bind JSON",
			"Data":    gin.H{},
		})
		return
	}
	var checkUser model.Users
	perror := database.DB.Where("email=?", userDetails.Email).First(&checkUser)
	if checkUser.Status == "Blocked" {
		c.JSON(404, gin.H{
			"Status":  "failed",
			"Code":    404,
			"Message": "your account is blocked",
			"Data":    gin.H{},
		})
	}
	if perror.Error != nil {
		c.JSON(401, gin.H{
			"Status":  "failed",
			"Code":    401,
			"Message": "user not found",
			"Data":    gin.H{},
		})
		return
	}
	hashpassword := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(userDetails.Password))
	if hashpassword != nil {
		c.JSON(401, gin.H{
			"Status":  "failed",
			"Code":    401,
			"Message": "incorrect password",
			"Data":    gin.H{},
		})
		return
	}
	token, err := middleware.JwtToken(c, checkUser.ID, checkUser.Email, RoleUser)

	if err != nil {
		c.JSON(403, gin.H{
			"Status":  "failed",
			"Code":    403,
			"Message": "failed to create token",
			"Data":    gin.H{},
		})
		return
	}

	c.SetCookie("JwtTokenUser", token, int((time.Hour * 100).Seconds()), "/", "pkartz.shop", false, false)
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "successfully logged in",
		"Data": gin.H{
			"Token": token,
			"Id":    checkUser.ID,
		},
	})
}
func Logout(c *gin.Context) {

	c.SetCookie("JwtUser", "", -1, "/", "pkartz.shop", false, true)
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Logged out successfully!",
		"Data":    gin.H{},
	})
}
