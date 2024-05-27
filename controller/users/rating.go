package controller

import (
	"fmt"
	"pkart/database"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"

)

func AddRating(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDING RATING------------------------")

	var rating, alrdyRate model.Rating
	var rates []model.Rating
	var product model.Products
	var sum float32

	userId := c.GetUint("userid")
	ProductId, _ := strconv.Atoi(c.Param("ID"))
	// id, _ := strconv.Atoi(c.Query("id"))
	er := database.DB.Where("Product_Id=? AND User_Id=?", ProductId, userId).First(&alrdyRate).Error

	rate, _ := strconv.Atoi(c.Request.FormValue("rating"))
	rating.Review = c.Request.FormValue("review")
	fmt.Println(rate)
	fmt.Println(rating.Review)
	rating.Rating = float32(rate)

	if rating.Rating > 5 {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "Rating should be in between 1 and 5",
			"Data":    gin.H{},
		})
		return
	}

	err := database.DB.First(&product, ProductId).Error
	database.DB.Find(&rates, "Product_Id=?", uint(ProductId))

	if er != nil {
		if err != nil {
			c.JSON(404, gin.H{
				"Status":  "Error!",
				"Code":    404,
				"Message": "Product not found!",
				"Data":    gin.H{},
			})
		} else {
			rating.ProductId = uint(ProductId)
			rating.UserId = userId

			database.DB.Create(&rating)
			for _, v := range rates {
				sum += v.Rating
			}
			product.AvrgRating = sum / float32(len(rates))
			database.DB.Save(&product)
			c.JSON(201, gin.H{
				"Status":  "Success!",
				"Code":    201,
				"Message": "Rating and review added successfully!",
				"Data":    gin.H{},
			})
		}
	} else {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "Rating or review  already exists, Try to update it instead of adding again!",
			"Data":    gin.H{},
		})
	}

}

func EditRating(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------EDIT RATING------------------------")

	// id, _ := strconv.Atoi(c.Query("Id"))
	userId := c.GetUint("userid")
	ProductId, _ := strconv.Atoi(c.Param("ID"))

	fmt.Println("productid:",ProductId)
	fmt.Println("userid:",userId)
	var rate model.Rating
	var rates []model.Rating
	var sum float32

	fmt.Println("rates:",rates)
	if err := database.DB.Where("Product_Id=?", uint(ProductId)).Find(&rates).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Couldn't find any rating for review for this product!",
			"Data":    gin.H{},
		})
		return
	}
	fmt.Println("rates:",rates)

	if err := database.DB.Preload("Product").Where("User_Id=? AND Product_Id=?", userId, uint(ProductId)).First(&rate).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Rating and review not found! Please add first!",
			"Data":    gin.H{},
		})
		return
	}
	fmt.Println("rate:",rate )

	rating, _ := strconv.Atoi(c.Request.FormValue("rating"))
	rate.Review = c.Request.FormValue("review")
	rate.Rating = float32(rating)
	if rate.Rating > 5 {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "Rating should be less than or equal to 5!",
			"Data":    gin.H{},
		})
		return
	}
	fmt.Println("rating:",rating)
	fmt.Println("newrate:",rate)
	// if aerr := database.DB.Where("product_id=? AND user_id=?", ProductId, userId).Updates(&rating).Error; aerr!= nil{
	// 	c.JSON(500,gin.H{
	// 		"message": "failed to update review",
	// 	})
	// }

	// if err := database.DB.Model(&rate).Updates(&rating).Error; err != nil {
	if err := database.DB.Model(&rate).Updates(&rating).Error; err != nil {
		c.JSON(500, gin.H{
			"Status":  "Error!",
			"Code":    500,
			"Message": "Couldn't update the rating or review! try again",
			"Data":    gin.H{},
		})
		return
	}
	for _, v := range rates {
		sum += v.Rating
	}
	rate.Product.AvrgRating = sum / float32(len(rates))
	if err := database.DB.Save(&rate.Product).Error; err != nil {
		c.JSON(500, gin.H{
			"Status":  "Error!",
			"Code":    500,
			"Error":   err.Error(),
			"Message": "Couldn't update the rating or review!",
			"Data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "Error!",
		"Code":    404,
		"Message": "The rating has been updated!",
		"Data":    gin.H{},
	})
}
