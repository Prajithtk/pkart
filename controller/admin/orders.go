package controller

import (
	"fmt"
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
		c.JSON(404, gin.H{"Message": "No orders found!"})
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
		"success": "true",
		"message": "Order items are:",
		"values":  show})
}

func OrdersStatusChange(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ORDER STATUS CHANGING------------------------")

	status := c.Request.FormValue("status")

	ord, _ := strconv.Atoi(c.Query("order"))

	fmt.Println(status)

	var order model.OrderItem

	if err := database.DB.Preload("Order").Preload("Order.Coupon").Preload("Product").First(&order, uint(ord)).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "No such Order found!",
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
			"Status":  "Error!",
			"Code":    400,
			"Message": "This status can't be assigned!",
			"Data":    gin.H{},
		})
		return
	}
	er := database.DB.Save(&order).Error
	if er != nil {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   er.Error(),
			"Message": "Couldn't change the order status!",
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

// func EditOrder(c *gin.Context) {

// 	var req struct {
// 		OdrId  uint   `json:"id"`
// 		Status string `json:"status"`
// 	}
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(404, gin.H{
// 			"success": "false",
// 			"message": "failed to bind json"})
// 		return
// 	}

// 	var orderitem model.OrderItem
// 	var payment model.Payment
// 	var wallet model.Wallet

// 	if err := database.DB.Preload("Order").Preload("Order.Coupon").Preload("Product").First(&orderitem, req.OdrId).Error; err != nil {
// 		c.JSON(404, gin.H{
// 			"success": "false",
// 			"message": "Order not found!"})
// 		return
// 	}
// 	if err := database.DB.First(&payment, "Order_Id=?", orderitem.OrderId).Error; err != nil {
// 		c.JSON(500, gin.H{
// 			"success": "false",
// 			"message": "No such payment!"})
// 		return
// 	}
// 	if err := database.DB.First(&wallet, "User_Id=?", orderitem.Order.UserId).Error; err != nil {
// 		c.JSON(501, gin.H{
// 			"success": "false",
// 			"message": "Failed to find the user wallet!"})
// 		return
// 	}
// 	if req.Status == "cancelled" {
// 		if orderitem.Status == "cancelled" {
// 			c.JSON(409, gin.H{
// 				"success": "true",
// 				"message": "This order is already cancelled"})
// 			return
// 		}

// 		orderitem.Status = req.Status

// 		orderitem.Order.Total = orderitem.Order.Total - orderitem.Product.Price*orderitem.Quantity
// 		if orderitem.Order.Total < (orderitem.Order.Coupon.Min) {
// 			orderitem.Order.Amount = orderitem.Order.Total
// 			orderitem.Order.CouponId = 1
// 		} else {
// 			orderitem.Order.Amount = orderitem.Order.Total - (orderitem.Order.Total * (orderitem.Order.Coupon.Value) / 100)
// 		}
// 		if er := database.DB.Save(&orderitem.Order).Error; er != nil {
// 			c.JSON(500, gin.H{
// 				"success": "false",
// 				"message": "Can't decrease the order amount!"})
// 			return
// 		}
// 		orderitem.Product.Quantity += orderitem.Quantity
// 		if er := database.DB.Save(&orderitem.Product).Error; er != nil {
// 			c.JSON(500, gin.H{
// 				"success": "false",
// 				"message": "Can't increase product quantity!"})
// 			return
// 		}
// 		if payment.Status == "recieved" {
// 			wallet.Amount += ((orderitem.Product.Price) * (orderitem.Quantity)) - ((orderitem.Product.Price) * (orderitem.Quantity) * (orderitem.Order.Coupon.Value) / 100)
// 			if err := database.DB.Save(&wallet).Error; err != nil {
// 				c.JSON(500, gin.H{
// 					"success": "false",
// 					"message": "Couldn't update wallet!"})
// 				return
// 			}
// 			payment.Status = "partially refunded"
// 			if err := database.DB.Model(&payment).Update("Status", payment.Status).Error; err != nil {
// 				c.JSON(500, gin.H{
// 					"success": "false",
// 					"message": "Failed to set payment as refunded"})
// 				return
// 			}
// 		}
// 	} else if req.Status == "shipped" {
// 		orderitem.Status = req.Status
// 	} else if req.Status == "delivered" {
// 		orderitem.Status = req.Status
// 	} else {
// 		c.JSON(400, gin.H{
// 			"success": "false",
// 			"message": "This status can't be assigned!"})
// 		return
// 	}
// 	er := database.DB.Save(&orderitem).Error
// 	if er != nil {
// 		c.JSON(401, gin.H{
// 			"success": "false",
// 			"message": "Couldn't change the order status!"})
// 		return
// 	}
// 	if req.Status == "cancelled" {
// 		c.JSON(200, gin.H{
// 			"success": "false",
// 			"Message": "Order cancelled succesfully!"})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"success": "false",
// 		"Message": "Order status updated successfully!"})
// }
