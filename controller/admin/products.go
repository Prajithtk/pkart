package controller

import (
	"net/http"
	"pkart/database"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ViewProducts(c *gin.Context) {
	var productList []model.Products
	var productinfo []gin.H
	database.DB.Preload("Category").Order("ID asc").Find(&productList)

	for _, val := range productList {
		productdetails := gin.H{
			"id":          val.ID,
			"name":        val.Name,
			"color":       val.Color,
			"quantity":    val.Quantity,
			"price":	val.Price,
			"offer":	val.Offer,
			"description": val.Description,
			"category":    val.Category.Name,
			"status":      val.Status,
			"images1":     val.Image1,
			"images2":     val.Image2,
			"images3":     val.Image3,
		}
		productinfo = append(productinfo, productdetails)
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code": 200,
		"Message": "product details are:",
		"Data":  productinfo})
}
func AddProducts(c *gin.Context) {
	// var Product model.Products
	if err := c.ShouldBindJSON(&Product); err != nil {
		c.JSON(400, gin.H{
			"Status": "failed",
			"Code": 400,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})		
	return
	}

	// Checking if a product with the same name already exists

	var existingProduct model.Products
	if result := database.DB.Where("name=?", Product.Name).First(&existingProduct); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "failed",
			"Code": 400,
			"Message": "product already exists, try edit product",
			"Data": gin.H{},
		})
		return
	}
	// If no existing product found, proceed with adding the new product

	c.JSON(http.StatusSeeOther, gin.H{
		"status": "success",
		"Code": 303,
		"message": "please upload images",
		"Data": gin.H{},
	})

}
func ProductImage(c *gin.Context) {
	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "failed",
			"Code": 400,
			"Message": "Failed to fetch images",
			"Data": gin.H{},
		})
		return
	}
	files := file.File["images"]
	var imagePaths []string

	for i, val := range files {
		filePath := "./images/" + strconv.Itoa(i) + "_" + val.Filename
		if err := c.SaveUploadedFile(val, filePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status": "failed",
				"Code": 400,
				"Message": "failed to save images",
				"Data": gin.H{},
			})
			return
		}
		imagePaths = append(imagePaths, filePath)
	}
	Product.Image1 = imagePaths[0]
	Product.Image2 = imagePaths[1]
	Product.Image3 = imagePaths[2]

	if err := database.DB.Create(&Product).Error; err != nil {
		c.JSON(501, gin.H{
			"Status": "failed",
			"Code": 501,
			"Message": "Failed to add product to database",
		"Data": gin.H{},
	})
		return
	}

	c.JSON(200, gin.H{
		"Status": "success",
		"Code": 200,
		"Message": "product added successfully",
		"Data": gin.H{},
	})
	Product = model.Products{}

}

func EditProducts(c *gin.Context) {
	var productinfo model.Products
	if err := c.ShouldBindJSON(&productinfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "failed",
			"Code": 400,
			"message": "failed to bind json",
			"Data" : gin.H{},
		})
		return
	}
	id := c.Param("ID")
	if err := database.DB.Where("id=?", id).Updates(&productinfo); err.Error != nil {
		c.JSON(404, gin.H{
			"Status": "success",
			"Code" : 404,
			"Message": "failed to edit product",
			"Data": gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status": "success",
		"Code": 200,
		"Message": "successfully edited the product",
		"Data": gin.H{},
	})
}
func DeleteProducts(c *gin.Context) {
	var product model.Products
	id := c.Param("ID")
	err := database.DB.Where("id=?", id).Delete(&product)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status": "failed",
			"Code": 500,
			"Message": "failed to delete product",
			"Data": gin.H{},
		})
		return
	}
	c.JSON(http.StatusSeeOther, gin.H{
		"Status": "success",
		"Code": 303,
		"message": "Product Deleted Successfully",
		"Data": gin.H{},
	})
}
