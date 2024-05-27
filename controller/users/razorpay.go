package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"pkart/database"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

func PaymentHandler(orderId int, amount int) (string, error) {

	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY"), os.Getenv("RAZORPAY_SECRET"))
	orderParams := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  strconv.Itoa(orderId),
	}
	order, err := client.Order.Create(orderParams, nil)
	if err != nil {
		return "", errors.New("PAYMENT NOT INITIATED")
	}
	razorId, _ := order["id"].(string)
	return razorId, nil
}

func PaymentConfirmation(c *gin.Context) {
	
	var paymentStore model.Payment
	var paymentDetails = make(map[string]string)
	// var cartItems model.Cart
	if err := c.BindJSON(&paymentDetails); err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
			"error":  "Invalid request body",
			"code":   400,
		})
		return
	}
	pd := paymentDetails
	fmt.Println(pd)
	//============== verify the signature ================
	err := RazorPaymentVerification(pd["signature"], pd["order_id"], pd["payment_id"])
	if err != nil {
		fmt.Println("-------", err)
		return
	}
	if err := database.DB.First(&paymentStore, "pay_id=?", pd["order_id"]).Error; err != nil {
		fmt.Println("can't find payment details")
		return
	}
	paymentStore.TransId = pd["payment_id"]
	paymentStore.Status = "success"
	database.DB.Save(&paymentStore)
	
	//============ quantity remove ================
	var productQuantity model.Products
	var productCheck []model.OrderItem
	if err := database.DB.Where("order_id=?", paymentStore.OrderId).Find(&productCheck).Error; err != nil {
		fmt.Println("cant find items")
	}
	fmt.Println(productCheck)
	for _, val := range productCheck {
		database.DB.First(&productQuantity, val.Id)
		productQuantity.Quantity -= val.Quantity
		if err := database.DB.Save(&productQuantity).Error; err != nil {
			fmt.Println("failed to save  updated quantity of products in db")
		}
	}
	fmt.Println("payment done , order placed successfully")
}

func RazorPaymentVerification(sign, orderId, paymentId string) error {
	signature := sign
	secret := os.Getenv("RAZORPAY_SECRET")
	data := orderId + "|" + paymentId
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
		return errors.New("PAYMENT FAILED")
	} else {
		return nil
	}
}
