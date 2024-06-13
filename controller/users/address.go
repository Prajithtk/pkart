package controller

import (
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {
	userId := c.GetUint("userid")

	var address model.Address

	if err := c.BindJSON(&address); err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
		return
	}
	if len(address.PinCode) != 6 {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "pincode should be 6 digits",
			"Data":    gin.H{},
		})
		return
	}
	err := database.DB.Create(&model.Address{
		BuildingName: address.BuildingName,
		Street:       address.Street,
		City:         address.City,
		State:        address.State,
		Landmark:     address.Landmark,
		PinCode:      address.PinCode,
		UserId:       userId,
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": err.Error(),
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "successfully added an address",
		"Data":    gin.H{},
	})
}

func EditAddress(c *gin.Context) {
	var addressDetails model.Address
	err := c.ShouldBindJSON(&addressDetails)
	if err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
		return
	}
	userId := c.GetUint("userid")
	addId := c.Param("ID")
	aerr := database.DB.Where("address_id=? AND user_id=?", addId, userId).Updates(&addressDetails)
	if aerr.Error != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to edit address",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "successfully edited the address",
		"Data":    gin.H{},
	})
}

func DeleteAddress(c *gin.Context) {
	var address model.Address
	adrId := c.Param("ID")
	userId := c.GetUint("userid")
	if err := database.DB.Where("address_id=? AND user_id=?", adrId, userId).First(&address).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to fetch address",
			"Data":    gin.H{},
		})
		return
	}
	res := database.DB.Delete(&address)
	if res.Error != nil {
		c.JSON(400, gin.H{
			"Status":  "failed",
			"Code":    400,
			"Message": "failed to delete address",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "successfully deleted address",
		"Data":    gin.H{},
	})
}

func ListAddress(c *gin.Context) {
	userId := c.GetUint("userid")
	var addressList []model.Address
	var list []gin.H
	database.DB.Order("address_id asc").Find(&addressList).Where("userid=?", userId)
	for _, val := range addressList {
		list = append(list, gin.H{
			"id":            val.AddressId,
			"building name": val.BuildingName,
			"street":        val.Street,
			"city":          val.City,
			"state":         val.State,
			"landmark":      val.Landmark,
			"pincode":       val.PinCode,
		})
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "address details are : ",
		"Data": gin.H{
			"categories": list,
		},
	})
}
