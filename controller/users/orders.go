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
		"Status": "success",
		"Code": 200,
		"Message": "order items are:",
		"Data": gin.H{
			"orders": orderShow,
		},
	})
}

func OrderDetails(c *gin.Context) {
	var orderitems []model.OrderItem
	var payment model.Payment
	var orderItemShow []gin.H
	orderId := c.Param("ID")
	if err := database.DB.Where("order_items.order_id=?", orderId).Preload("Order").Find(&orderitems).Error; err != nil {
		c.JSON(400, gin.H{
			"Status": "failed",
			"Code":   400,
			"Message":  "can't find order details",
			"Data": gin.H{},
		})
		return
	}
	if error := database.DB.Where("payments.order_id=?", orderId).First(&payment).Error; error != nil {
		c.JSON(400, gin.H{
			"Status": "failed",
			"Code":   400,
			"Message":  "can't find payment details",
			"Data": gin.H{},
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
		"code": 200,
		"Message": "order details are:",
		"Data" : gin.H{
		"orderDetails":    orderItemShow,
		"totalAmount":     totalAmount,
		"subTotal":        subTotal,
		"coupon discount": subTotal - totalAmount,
		"paymentStatus":   payment.Status,
		},
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
			"Status":  "failed",
			"Code":    402,
			"Message": "please provide a valid cancellation reason.",
			"Data": gin.H{},
		})
	} else {
		if err := tx.First(&orderItem, orderItemId).Error; err != nil {
			c.JSON(404, gin.H{
				"Status": "fail",
				"Code":   404,
				"Message":  "can't find order",
				"Data": gin.H{},
			})
			tx.Rollback()
			return
		}
		// fmt.Println("orderitems:", orderItem)
		if orderItem.Status == "cancelled" {
			c.JSON(202, gin.H{
				"Status":  "failed",
				"Code":    202,
				"Message": "product already cancelled",
				"Data": gin.H{},
			})
			return
		}
		// ======= update status as cancelled ======
		orderItem.Status = "cancelled"
		orderItem.CancelReason = reason
		if err := tx.Save(&orderItem).Error; err != nil {
			c.JSON(500, gin.H{
				"Status": "failed",
				"Code":   500,
				"Message":  "failed to  save changes to database.",
				"Data": gin.H{},
			})
			tx.Rollback()
			return
		}

		var orderAmount model.Orders
		if err := tx.First(&orderAmount, orderItem.OrderId).Error; err != nil {
			c.JSON(404, gin.H{
				"Status": "failed",
				"Code":   404,
				"Message":  "failed to find order details",
				"Data": gin.H{},
			})
			tx.Rollback()
			return
		}
		//========== check coupon condition ============
		var couponRemove model.Coupons
		if orderAmount.CouponCode != "" {
			if err := database.DB.First(&couponRemove, "code=?", orderAmount.CouponCode).Error; err != nil {
				c.JSON(404, gin.H{
					"Status": "failed",
					"Code":   404,
					"Message":  "can't find coupon code",
					"Data": gin.H{},
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
				"Status": "failed",
				"Code":   500,
				"Message":  "failed to update order details",
				"Data": gin.H{},
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
				"Data": gin.H{},
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
				"Data": gin.H{},
			})
			tx.Rollback()
		} else {
			c.JSON(201, gin.H{
				"Status":  "Success",
				"Code" : 201,
				"Message": "order cancelled",
				"Data":    orderItem.Status,
			})
		}
	}
}
