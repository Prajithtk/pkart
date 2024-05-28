package controller

import (
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func BestSelling(c *gin.Context) {
	var BestProduct []model.Products
	var BestList []gin.H
	query := c.Query("type")
	switch query {
	case "product":
		if err := database.DB.Table("order_items oi").Select("p.name, p.price , COUNT(oi.quantity) quantity").
			Joins("JOIN products p ON p.id = oi.product_id").
			Group("p.name, p.price").
			Order("quantity DESC").
			Limit(10).
			Scan(&BestProduct).Error; err != nil {
			c.JSON(500, gin.H{
				"status":  "Fail",
				"message": err.Error(),
				"code":    500,
			})
			return
		}
		for _, v := range BestProduct {
			BestList = append(BestList, gin.H{
				"productName": v.Name,
				"salesVolume": v.Quantity,
			})
		}

	case "category":
		var BestCategory []model.Category
		if err := database.DB.Table("order_items oi").
			Select("c.name, COUNT(oi.quantity) AS quantity").
			Joins("JOIN products p ON oi.product_id = p.id").Joins("JOIN categories c ON  c.id=p.category_id").
			Group("c.name").
			Order("quantity DESC").
			Limit(10).
			Scan(&BestCategory).Error; err != nil {
			c.JSON(500, gin.H{
				"status":  "Fail",
				"message": err,
				"code":    500,
			})
			return
		}
		for _, v := range BestCategory {
			BestList = append(BestList, gin.H{
				"categoryName": v.Name,
			})
		}
	}
	c.JSON(200, gin.H{
		"data":   BestList,
		"status": "Success",
	})
}
