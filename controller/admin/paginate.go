package controller

import (
	"net/http"
	"pkart/database"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PaginateProducts(c *gin.Context) {
	// Parse query parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Retrieve paginated records
	var products []model.Products
	var productinfo []gin.H
	result := database.DB.Preload("Category").Limit(pageSize).Offset(offset).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}
	for _, val := range products {
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

	// Retrieve total count of records (for pagination metadata)
	var totalCount int64
	if err := database.DB.Model(&model.Products{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve total count"})
		return
	}

	// Calculate total number of pages
	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	// Prepare response
	response := gin.H{
		"data":        productinfo,
		"page":        page,
		"page_size":   pageSize,
		"total_items": totalCount,
		"total_pages": totalPages,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "true",
		"message": "product details are:",
		"values":  response})
}
