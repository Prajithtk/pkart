package controller

import (
	"fmt"
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func SearchProduct(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SEARCH PRODUCT------------------------")

	var products []model.Products
	var show []gin.H

	searchQuery := c.Query("search")
	fmt.Println(searchQuery)

	database.DB.Where("name ILIKE ?", "%"+searchQuery+"%").Find(&products)
	if len(products) == 0 {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Message": "Products not found!",
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
				"Image1":  product.Image1,
				"Image2":  product.Image2,
				"Image3":  product.Image3,
				"Name":   product.Name,
				"Price":  product.Offer,
				"Rating": rate,
			})
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Showing searched products!",
		"Data":    show,
	})
}

func FilterProduct(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------FILTER PRODUCT------------------------")

	var Product []model.Products
	var show []gin.H
	// var img string
	var rate float32

	Category := c.Query("category")
	fmt.Println(Category)

	if err := database.DB.Preload("Category").Find(&Product).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Couldn't find any product!"})
		return
	}

	for _, v := range Product {
		if v.Category.Name == Category {
		
			if v.AvrgRating != 0 {
				rate = v.AvrgRating
			} else {
				rate = 0
			}
			show = append(show, gin.H{
				"Image1":  v.Image1,
				"Image2":  v.Image2,
				"Image3":  v.Image3,
				"Name":   v.Name,
				"Price":  v.Price,
				"Rating": rate,
			})
		}
	}
	if show == nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No products found in this category!",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Showing products of specific category!",
		"Data": gin.H{
			"Products": show,
		},
	})
}
