package controller

import (
	"pkart/database"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddWishlist(c *gin.Context) {

	var Wishlist model.Wishlist

	userId := c.GetUint("userid")
	ProductID, _ := strconv.Atoi(c.Param("ID"))

	if err := database.DB.Where("Product_Id=? AND User_Id=?", uint(ProductID), userId).First(&Wishlist).Error; err == nil {
		c.JSON(409, gin.H{
			"Status":  "error",
			"Code":    409,
			"Error":   err,
			"Message": "product already exist in wishlist",
			"Data":    gin.H{},
		})
		return
	}
	Wishlist = model.Wishlist{
		UserId:    userId,
		ProductId: uint(ProductID),
	}
	if err := database.DB.Create(&Wishlist).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "error",
			"Code":    400,
			"Error":   "requested product is doesnot exists",
			"Message": "couldn't create the wishlist",
			"Data":    gin.H{},
		})
		return
	}
	
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "added to wishlist",
		"Data":    gin.H{},
	})
}

func RemoveWishlist(c *gin.Context) {

	var Wishlist model.Wishlist

	userId := c.GetUint("userid")
	ProductID, _ := strconv.Atoi(c.Param("ID"))

	if err := database.DB.Where("Product_Id=? AND User_Id=?", uint(ProductID), userId).First(&Wishlist).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "error",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "product not found in wishlist",
			"Data":    gin.H{},
		})
		return
	}

	if err := database.DB.Delete(&Wishlist).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "error",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "couldn't delete the wishlist",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "deleted from wishlist",
		"Data":    gin.H{
			"Product Id": ProductID,
		},
	})
}

func ShowWishlist(c *gin.Context) {

	var wishlist []model.Wishlist
	var show []gin.H
	userId := c.GetUint("userid")
	if err := database.DB.Preload("Product").Where("User_Id=?", userId).Find(&wishlist).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "error",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "no products found in wishlist",
			"Data":    gin.H{},
		})
		return
	}
	for _, v := range wishlist {
		
		show = append(show, gin.H{
			"Id":    v.Product.ID,
			"Name":  v.Product.Name,
			"Price": v.Product.Price,
			"Image1": v.Product.Image1,
		})
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "retrieved whishlist data",
		"Data": gin.H{
			"Products": show,
		},
	})
}

