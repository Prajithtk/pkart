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
		"status":  "true",
		"message": "product details are:",
		"values":  productinfo})
}
func AddProducts(c *gin.Context) {
	// var Product model.Products
	if err := c.ShouldBindJSON(&Product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}

	// Checking if a product with the same name already exists

	var existingProduct model.Products
	if result := database.DB.Where("name=?", Product.Name).First(&existingProduct); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "product already exists!!! try edit product"})
		return
	}
	// If no existing product found, proceed with adding the new product

	c.JSON(http.StatusSeeOther, gin.H{
		"success": "true",
		"message": "please upload images"})

}
func ProductImage(c *gin.Context) {
	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Failed to fetch images"})
		return
	}
	files := file.File["images"]
	var imagePaths []string

	for i, val := range files {
		filePath := "./images/" + strconv.Itoa(i) + "_" + val.Filename
		if err := c.SaveUploadedFile(val, filePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": "false",
				"message": "Failed to save images"})
			return
		}
		imagePaths = append(imagePaths, filePath)
	}
	Product.Image1 = imagePaths[0]
	Product.Image2 = imagePaths[1]
	Product.Image3 = imagePaths[2]

	if err := database.DB.Create(&Product).Error; err != nil {
		c.JSON(501, gin.H{
			"success": "false",
			"message": "Failed to add product to database"})
		return
	}

	c.JSON(200, gin.H{
		"success": "true",
		"message": "Product added successfully"})
	Product = model.Products{}

}

func EditProducts(c *gin.Context) {
	var productinfo model.Products
	if err := c.ShouldBindJSON(&productinfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "failed to bind json"})
		return
	}
	id := c.Param("ID")
	if err := database.DB.Where("id=?", id).Updates(&productinfo); err.Error != nil {
		c.JSON(404, gin.H{
			"success": "false",
			"message": "failed to edit product"})
		return
	}
	c.JSON(200, gin.H{
		"success": "true",
		"message": "successfully editted"})
}
func DeleteProducts(c *gin.Context) {
	var product model.Products
	id := c.Param("ID")
	err := database.DB.Where("id=?", id).Delete(&product)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "failed to delete product"})
		return
	}
	c.JSON(http.StatusSeeOther, gin.H{
		"success": "true",
		"message": "Product Deleted Successfully"})
}
