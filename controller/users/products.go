package controller

import (
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func UserViewProducts(c *gin.Context) {
	var productList []model.Products
	var list []gin.H
	database.DB.Preload("Category").Order("ID asc").Find(&productList)
	for _, val := range productList {
	list = 	append(list, gin.H{
			"id":          val.ID,
			"name":        val.Name,
			"color":       val.Color,
			"quantity":    val.Quantity,
			"price":       val.Price,
			"offer":	val.Offer,
			"description": val.Description,
			// "categiryid":  val.CategoryId,
			"category": val.Category.Name,
			"status":   val.Status,
			"image1":   val.Image1,
			"image2":   val.Image2,
			"image3":   val.Image3,
		})
	}
	c.JSON(201, gin.H{
		"Status":  "Success",
		"Code" : 201,
		"Message": "products are:",
		"Data":    list,
	})
}
