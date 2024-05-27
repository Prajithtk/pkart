package controller

import (
	"fmt"
	"os"
	"pkart/database"
	"pkart/helper"
	"pkart/model"

	"github.com/gin-gonic/gin"
)

type Razor struct {
	Order     string `json:"OrderID"`
	Payment   string `json:"PaymentID"`
	Signature string `json:"Signature"`
}

func RazorPay(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------RAZOR PAY------------------------")

	var payment model.Payment
	var orderitems []model.OrderItem
	var detail string

	orderId := c.Query("id")
	Logged := c.GetUint("userid")
	// Logged := c.MustGet("Id").(uint)

	if err := database.DB.Preload("User").First(&payment, "Pay_Id=?", orderId).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Order not found!",
			"Data":    gin.H{},
		})
		return
	}
	if payment.UserId != Logged {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Order not found!",
			"Data":    gin.H{},
		})
		return
	}
	if err := database.DB.Preload("Product").Find(&orderitems, "Order_Id=?", payment.OrderId).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Message": "Cannot fetch Order Items!",
			"Data":    gin.H{},
		})
		return
	}
	for i := 0; i < len(orderitems); i++ {
		if i == (len(orderitems) - 1) {
			detail += "and " + orderitems[i].Product.Name
		} else {
			detail += orderitems[i].Product.Name + ", "
		}
	}

	c.HTML(200, "razor.html", gin.H{
		"Order":   orderId,
		"Amounr":  payment.Amount,
		"Key":     os.Getenv("RAZOR_KEY"),
		"Name":    payment.User.Name,
		"Eamil":   payment.User.Email,
		"Phone":   payment.User.Phone,
		"Product": "Your products " + detail + ". Pay for them now!",
	})
}

func RazorPayVerify(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PAYMENT VERIFY------------------------")

	var verify Razor
	var order []model.OrderItem
	var payment model.Payment
	var ca []model.Cart

	Logged := c.MustGet("Id").(uint)

	err := c.ShouldBindJSON(&verify)
	if err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Couldn't find any data bind!",
			"Data":    gin.H{},
		})
		fmt.Println("Couldn't find any data bind!")
		return
	}

	er := database.DB.First(&payment, "Payment_Id=?", verify.Order).Error
	if er != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   er.Error(),
			"Message": "No such order found!",
			"Data":    gin.H{},
		})
		fmt.Println("No such order found!")
		return
	}

	if err := database.DB.Preload("Order").Preload("Product").Find(&order, "Order_Id=?", payment.OrderId).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "Couldn't find order items from databse!",
			"Data":    gin.H{},
		})
		fmt.Println("Couldn't find order items from databse!")
		return
	}
	eror := helper.RazorPaymentVerification(verify.Signature, verify.Order, verify.Payment)
	if eror != nil {
		c.JSON(402, gin.H{
			"Status":  "Error!",
			"Code":    402,
			"Error":   eror.Error(),
			"Message": "Payment failed!",
			"Data":    gin.H{},
		})
		fmt.Println("Payment failed!")
		return
	}
	payment.TransId = verify.Payment
	payment.Status = "recieved"
	erorr := database.DB.Save(&payment).Error

	errors := database.DB.Preload("Product").Find(&ca, "User_Id=?", Logged).Error
	if errors != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   errors.Error(),
			"Message": "Nothing in cart!",
			"Data":    gin.H{},
		})
		return
	}
	for _, v := range ca {
		v.Product.Quantity -= int(v.Quantity)
		erro := database.DB.Model(&v.Product).Updates(&v.Product).Error
		if erro != nil {
			c.JSON(402, gin.H{
				"Status":  "Error!",
				"Code":    402,
				"Error":   erro.Error(),
				"Message": "Error while updating product!",
				"Data":    gin.H{},
			})
			return
		}
	}

	for _, v := range ca {
		database.DB.Delete(&v)
	}

	if erorr != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   erorr.Error(),
			"Message": "Couldn't update payment success in databse!",
			"Data":    gin.H{},
		})
		fmt.Println("Couldn't update payment success in databse!")
		return
	}
	// if err := Invoice(c, payment.OrderId); err != nil {
	// 	c.JSON(400, gin.H{
	// 		"Status":  "Error!",
	// 		"Code":    400,
	// 		"Error":   err.Error(),
	// 		"Message": "Error on invoice create!",
	// 		"Data":    gin.H{},
	// 	})
	// 	fmt.Println("Error on invoice create!  ", err)
	// 	return
	// }
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    400,
		"Message": "Payment Succesfull!",
		"Data":    gin.H{},
	})
	fmt.Println("Payment Succesfull!")
}
