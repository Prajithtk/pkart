package controller

import (
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func ShowWallet(c *gin.Context) {

	// Logged := c.MustGet("Id").(uint)
	userId := c.GetUint("userid")
	var wallet model.Wallet

	if err := database.DB.First(&wallet, "User_Id=?", userId).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "error",
			"Code":    404,
			"Message": "wallet not found",
			"Data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "wallet found",
		"Data": gin.H{
			"Balance": wallet.Amount,
		},
	})

}
