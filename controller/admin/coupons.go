package controller

import (
	"pkart/database"
	"pkart/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Couponst struct {
	model.Coupons
	AddDay int `json:"day"`
}

func AddCoupon(c *gin.Context) {

	var addCoupon Couponst
	if err := c.ShouldBindJSON(&addCoupon); err != nil {
		c.JSON(404, gin.H{"error": "failed to bind json"})
		return
	}

	var coupon model.Coupons

	coupon.Name = addCoupon.Name
	coupon.Desc = addCoupon.Desc
	coupon.Code = addCoupon.Code
	coupon.Value = addCoupon.Value
	coupon.Min = addCoupon.Min
	coupon.Exp = time.Now().AddDate(0, 0, addCoupon.AddDay)

	err := database.DB.Create(&coupon)
	if err.Error != nil {
		c.JSON(409, gin.H{"message": "Coupon name or code already exist, please try to edit"})
	} else {
		c.JSON(200, gin.H{"message": "Coupon added successfully"})
	}
}

func ViewCoupon(c *gin.Context) {
	var couponlist []model.Coupons
	var couponinfo []gin.H

	database.DB.Order("ID asc").Find(&couponlist)
	// fmt.Println(coupons)
	for _, val := range couponlist {
		coupondetails := gin.H{
			"ID":          val.Id,
			"NAME":        val.Name,
			"DESCRIPTION": val.Desc,
			"CODE":        val.Code,
			"VALUE":       val.Value,
			"EXPIRY":      val.Exp,
		}
		couponinfo = append(couponinfo, coupondetails)
	}
	c.JSON(200, gin.H{
		"status":  "True",
		"message": "The coupon details are :",
		"values":  couponinfo})
}

func EditCoupon(c *gin.Context) {

	Id, _ := strconv.Atoi(c.Param("Id"))

	var edit model.Coupons
	var coupon model.Coupons

	c.BindJSON(&edit)

	database.DB.First(&coupon, Id)
	database.DB.Model(&model.Coupons{}).Where("Id=?", Id).Updates(edit)

	if coupon.Id == 0 {
		c.JSON(404, gin.H{"error": "Coupon not found."})
	} else {
		c.JSON(200, gin.H{"message": "Coupon edited succesfully."})
	}

}

func DeleteCoupon(c *gin.Context) {

	Id, _ := strconv.Atoi(c.Param("Id"))

	var coupon model.Coupons

	database.DB.First(&coupon, Id)

	if coupon.Id == 0 {
		c.JSON(404, gin.H{
			"success": "False",
			"message": "Coupon not found.",
			"data":    "{ }"})
	} else {
		database.DB.Delete(&coupon)
		c.JSON(200, gin.H{
			"success": "True",
			"message": "Coupon deleted successfully",
			"data":    "{ }"})
	}
}
