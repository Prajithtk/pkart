package controller

import (
	"net/http"
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func ViewCategory(c *gin.Context) {
	var categoryList []model.Category
	database.DB.Order("ID asc").Find(&categoryList)

	for _, val := range categoryList {
		c.JSON(200, gin.H{
			"id":          val.ID,
			"name":        val.Name,
			"description": val.Description,
			"status":      val.Status,
		})
	}
}
func AddCategory(c *gin.Context) {
	var categoryinfo model.Category
	err := c.ShouldBindJSON(&categoryinfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}
	addCategory := database.DB.Create(&model.Category{
		Name:        categoryinfo.Name,
		Description: categoryinfo.Description,
		Status:      categoryinfo.Status,
	})
	if addCategory.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to add category"})
	} else {
		c.JSON(200, gin.H{"message": "Category added successfully"})
	}
}
func EditCategory(c *gin.Context) {
	var categoryinfo model.Category
	err := c.ShouldBindJSON(&categoryinfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
	}
	id := c.Param("ID")
	cerr := database.DB.Where("id=?", id).Updates(&categoryinfo)
	if cerr.Error != nil {
		c.JSON(404, gin.H{"error": "failed to edit category"})
	}
	c.JSON(200, gin.H{"message": "successfully editted"})
}
func BlockCategory(c *gin.Context) {
	var category model.Category
	id := c.Param("ID")
	database.DB.First(&category, id)
	if category.Status == "blocked" {
		database.DB.Model(&category).Update("status", "active")
		c.JSON(http.StatusOK, gin.H{"message": "Category Active"})
	} else {
		database.DB.Model(&category).Update("status", "blocked")
		c.JSON(http.StatusOK, gin.H{"message": "Category Blocked"})
	}
}
func DeleteCategory(c *gin.Context) {
	var category model.Category
	id := c.Param("ID")
	err := database.DB.Where("id=?", id).Delete(&category)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}
	c.JSON(200, gin.H{"message": "Category Deleted Successfully"})
}
