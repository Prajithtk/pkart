package controller

import (
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

func OrderView(c *gin.Context) {
	var orders []model.Orders
	var orderShow []gin.H
	userId := c.GetUint("userid")
	database.DB.Where("user_id=?", userId).Find(&orders)
	for _, val := range orders {
		orderShow = append(orderShow, gin.H{
			"orderId":     val.Id,
			"userName":    val.UserId,
			"addressId":   val.AddressId,
			"orderAmount": val.Amount,
			"orderDate":   val.CreatedAt,
		})
	}
	c.JSON(200, gin.H{
		"status": "success",
		"orders": orderShow,
	})
}

func OrderDetails(c *gin.Context) {
	var orderitems []model.OrderItem
	var payment model.Payment
	var orderItemShow []gin.H
	orderId := c.Param("ID")
	if err := database.DB.Where("order_items.order_id=?", orderId).Preload("Order").Find(&orderitems).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Can't find order details",
			"code":   400,
		})
		return
	}
	if error := database.DB.Where("payments.order_id=?", orderId).First(&payment).Error; error != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Can't find payment details",
			"code":   400,
		})
		return
	}
	// fmt.Println("orderitems:", orderitems)
	var subTotal float64
	var totalAmount float64
	// var orderPayment string
	for _, val := range orderitems {
		subTotal = float64(val.Order.Total)
		totalAmount = float64(val.Order.Amount)
		// orderPayment = val.Order.OrderPaymentMethod
		orderItemShow = append(orderItemShow, gin.H{
			"OrderItemId":     val.Id,
			"product":         val.Product.Name,
			"quantity":        val.Quantity,
			"Amount":          val.SubTotal,
			"Delivery status": val.Status,
			"orderDate":       val.Order.CreatedAt,
			"couponcode":      val.Order.CouponId,
			"addressId":       val.Order.AddressId,
		})
	}

	c.JSON(200, gin.H{
		"status":          "Success",
		"orderDetails":    orderItemShow,
		"totalAmount":     totalAmount,
		"subTotal":        subTotal,
		"coupon discount": subTotal - totalAmount,
		"paymentStatus":   payment.Status,
	})
}

func CancelOrder(c *gin.Context) {
	var orderItem model.OrderItem
	orderItemId := c.Param("ID")
	reason := c.Request.FormValue("reason")
	// fmt.Println("orderitemid:", orderItemId)
	tx := database.DB.Begin()
	if reason == "" {
		c.JSON(402, gin.H{
			"status":  "Fail",
			"message": "Please provide a valid cancellation reason.",
			"code":    402,
		})
	} else {
		if err := tx.First(&orderItem, orderItemId).Error; err != nil {
			c.JSON(404, gin.H{
				"status": "Fail",
				"error":  "can't find order",
				"code":   404,
			})
			tx.Rollback()
			return
		}
		// fmt.Println("orderitems:", orderItem)
		if orderItem.Status == "cancelled" {
			c.JSON(202, gin.H{
				"status":  "Fail",
				"message": "product already cancelled",
				"code":    202,
			})
			return
		}
		// ======= update status as cancelled ======
		orderItem.Status = "cancelled"
		orderItem.CancelReason = reason
		if err := tx.Save(&orderItem).Error; err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "Failed to  save changes to database.",
				"code":   500,
			})
			tx.Rollback()
			return
		}

		var orderAmount model.Orders
		if err := tx.First(&orderAmount, orderItem.OrderId).Error; err != nil {
			c.JSON(404, gin.H{
				"status": "Fail",
				"error":  "failed to find order details",
				"code":   404,
			})
			tx.Rollback()
			return
		}
		//========== check coupon condition ============
		var couponRemove model.Coupons
		if orderAmount.CouponCode != "" {
			if err := database.DB.First(&couponRemove, "code=?", orderAmount.CouponCode).Error; err != nil {
				c.JSON(404, gin.H{
					"status": "Fail",
					"error":  "can't find coupon code",
					"code":   404,
				})
				tx.Rollback()
			}
		}
		if couponRemove.Min > int(orderAmount.Total) {
			orderAmount.Amount += couponRemove.Value
			orderAmount.Amount -= int(orderItem.SubTotal)
			orderAmount.CouponCode = ""
		}
		if err := tx.Save(&orderAmount).Error; err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to update order details",
				"code":   500,
			})
			tx.Rollback()
			return
		}
		var walletUpdate model.Wallet
		if err := tx.First(&walletUpdate, "user_id=?", orderAmount.UserId).Error; err != nil {
			c.JSON(501, gin.H{
				"status": "Fail",
				"error":  "failed to fetch wallet details",
				"code":   501,
			})
			tx.Rollback()
			return
		} else {
			walletUpdate.Amount += int(orderItem.Amount)
			tx.Save(&walletUpdate)
		}
		if err := tx.Commit().Error; err != nil {
			c.JSON(201, gin.H{
				"status":  "Fail",
				"message": "failed to commit transaction",
				"code":    201,
			})
			tx.Rollback()
		} else {
			c.JSON(201, gin.H{
				"status":  "Success",
				"message": "Order Cancelled",
				"data":    orderItem.Status,
			})
		}
	}
}
