package controller

import (
	"fmt"
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func ShowWallet(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------WALLET SHOWING----------------------")

	// Logged := c.MustGet("Id").(uint)
	userId := c.GetUint("userid")


	var wallet model.Wallet

	if err := database.DB.First(&wallet, "User_Id=?", userId).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Wallet not found!",
			"Data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Wallet found!",
		"Data":    gin.H{"Balance": wallet.Amount},
	})
}
