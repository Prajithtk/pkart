package controller

import (
	"pkart/database"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("ID"))
	userId := c.GetUint("userid")
	var product model.Products
	var cart model.Cart
	if err := database.DB.Where("id=?", id).First(&product).Error; err != nil {
		c.JSON(404, gin.H{"error": "can't find the product"})
	} else {
		err := database.DB.Where("product_id=?", id).First(&cart)
		if err.Error == nil {
			if cart.Quantity < 10 && cart.Quantity < product.Quantity {
				cart.Quantity++
				database.DB.Save(&cart)
				c.JSON(200, gin.H{"message": "quantity added to the cart"})
			} else {
				c.JSON(404, gin.H{"error": "can't add more of this product"})
				return
			}
		} else {
			cart = model.Cart{
				UserId:    userId,
				ProductId: uint(id),
				Quantity:  1,
			}
			database.DB.Create(&cart)
			c.JSON(200, gin.H{"message": "product added to cart successfully"})
		}
	}
}

func RemoveCart(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("ID"))
	userId := c.GetUint("userid")
	var cart model.Cart
	err := database.DB.First(&cart, "user_id=? AND product_id=?", userId, id)
	if err != nil {
		if cart.Quantity <= 1 {
			database.DB.Delete(&cart)
			c.JSON(200, gin.H{"message": "product is removed from cart"})
		} else {
			cart.Quantity--
			database.DB.Save(&cart)
			c.JSON(200, gin.H{"message": "quantity is reduced by 1"})
		}
	} else {
		c.JSON(404, gin.H{"error": "product not found in cart"})
	}
}

func ViewCart(c *gin.Context) {
	type showcart struct {
		Id          uint
		Product     string
		Quantity    uint
		Description string
		Price       int
		OfferPrice	int
	}
	var cart []model.Cart
	var products []model.Products
	var show []showcart
	var total int
	userId := c.GetUint("userid")
	database.DB.Find(&cart, "user_id=?", userId)
	for i := 0; i < len(cart); i++ {
		var product model.Products
		database.DB.First(&product, cart[i].ProductId)
		products = append(products, product)
	}
	for i := 0; i < len(cart); i++ {
		l := showcart{
			Id:          products[i].ID,
			Product:     products[i].Name,
			Quantity:    uint(cart[i].Quantity),
			Description: products[i].Description,
			Price:       int(products[i].Price),
			OfferPrice: int(products[i].Price) - products[i].Offer,
		}
		total += int(l.Quantity) * l.OfferPrice
		show = append(show, l)
	}
	c.JSON(200, gin.H{
		"Products":     show,
		"Total Amount": total,
	})
}
