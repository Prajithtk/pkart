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
		c.JSON(400, gin.H{
			"Status": "error",
			"Code": 400,
			"Message": "invalid page number",
			"Data":    gin.H{},
		})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(400, gin.H{
			"Status": "error",
			"Code": 400,
			"Message": "invalid page size",
			"Data":    gin.H{},
		})
		return
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Retrieve paginated records
	var products []model.Products
	var productinfo []gin.H
	result := database.DB.Preload("Category").Limit(pageSize).Offset(offset).Find(&products)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Status": "failed",
			"Code": 500,
			"Message": "failed to fetch records",
			"Data":    gin.H{},
		})
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
		c.JSON(500, gin.H{
			"Status": "failed",
			"Code": 500,
			"Message": "failed to retreive total count",
			"Data":    gin.H{},
		})
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
		"Status":  "success",
		"Code": 200,
		"Message": "product details are:",
		"Data":  response})
}
