package controller

import (
	"pkart/database"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ShowOrders(c *gin.Context) {

	var order []model.OrderItem
	var show []gin.H

	err := database.DB.Preload("Order").Preload("Product").Preload("Order.User").Find(&order).Error

	if err != nil {
		c.JSON(404, gin.H{
			"Status":  "failed",
			"Code":    404,
			"Message": "no orders found",
			"Data":    gin.H{},
		})
		return
	}
	for _, val := range order {

		show = append(show, gin.H{
			"Id":           val.Id,
			"OrderId":      val.OrderId,
			"Username":     val.Order.User.Name,
			"User_Email":   val.Order.User.Email,
			"Product_Name": val.Product.Name,
			"Image1":       val.Product.Image1,
			"Quantity":     val.Quantity,
			"Status":       val.Status,
		})
	}
	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "order items are:",
		"Data":    show})
}

func OrdersStatusChange(c *gin.Context) {

	status := c.Request.FormValue("status")
	ord, _ := strconv.Atoi(c.Query("order"))
	var order model.OrderItem

	if err := database.DB.Preload("Order").Preload("Order.Coupon").Preload("Product").First(&order, uint(ord)).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "error",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "No such Order found",
			"Data":    gin.H{},
		})
		return
	}
	if status == "shipped" {
		order.Status = status
	} else if status == "delivered" {
		order.Status = status
	} else {
		c.JSON(400, gin.H{
			"Status":  "error",
			"Code":    400,
			"Message": "this status can't be assigned",
			"Data":    gin.H{},
		})
		return
	}
	er := database.DB.Save(&order).Error
	if er != nil {
		c.JSON(401, gin.H{
			"Status":  "error",
			"Code":    400,
			"Error":   er.Error(),
			"Message": "couldn't change the order status",
			"Data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "Success",
		"Code":    200,
		"Message": "Order status updated successfully!",
		"Data":    gin.H{},
	})
}
