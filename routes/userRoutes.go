package routes

import (
	"net/http"
	controller "pkart/controller/users"
	"pkart/middleware"

	"github.com/gin-gonic/gin"
)
var roleuser = "User"

func UserRoutes(r *gin.RouterGroup) {

	/////////login
	r.POST("/signup", controller.UserSignUp)
	r.POST("/login", controller.UserLogin)
	r.GET("/otp", controller.OtpSignUp)
	r.GET("/resendotp", controller.ResendOtp)
 	r.DELETE("/logout", middleware.AuthMiddleware(roleuser), controller.Logout)

	///products
	r.GET("/products", controller.UserViewProducts)
	r.GET("/products/sort", middleware.AuthMiddleware(roleuser), controller.SortProduct)
	r.GET("/product/search",middleware.AuthMiddleware(roleuser), controller.SearchProduct)
	r.GET("/product/filter",middleware.AuthMiddleware(roleuser), controller.FilterProduct)

	//---------------Rating
	r.POST("/rating/:ID", middleware.AuthMiddleware(roleuser), controller.AddRating)
	r.PUT("/rating/:ID", middleware.AuthMiddleware(roleuser), controller.EditRating)///someissue

	///address
	r.POST("/address", middleware.AuthMiddleware(roleuser), controller.AddAddress)
	r.PUT("/address/:ID", middleware.AuthMiddleware(roleuser), controller.EditAddress)
	r.DELETE("/address/:ID", middleware.AuthMiddleware(roleuser), controller.DeleteAddress)
	r.GET("/listaddress", middleware.AuthMiddleware(roleuser), controller.ListAddress)

	///profile
	r.GET("/profile", middleware.AuthMiddleware(roleuser), controller.ShowProfile)
	r.PATCH("/profile", middleware.AuthMiddleware(roleuser), controller.EditProfile)

	///forget password
	r.POST("/forgetpassword", controller.ForgetPassword)
	r.GET("/checkotp", controller.CheckOtp)
	r.PATCH("/newpassword", controller.NewPassword)

	///cart
	r.GET("/cart", middleware.AuthMiddleware(roleuser), controller.ViewCart)
	r.POST("/cart/:ID", middleware.AuthMiddleware(roleuser), controller.AddToCart)
	r.PATCH("/cart/:ID", middleware.AuthMiddleware(roleuser), controller.RemoveCart)

	////// checkout
	r.POST("/checkout", middleware.AuthMiddleware(roleuser), controller.CartCheckOut)

	////// orders
	r.PATCH("/cancelorder/:ID", middleware.AuthMiddleware(roleuser), controller.CancelOrder)
	r.GET("/orders", middleware.AuthMiddleware(roleuser), controller.OrderView)
	r.GET("/orderdetails/:ID", middleware.AuthMiddleware(roleuser), controller.OrderDetails)

	////////payment
	r.GET("/payment", func(c *gin.Context) { c.HTML(http.StatusOK, "pay.html", nil) })
	r.POST("/payment/confirm", controller.PaymentConfirmation)

	/////// Wishlist
	r.GET("/wishlist", middleware.AuthMiddleware(roleuser), controller.ShowWishlist)
	r.POST("/wishlist/:ID", middleware.AuthMiddleware(roleuser), controller.AddWishlist)
	r.DELETE("/wishlist/:ID", middleware.AuthMiddleware(roleuser), controller.RemoveWishlist)

	///////Wallet
	r.GET("/wallet", middleware.AuthMiddleware(roleuser), controller.ShowWallet)

	////////Invoice
	r.GET("/order/invoice/:ID", middleware.AuthMiddleware(roleuser), controller.CreateInvoice)


}
