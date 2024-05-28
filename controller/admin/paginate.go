package controller

import (
	"net/http"
	"pkart/database"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Paginate(c *gin.Context) {
	var products model.Products
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	result := database.DB.Limit(pageSize).Offset(offset).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var totalRecords int64
	database.DB.Model(&model.Products{}).Count(&totalRecords)

	c.JSON(http.StatusOK, gin.H{
		"data":        products,
		"page":        page,
		"page_size":   pageSize,
		"total_items": totalRecords,
		"total_pages": (totalRecords + int64(pageSize) - 1) / int64(pageSize),
	})
}
