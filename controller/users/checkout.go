package controller

import (
	"net/http"
	"pkart/database"
	"pkart/helper"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CartCheckOut(c *gin.Context) {

	var req struct {
		Coupon    string `json:"coupon"`
		Payment   string `json:"payment"`
		AddressId uint   `json:"addressid"`
	}
	var coupon model.Coupons
	var order model.Orders
	var subTotal int

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, gin.H{
			"Status":  "error",
			"Code":    404,
			"Message": "failed to bind JSON",
			"Data":    gin.H{},
		})
		return
	}
	userId := c.GetUint("userid")

	var address model.Address
	var cartItems []model.Cart
	database.DB.First(&address, "address_id=?", uint(req.AddressId)).Where("User_Id = ? ", userId)
	database.DB.Preload("Product").Find(&cartItems, "user_id=?", userId)

	if req.AddressId == 0 || address.UserId != userId {
		c.JSON(404, gin.H{
			"Status":  "error",
			"Code":    404,
			"Message": "no address found",
			"Data":    gin.H{},
		})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(404, gin.H{
			"Status":  "error",
			"Code":    404,
			"Message": "your cart is empty",
			"Data":    gin.H{},
		})
		return
	}
	for _, item := range cartItems {
		subTotal += int(item.Quantity) *(int(item.Product.Price)-item.Product.Offer)
	}
	// fmt.Println("subtotal:",subTotal)
	
	fetchCoupon := database.DB.First(&coupon, "code=?", req.Coupon)
	
	order.UserId = userId
	order.AddressId = req.AddressId
	order.Total = subTotal
	order.Amount = subTotal - (order.Total * coupon.Value / 100)
	if coupon.Min < subTotal && fetchCoupon.Error == nil {
		order.CouponId = coupon.Id
	} else if coupon.Min > subTotal {
		c.JSON(401, gin.H{
			"Status":  "Fail!",
			// "Code":    401,
			"Message": "this coupon is not for this amount",
			"Data":    gin.H{},
		})
		return

		} else if fetchCoupon.Error != nil {
			if req.Coupon == "" {
			order.CouponId = 1
			c.JSON(200, gin.H{
				"Status":  "error",
				// "Code":    404,
				"Message": "no coupon provided",
				"Data":    gin.H{},
			})

			} else {
				c.JSON(404, gin.H{
				"Status":  "error",
				"Code":    404,
				"Message": "not a valid coupon code!",
				"Data":    gin.H{},
			})
			return
		}
	}
		if req.Payment == "COD" {
			if order.Amount < 1000 {
				c.JSON(401, gin.H{
					"Status":  "Error!",
					"Code":    401,
					"Message": "Minimum order amount for Cash On Delivery is 1000",
					"Data":    gin.H{},
				})
				return
			}
		}
	num := helper.GenerateInt()
	numb, _ := strconv.Atoi(num)
	order.Id, _ = strconv.Atoi(num)

	database.DB.Create(&order)

	for _, list := range cartItems {
		orderitem := model.OrderItem{
			OrderId:   uint(order.Id),
			ProductId: list.ProductId,
			Quantity:  list.Quantity,
			SubTotal: float64(list.Product.Price-(list.Product.Offer))*float64(list.Quantity),
			Status:    "pending",
		}
		if err := database.DB.Create(&orderitem); err.Error != nil {
			c.JSON(403, gin.H{
				"Status":  "error",
				"Code":    403,
				"Message": "couldn't place the order. Please try again later.",
				"Error":   err.Error,
				"Data":    gin.H{},
			})
			return
		}
		list.Product.Quantity -= int(list.Quantity)
		database.DB.Model(list.Product).Update("quantity", list.Product.Quantity)
	}
	payment := model.Payment{
		OrderId: uint(order.Id),
		UserId:  userId,
		Amount:  order.Amount,
		Status:  "pending",
	}
	switch req.Payment {
	case "COD":
		// Handle Cash on Delivery
		payment.PayMeth = "Cash on Delivery"

		for _, val := range cartItems {
			database.DB.Delete(&val)
		}
		if err := database.DB.Create(&payment); err.Error != nil {
			c.JSON(403, gin.H{
				"Status":  "error",
				"Code":    403,
				"Message": "failed to process COD!! Try again later.",
				"Error":   err.Error,
				"Data":    gin.H{},
			})
			return
		}
		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": "order placed successfully.",
			"Data":    gin.H{},
		})

	case "PAY NOW":
		// Handle online payment
		payment.PayMeth = "Razor Pay"
		for _, val := range cartItems {
			database.DB.Delete(&val)
		}
		razorIde, err := PaymentHandler(numb, payment.Amount)

		if err != nil {
				c.JSON(406, gin.H{
				"Status":  "error",
				"Code":    406,
				"Message": "payment gateway not initiated.",
				"Data":    gin.H{},
			})
			return
		}
		payment.PayId = razorIde
		if error := database.DB.Create(&payment); error.Error != nil {
			c.JSON(403, gin.H{
				"Status":  "error",
				"Code":    403,
				"Message": "payment creation failed! Try again later.",
				"Data":    gin.H{},
			})
			return
		}

		c.JSON(200, gin.H{
			"Status":  "success",
			"Code":    200,
			"Message": "Complete the payment to place the order",
			"Data": gin.H{
				"Payment": razorIde,
				"Amount":	order.Amount,
			},
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}
}
