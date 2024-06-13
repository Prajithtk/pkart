package controller

import (
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func ViewCategory(c *gin.Context) {
	var categoryList []model.Category
	var list []gin.H
	database.DB.Order("ID asc").Find(&categoryList)

	for _, val := range categoryList {
		list = append(list, gin.H{
			"id":          val.ID,
			"name":        val.Name,
			"description": val.Description,
			"status":      val.Status,
		})
	}
	c.JSON(200, gin.H{
		"Status":  "success!",
		"Code":    200,
		"Message": "retrieved category details!",
		"Data": gin.H{
			"categories": list,
		},
	})
}
func AddCategory(c *gin.Context) {
	var categoryinfo model.Category
	err := c.ShouldBindJSON(&categoryinfo)
	if err != nil {
		c.JSON(200, gin.H{
			"Status":  "failed!",
			"Code":    400,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
		return
	}
	addCategory := database.DB.Create(&model.Category{
		Name:        categoryinfo.Name,
		Description: categoryinfo.Description,
		Status:      categoryinfo.Status,
	})
	if addCategory.Error != nil {
		c.JSON(200, gin.H{
			"Status":  "failed!",
			"Code":    400,
			"Message": "failed to add category",
			"Data":    gin.H{},
		})
	} else {
		c.JSON(200, gin.H{
			"Status":  "success!",
			"Code":    200,
			"Message": "category created successfully!",
			"Data":    gin.H{},
		})
	}
}
func EditCategory(c *gin.Context) {
	var categoryinfo model.Category
	err := c.ShouldBindJSON(&categoryinfo)
	if err != nil {
		c.JSON(200, gin.H{
			"Status":  "failed!",
			"Code":    400,
			"Message": "failed to bind json",
			"Data":    gin.H{},
		})
	}
	id := c.Param("ID")
	cerr := database.DB.Where("id=?", id).Updates(&categoryinfo)
	if cerr.Error != nil {
		c.JSON(200, gin.H{
			"Status":  "failed!",
			"Code":    404,
			"Message": "failed to edit category",
			"Data":    gin.H{},
		})
	}
	c.JSON(200, gin.H{
		"Status":  "success!",
		"Code":    200,
		"Message": "successfully editted",
		"Data":    gin.H{},
	})
}
func BlockCategory(c *gin.Context) {
	var category model.Category
	id := c.Param("ID")
	database.DB.First(&category, id)
	if category.Status == "blocked" {
		database.DB.Model(&category).Update("status", "active")
		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": "category active",
			"Data":    gin.H{},
		})
	} else {
		database.DB.Model(&category).Update("status", "blocked")
		c.JSON(200, gin.H{
			"Status":  "success!",
			"Code":    200,
			"Message": "category blocked!",
			"Data":    gin.H{},
		})
	}
}
func DeleteCategory(c *gin.Context) {
	var category model.Category
	id := c.Param("ID")
	err := database.DB.Where("id=?", id).Delete(&category)
	if err.Error != nil {
		c.JSON(200, gin.H{
			"Status":  "failed",
			"Code":    500,
			"Message": "failed to delete category",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "success!",
		"Code":    200,
		"Message": "category deleted successfully!",
		"Data":    gin.H{},
	})
}
