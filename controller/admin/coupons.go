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
		c.JSON(404, gin.H{
			"Status":  "failed",
			"Code":    404,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
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
		c.JSON(409, gin.H{
			"Status":  "failed",
			"Code":    409,
			"Message": "coupon name or code already exist, please try to edit",
			"Data":    gin.H{},
		})
	} else {
		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": "coupon added successfully",
			"Data":    gin.H{},
		})
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
		"Status":  "success",
		"Code": 200,
		"Message": "the coupon details are :",
		"Data":  couponinfo})
}

func EditCoupon(c *gin.Context) {

	Id, _ := strconv.Atoi(c.Param("ID"))

	var edit model.Coupons
	var coupon model.Coupons

	c.BindJSON(&edit)

	database.DB.First(&coupon, Id)
	database.DB.Model(&model.Coupons{}).Where("Id=?", Id).Updates(edit)

	if coupon.Id == 0 {
		c.JSON(404, gin.H{
			"Status":  "failed",
			"Code":    404,
			"Message": "coupon not found",
			"Data":    gin.H{},
		})
	} else {
		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": "coupon edited successfully",
			"Data":    gin.H{},
		})
	}

}

func DeleteCoupon(c *gin.Context) {

	Id, _ := strconv.Atoi(c.Param("ID"))
	var coupon model.Coupons
	database.DB.First(&coupon, Id)

	if coupon.Id == 0 {
		c.JSON(404, gin.H{
			"Status": "failed",
			"Code": 404,
			"Message": "Coupon not found.",
			"Data":    gin.H{},
		})
	} else {
		database.DB.Delete(&coupon)
		c.JSON(200, gin.H{
			"Status": "success",
			"Code": 200,
			"Message": "coupon deleted successfully",
			"Data":    gin.H{},
		})
	}
}
