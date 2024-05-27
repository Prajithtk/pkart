package controller

import (
	"net/http"
	"pkart/database"
	"strings"

	"github.com/gin-gonic/gin"
)

func SortProduct(c *gin.Context) {

	type products struct {
		Name        string `json:"name"`
		Price       uint   `json:"price"`
		Color       string `json:"color"`
		Quantity    uint   `json:"quantity"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Image1      string
	}

	var req struct {
		Sort string `json:"sort"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse JSON"})
		return
	}

	sort := strings.ToLower(req.Sort)
	var product []products

	switch sort {
	case "asc":
		database.DB.Order("name asc").Find(&product)
	case "desc":
		database.DB.Order("name desc").Find(&product)
	case "highlow":
		database.DB.Order("price desc").Find(&product)
	case "lowhigh":
		database.DB.Order("price asc").Find(&product)
	case "latest":
		database.DB.Order("created_at desc").Find(&product)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please give correct options"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Products": product})
}
