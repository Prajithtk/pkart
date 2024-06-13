package controller

import (
	"fmt"
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func SearchProductAd(c *gin.Context) {

	var products []model.Products
	var show []gin.H

	searchQuery := c.Query("search")
	fmt.Println(searchQuery)

	database.DB.Where("name ILIKE ?", "%"+searchQuery+"%").Find(&products)
	if len(products) == 0 {
		c.JSON(404, gin.H{
			"Status":  "failed",
			"Code":    404,
			"Message": "products not found",
			"Data":    gin.H{},
		})
		return
	}
	for _, product := range products {
		var rate float32
		var r []model.Rating
		database.DB.Find(&r, "Product_Id=?", product.ID)
		for _, k := range r {
			rate += k.Rating
		}
		if len(r) == 0 {
			rate = 0
		} else {
			rate = rate / float32(len(r))
		}

		show = append(show, gin.H{
			"Id": product.ID,
			"Image1": product.Image1,
			"Name":   product.Name,
			"Price":  product.Offer,
			"Rating": rate,
		})
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "showing searched products",
		"Data":    show,
	})
}
