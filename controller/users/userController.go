package controller

import (
	"fmt"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "failed to bind json"})
		return
	}
	var ifUser model.Users
	uerr := database.DB.Where("email=?", UserInfo.Email).First(&ifUser)
	if uerr.Error == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "User Already exist"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(UserInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(401, gin.H{"success": "false", "message": "Failed to hash password"})
		return
	}
	cerr:= database.DB.Where("referal_code", UserInfo.ReferalCode).First(&CommisionUser)
	if cerr != nil {
		c.JSON(401, gin.H{"success": "true", "message": "Found commission user"})
		// return
	}
	UserInfo.Password = string(hashedPassword)
	UserInfo.Status = "Active"
	Code, _ := helper.GenerateRandomAlphanumericCode(6)
	// CommisionCode := UserInfo.ReferalCode
	UserInfo.ReferalCode = Code
	fmt.Println("commission code",CommisionUser)
	otp := onetp.GenerateOTP(6)
	err = onetp.SendOtp(UserInfo.Email, otp)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"err":    err.Error(),
			"error":  "failed to send otp",
			"code":   400,
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
	c.JSON(200, gin.H{"success": "true", "message": "OTP sent succesfully",
		"otp": otp})
}

func OtpSignUp(c *gin.Context) {
	var otp model.Otp
	err := c.BindJSON(&otp)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "failed to bind json"})
		return
	}
	var notp model.Otp
	oerr := database.DB.Where("Email=?", UserInfo.Email).Find(&notp)
	if oerr.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "failed to fetch OTP"})
		return
	}

	// currentTime := time.Now()
	// if currentTime.After(notp.Expires) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "OTP expired"})
	// 	return
	// }

	if otp.Otp != notp.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "invalid otp"})
		return
	}

	fmt.Println(UserInfo)
	create := database.DB.Create(&UserInfo)
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "message": "failed to create user"})
		return
		} else {
			c.JSON(200, gin.H{"success": "true", "message": "user created successfully"})
		}
	var userFetchData model.Users
	if err := database.DB.First(&userFetchData, "email=?", UserInfo.Email).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to fetch user details",
			"code":   400,
		})
		return
	}
	fmt.Println(userFetchData)
	var refFetch model.Users
	var wallet model.Wallet
	fmt.Println("commcode", CommisionUser)
	if UserInfo.ReferalCode != "" {

		if err := database.DB.First(&refFetch, "referal_code=?", CommisionUser.ReferalCode).Error; err != nil {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "failed to fetch user details for referal code",
				"code":   400,
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
		"status":  "Success",
		"message": "user created successfully",
	})
	UserInfo= model.Users{}
}

func ResendOtp(c *gin.Context) {
	var resend model.Otp
	err := c.ShouldBindJSON(&resend)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "failed to fetch otp"})
		return
	}
	var exotp model.Otp
	perr := database.DB.Where("email=?", resend.Email).First(&exotp)
	if perr.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "email not found"})
		return
	}

	newOtp := onetp.GenerateOTP(6)

	res := database.DB.Model(&model.Otp{}).Where("email=?", resend.Email).Updates(model.Otp{Otp: newOtp, Expires: time.Now().Add(1 * time.Minute)})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "message": "failed to resend otp"})
		return
	}
	onetp.SendOtp(resend.Email, newOtp)
	c.JSON(200, gin.H{"success": "true", "message": "otp resent successfully"})
}

func UserLogin(c *gin.Context) {
	var userDetails model.Users
	err := c.ShouldBindJSON(&userDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "failed to bind json"})
		return
	}
	var checkUser model.Users
	perror := database.DB.Where("email=?", userDetails.Email).First(&checkUser)
	if checkUser.Status == "blocked" {
		c.JSON(404, gin.H{"success": "false", "message": "Your account is blocked"})
	}
	if perror.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "User not found"})
		return
	}
	hashpassword := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(userDetails.Password))
	if hashpassword != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "Incorrect Password"})
		return
	}
	token, err := middleware.JwtToken(c, checkUser.ID, checkUser.Email, RoleUser)

	if err != nil {
		c.JSON(403, gin.H{"success": "false", "message": "failed to create token"})
		return
	}

	c.SetCookie("JwtTokenUser", token, int((time.Hour * 100).Seconds()), "/", "localhost", false, false)
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Successfully Logged in!",
		"Data": gin.H{
			"Token": token,
			"Id":    checkUser.ID,
		},
	})
}
func Logout(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER LOGGING OUT----------------------")

	c.SetCookie("JwtUser", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Logged out successfully!",
		"Data":    gin.H{},
	})
}
