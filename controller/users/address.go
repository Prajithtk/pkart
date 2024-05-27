package controller

import (
	"fmt"
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {
	userId := c.GetUint("userid")

	var address model.Address

	if err := c.BindJSON(&address); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if len(address.PinCode) != 6 {
		c.JSON(400, gin.H{"error": "pincode should be 6 digits"})
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
		c.JSON(400, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, "message: successfully added an address")
}

func EditAddress(c *gin.Context) {
	var addressDetails model.Address
	err := c.ShouldBindJSON(&addressDetails)
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to bind json"})
		return
	}
	userId := c.GetUint("userid")
	addId := c.Param("ID")
	aerr := database.DB.Where("address_id=? AND user_id=?", addId, userId).Updates(&addressDetails)
	if aerr.Error != nil {
		c.JSON(400, gin.H{"error": "failed to edit address"})
		return
	}
	c.JSON(200, gin.H{"message": "successfully edited the address"})
}

func DeleteAddress(c *gin.Context) {
	var address model.Address
	adrId := c.Param("ID")
	userId := c.GetUint("userid")
	if err := database.DB.Where("address_id=? AND user_id=?", adrId, userId).First(&address).Error; err != nil {
		c.JSON(400, gin.H{"error": "failed to fetch address"})
		fmt.Println(err)
		return
	}
	res := database.DB.Delete(&address)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": "failed to delete address"})
		return
	}
	c.JSON(200, gin.H{"message": "successfully deleted address"})
}

func ListAddress(c *gin.Context) {
	userId := c.GetUint("userid")
	var addressList []model.Address
	database.DB.Order("address_id asc").Find(&addressList).Where("userid=?", userId)
	for _, val := range addressList {
		c.JSON(200, gin.H{
			"id":            val.AddressId,
			"building name": val.BuildingName,
			"street":        val.Street,
			"city":          val.City,
			"state":         val.State,
			"landmark":      val.Landmark,
			"pincode":       val.PinCode,
		})
	}
}
