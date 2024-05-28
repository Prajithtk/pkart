package controller

import (
	"fmt"
	"net/http"
	"pkart/database"
	"pkart/helper"
	"pkart/middleware"
	"pkart/model"
	onetp "pkart/onetp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var UserInfo model.Users

var RoleUser = "User"

func UserSignUp(c *gin.Context) {
	// var userInfo model.Users
	err := c.ShouldBindJSON(&UserInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "failed to bind json"})
		return
	}
	var ifUser model.Users
	uerr := database.DB.Where("email=?", UserInfo.Email).First(&ifUser)
	if uerr.Error == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "User Already exist"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(UserInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(401, gin.H{"success": "false", "message": "Failed to hash password"})
		return
	}
	UserInfo.Password = string(hashedPassword)

	UserInfo.Status = "Active"
	Code, _ := helper.GenerateRandomAlphanumericCode(6)
	UserInfo.ReferalCode = Code

	otp := onetp.GenerateOTP(6)
	newOtp := model.Otp{
		Otp:     otp,
		Email:   UserInfo.Email,
		Expires: time.Now().Add(1 * time.Minute),
	}

	database.DB.Create(&newOtp)

	onetp.SendOtp(newOtp.Email, newOtp.Otp)
	c.JSON(200, gin.H{"success": "true", "message": "OTP sent succesfully",
		"otp": otp})
}

func OtpSignUp(c *gin.Context) {
	var otp model.Otp
	err := c.BindJSON(&otp)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "failed to bind json"})
		return
	}
	var notp model.Otp
	oerr := database.DB.Where("Email=?", UserInfo.Email).Find(&notp)
	if oerr.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "failed to fetch OTP"})
		return
	}

	// currentTime := time.Now()
	// if currentTime.After(notp.Expires) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "OTP expired"})
	// 	return
	// }

	if otp.Otp != notp.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "invalid otp"})
		return
	}

	create := database.DB.Create(&UserInfo)
	fmt.Println(UserInfo)
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "message": "failed to create user"})
		return
	} else {
		c.JSON(200, gin.H{"success": "true", "message": "user created successfully"})
	}
	var userFetchData model.Users
	if err := database.DB.First(&userFetchData, "email=?", UserInfo.Email).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to fetch user details for wallet",
			"code":   400,
		})
		return
	}
	var userFetch model.Users
	if UserInfo.ReferalCode != ""{

		if err := database.DB.First(&userFetch, "referal_code=?", UserInfo.ReferalCode).Error; err != nil {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "failed to fetch user details for wallet",
				"code":   400,
			})
			return
		}
	}
	var wallet model.Wallet
	wallet.UserId =  userFetchData.ID
	wallet.Amount = 50
	database.DB.Create(&wallet)

	
	c.SetCookie("sessionId", "", -1, "/", "", false, false)
	c.JSON(201, gin.H{
		"status":  "Success",
		"message": "user created successfully",
	})
}

func ResendOtp(c *gin.Context) {
	var resend model.Otp
	err := c.ShouldBindJSON(&resend)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "failed to fetch otp"})
		return
	}
	var exotp model.Otp
	perr := database.DB.Where("email=?", resend.Email).First(&exotp)
	if perr.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "email not found"})
		return
	}

	newOtp := onetp.GenerateOTP(6)

	res := database.DB.Model(&model.Otp{}).Where("email=?", resend.Email).Updates(model.Otp{Otp: newOtp, Expires: time.Now().Add(1 * time.Minute)})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "message": "failed to resend otp"})
		return
	}
	onetp.SendOtp(resend.Email, newOtp)
	c.JSON(200, gin.H{"success": "true", "message": "otp resent successfully"})
}

func UserLogin(c *gin.Context) {
	var userDetails model.Users
	err := c.ShouldBindJSON(&userDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "failed to bind json"})
		return
	}
	var checkUser model.Users
	perror := database.DB.Where("email=?", userDetails.Email).First(&checkUser)
	if checkUser.Status == "blocked" {
		c.JSON(404, gin.H{"success": "false", "message": "Your account is blocked"})
	}
	if perror.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "User not found"})
		return
	}
	hashpassword := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(userDetails.Password))
	if hashpassword != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": "false", "message": "Incorrect Password"})
		return
	}
	token, err := middleware.JwtToken(c, checkUser.ID, checkUser.Email, RoleUser)

	if err != nil {
		c.JSON(403, gin.H{"success": "false", "message": "failed to create token"})
		return
	}

	c.SetCookie("JwtTokenUser", token, int((time.Hour * 100).Seconds()), "/", "localhost", false, false)
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Successfully Logged in!",
		"Data": gin.H{
			"Token": token,
			"Id":    checkUser.ID,
		},
	})
}
func Logout(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER LOGGING OUT----------------------")

	c.SetCookie("JwtUser", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Logged out successfully!",
		"Data":    gin.H{},
	})
}

// func UserViewProducts(c *gin.Context) {
// 	var productList []model.Products
// 	database.DB.Preload("Category").Order("ID asc").Find(&productList)
// 	for _, val := range productList {
// 		c.JSON(200, gin.H{
// 			"id":          val.ID,
// 			"name":        val.Name,
// 			"color":       val.Color,
// 			"quantity":    val.Quantity,
// 			"description": val.Description,
// 			// "categiryid":  val.CategoryId,
// 			"category": val.Category.Name,
// 			"status":   val.Status,
// 			"image1":   val.Image1,
// 			"image2":   val.Image2,
// 			"image3":   val.Image3,
// 		})
// 	}
// }

// func AddAddress(c *gin.Context) {
// 	userId := c.GetUint("userid")
// 	var address model.Address
// 	if err := c.BindJSON(&address); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if len(address.PinCode) != 6 {
// 		c.JSON(400, gin.H{"error": "pincode should be 6 digits"})
// 		return
// 	}
// 	err := database.DB.Create(&model.Address{
// 		BuildingName: address.BuildingName,
// 		Street:       address.Street,
// 		City:         address.City,
// 		State:        address.State,
// 		Landmark:     address.Landmark,
// 		PinCode:      address.PinCode,
// 		UserId:       userId,
// 	}).Error
// 	if err != nil {
// 		c.JSON(400, map[string]interface{}{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(200, "message: successfully added an address")
// }

// func EditAddress(c *gin.Context) {
// 	var addressDetails model.Address
// 	err := c.ShouldBindJSON(&addressDetails)
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "failed to bind json"})
// 		return
// 	}
// 	userId := c.GetUint("userid")
// 	addId := c.Param("ID")
// 	aerr := database.DB.Where("address_id=? AND user_id=?", addId, userId).Updates(&addressDetails)
// 	if aerr.Error != nil {
// 		c.JSON(400, gin.H{"error": "failed to edit address"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"message": "successfully edited the address"})
// }

// func DeleteAddress(c *gin.Context) {
// 	var address model.Address
// 	adrId := c.Param("ID")
// 	userId := c.GetUint("userid")
// 	if err := database.DB.Where("address_id=? AND user_id=?", adrId, userId).First(&address).Error; err != nil {
// 		c.JSON(400, gin.H{"error": "failed to fetch address"})
// 		fmt.Println(err)
// 		return
// 	}
// 	res := database.DB.Delete(&address)
// 	if res.Error != nil {
// 		c.JSON(400, gin.H{"error": "failed to delete address"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"message": "successfully deleted address"})
// }

// func ListAddress(c *gin.Context) {
// 	userId := c.GetUint("userid")
// 	var addressList []model.Address
// 	database.DB.Order("address_id asc").Find(&addressList).Where("userid=?", userId)
// 	for _, val := range addressList {
// 		c.JSON(200, gin.H{
// 			"id":            val.AddressId,
// 			"building name": val.BuildingName,
// 			"street":        val.Street,
// 			"city":          val.City,
// 			"state":         val.State,
// 			"landmark":      val.Landmark,
// 			"pincode":       val.PinCode,
// 		})
// 	}
// }

// -----------------------------profile management--------------------------------//
//--------------------------------------------------------------------------------//

// func ShowProfile(c *gin.Context) {
// 	userId := c.GetUint("userid")
// 	var userProfile model.Users
// 	// var userInfo []gin.H
// 	if err := database.DB.Preload("Address").Where("id=?", userId).First(&userProfile).Error; err != nil {
// 		c.JSON(500, gin.H{"error": "failed to fetch user profile"})
// 		return
// 	}
// 	response := gin.H{
// 		"Name":   userProfile.Name,
// 		"Email":  userProfile.Email,
// 		"Phone":  userProfile.Phone,
// 		"Gender": userProfile.Gender,
// 	}
// 	if len(userProfile.Address) > 0 {
// 		response["Address"] = userProfile.Address
// 	}
// 	c.JSON(200, gin.H{
// 		"success": "true",
// 		"message": "profile details are:",
// 		"values":  response})
// }

// func EditProfile(c *gin.Context) {
// 	var editdata model.Users
// 	fmt.Println("hello world")
// 	if err := c.ShouldBindJSON(&editdata); err != nil {
// 		c.JSON(400, gin.H{"error": "failed to bind json"})
// 		return
// 	}
// 	fmt.Println(editdata)
// 	fmt.Println("editdata")
// 	userId := c.GetUint("userid")
// 	fmt.Println(userId)
// 	if err := database.DB.Where("id=?", userId).Updates(&editdata); err.Error != nil {
// 		c.JSON(404, gin.H{"error": "failed to update address"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"message": "profile editted successfully"})
// }

// var Useremail model.Users

// func ForgetPassword(c *gin.Context) {
// 	var newOtp model.Otp
// 	var CheckOtp model.Otp
// 	var checkMail model.Users
// 	if err := c.ShouldBindJSON(&Useremail); err != nil {
// 		c.JSON(400, gin.H{"error": "failed to bind json"})
// 		return
// 	}
// 	if Useremail.Email == "" {
// 		c.JSON(303, gin.H{"error": "enter the email"})
// 		return
// 	}
// 	if err := database.DB.Where("email=?", Useremail.Email).First(&checkMail).Error; err != nil {
// 		c.JSON(400, gin.H{"error": "no account found,SignUp"})
// 		return
// 	}
// 	otp := onetp.GenerateOTP(6)
// 	newOtp = model.Otp{
// 		Otp:     otp,
// 		Email:   Useremail.Email,
// 		Expires: time.Now().Add(3 * time.Minute),
// 	}
// 	fmt.Println(newOtp)
// 	if err := onetp.SendOtp(Useremail.Email, otp); err != nil {
// 		c.JSON(500, gin.H{"error": "failed to send otp"})
// 		return
// 	}
// 	if res := database.DB.Where("email=?", Useremail.Email).First(&CheckOtp).Error; res != nil {
// 		if err := database.DB.Create(&newOtp).Error; err != nil {
// 			c.JSON(500, gin.H{"error": "failed to save otp"})
// 			return
// 		}
// 	} else {
// 		if err := database.DB.Model(&CheckOtp).Where("email=?", Useremail.Email).Updates(&newOtp).Error; err != nil {
// 			c.JSON(500, gin.H{"error": "failed to update otp table"})
// 			return
// 		}
// 	}
// 	c.JSON(303, gin.H{"message": "OTP sent to successfully", "otp": otp})
// }

// func CheckOtp(c *gin.Context) {
// 	var otp model.Otp
// 	if err := c.ShouldBindJSON(&otp); err != nil {
// 		c.JSON(404, gin.H{"error": "failed to bind json"})
// 		return
// 	}
// 	var notp model.Otp
// 	oerr := database.DB.Where("email=?", Useremail.Email).First(&notp)
// 	fmt.Println(Useremail.Email)
// 	if oerr.Error != nil {
// 		c.JSON(401, gin.H{"error": "failed to fetch OTP"})
// 		return
// 	}
// 	currentTime := time.Now()
// 	if currentTime.After(notp.Expires) {
// 		c.JSON(401, gin.H{"error": "OTP expired"})
// 		return
// 	}
// 	if otp.Otp != notp.Otp {
// 		c.JSON(400, gin.H{"error": "invalid otp"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"message": "Enter new password"})
// }

// func NewPassword(c *gin.Context) {
// 	var newPass model.Users
// 	if err := c.ShouldBindJSON(&newPass); err != nil {
// 		c.JSON(404, gin.H{"error": "failed to bind json"})
// 	}
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.JSON(401, gin.H{"error": "Failed to hash password"})
// 		return
// 	}
// 	newPass.Password = string(hashedPassword)
// 	if err := database.DB.Model(&newPass).Where("email=?", Useremail.Email).Updates(&newPass).Error; err != nil {
// 		c.JSON(404, gin.H{"error": "failed to update password"})
// 	}
// 	c.JSON(200, gin.H{"message": "successfully updated password"})
// }

// -------------------------------cart management---------------------------------//
//--------------------------------------------------------------------------------//

// func AddToCart(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("ID"))
// 	userId := c.GetUint("userid")
// 	var product model.Products
// 	var cart model.Cart
// 	if err := database.DB.Where("id=?", id).First(&product).Error; err != nil {
// 		c.JSON(404, gin.H{"error": "can't find the product"})
// 	} else {
// 		err := database.DB.Where("product_id=?", id).First(&cart)
// 		if err.Error == nil {
// 			if cart.Quantity < 10 && cart.Quantity < product.Quantity {
// 				cart.Quantity++
// 				database.DB.Save(&cart)
// 				c.JSON(200, gin.H{"message": "quantity added to the cart"})
// 			} else {
// 				c.JSON(404, gin.H{"error": "can't add more of this product"})
// 				return
// 			}
// 		} else {
// 			cart = model.Cart{
// 				UserId:    userId,
// 				ProductId: uint(id),
// 				Quantity:  1,
// 			}
// 			database.DB.Create(&cart)
// 			c.JSON(200, gin.H{"message": "product added to cart successfully"})
// 		}
// 	}
// }

// func RemoveCart(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("ID"))
// 	userId := c.GetUint("userid")
// 	var cart model.Cart
// 	err := database.DB.First(&cart, "user_id=? AND product_id=?", userId, id)
// 	if err != nil {
// 		if cart.Quantity <= 1 {
// 			database.DB.Delete(&cart)
// 			c.JSON(200, gin.H{"message": "product is removed from cart"})
// 		} else {
// 			cart.Quantity--
// 			database.DB.Save(&cart)
// 			c.JSON(200, gin.H{"message": "quantity is reduced by 1"})
// 		}
// 	} else {
// 		c.JSON(404, gin.H{"error": "product not found in cart"})
// 	}
// }

// func ViewCart(c *gin.Context) {
// 	type showcart struct {
// 		Product     string
// 		Quantity    uint
// 		Description string
// 		Price       int
// 	}
// 	var cart []model.Cart
// 	var products []model.Products
// 	var show []showcart
// 	var total int
// 	userId := c.GetUint("userid")
// 	database.DB.Find(&cart, "user_id=?", userId)
// 	for i := 0; i < len(cart); i++ {
// 		var product model.Products
// 		database.DB.First(&product, cart[i].ProductId)
// 		products = append(products, product)
// 	}
// 	for i := 0; i < len(cart); i++ {
// 		l := showcart{
// 			Product:     products[i].Name,
// 			Quantity:    uint(cart[i].Quantity),
// 			Description: products[i].Description,
// 			Price:       int(products[i].Price),
// 		}
// 		total += int(l.Quantity) * l.Price
// 		show = append(show, l)
// 	}
// 	c.JSON(200, gin.H{
// 		"Products":     show,
// 		"Total Amount": total,
// 	})
// }

// func SortProduct(c *gin.Context) {
// 	type products struct {
// 		Name        string `json:"name"`
// 		Price       uint   `json:"price"`
// 		Color       string `json:"color"`
// 		Quantity    uint   `json:"quantity"`
// 		Description string `json:"description"`
// 		Status      string `json:"status"`
// 		Image1      string
// 	}
// 	var req struct {
// 		Sort string `json:"sort"`
// 	}
// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse JSON"})
// 		return
// 	}
// 	sort := strings.ToLower(req.Sort)
// 	var product []products
// 	switch sort {
// 	case "asc":
// 		database.DB.Order("name asc").Find(&product)
// 	case "desc":
// 		database.DB.Order("name desc").Find(&product)
// 	case "htol":
// 		database.DB.Order("price desc").Find(&product)
// 	case "ltoh":
// 		database.DB.Order("price asc").Find(&product)
// 	case "latest":
// 		database.DB.Order("created_at desc").Find(&product)
// 	default:
// 		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please give correct options"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"Products": product})
// }

// -------------------------------order management---------------------------------//
// --------------------------------------------------------------------------------//

// func CartCheckOut(c *gin.Context) {
// 	var req struct {
// 		Coupon    string `json:"coupon"`
// 		Payment   string `json:"payment"`
// 		AddressId uint   `json:"addressid"`
// 	}
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(404, gin.H{"error": "failed to bind json"})
// 		return
// 	}
// 	userId := c.GetUint("userid")
// 	// userId, ok := c.MustGet("user").(uint)
// 	// if !ok {
// 	// 	c.JSON(501, gin.H{"error": "User ID not found in the context "})
// 	// }
// 	// aerr := database.DB.Where("address_id=? AND user_id=?", addId, userId).Updates(&addressDetails)
// 	var address model.Address
// 	var cartItems []model.Cart
// 	database.DB.First(&address, "address_id=?", uint(req.AddressId)).Where("User_Id = ? ", userId)
// 	database.DB.Preload("Product").Find(&cartItems, "user_id=?", userId)
// 	if req.AddressId == 0 || address.UserId != userId {
// 		c.JSON(404, gin.H{"Error": "No Address Found!"})
// 		return
// 	}
// 	if len(cartItems) == 0 {
// 		c.JSON(404, gin.H{"error": "Your cart is empty"})
// 		return
// 	}
// 	var coupon model.Coupons
// 	var order model.Orders
// 	var subTotal int
// 	for _, item := range cartItems {
// 		subTotal += int(item.Quantity) * int(item.Product.Price)
// 	}
// 	// fmt.Println("subtotal:",subTotal)
// 	fetchCoupon := database.DB.First(&coupon, "code=?", req.Coupon)
// 	order.UserId = userId
// 	order.AddressId = req.AddressId
// 	order.Total = subTotal
// 	order.Amount = subTotal - (order.Total * coupon.Value / 100)
// 	if coupon.Min < subTotal && fetchCoupon.Error == nil {
// 		order.CouponId = coupon.Id
// 	} else if coupon.Min > subTotal {
// 		c.JSON(401, gin.H{"message": "this coupon is not for this amount"})
// 		return
// 	} else if fetchCoupon.Error != nil {
// 		if req.Coupon == "" {
// 			order.CouponId = 1
// 		} else {
// 			c.JSON(404, gin.H{"message": "coupon code is not valid"})
// 			return
// 		}
// 	}
// 	num := helper.GenerateInt()
// 	order.Id, _ = strconv.Atoi(num)
// 	database.DB.Create(&order)
// 	for _, list := range cartItems {
// 		orderitem := model.OrderItem{
// 			OrderId:   uint(order.Id),
// 			ProductId: list.ProductId,
// 			Quantity:  list.Quantity,
// 			Status:    "pending",
// 		}
// 		if err := database.DB.Create(&orderitem); err.Error != nil {
// 			c.JSON(403, gin.H{"error": "failed to place the order,please try again"})
// 			return
// 		}
// 		list.Product.Quantity -= int(list.Quantity)
// 		database.DB.Model(list.Product).Update("quantity", list.Product.Quantity)
// 	}
// 	payment := model.Payment{
// 		OrderId: uint(order.Id),
// 		UserId:  userId,
// 		Amount:  order.Amount,
// 		Status:  "pending",
// 	}
// 	switch req.Payment {
// 	case "COD":
// 		// Handle Cash on Delivery
// 		payment.PayMeth = "Cash on Delivery"

// 		for _, val := range cartItems {
// 			database.DB.Delete(&val)
// 		}
// 		if err := database.DB.Create(&payment); err.Error != nil {
// 			c.JSON(403, gin.H{"Error": "Failed to process COD!! Try again later."})
// 			return
// 		}

// 		c.JSON(200, gin.H{"Message": "Sccessfully placed order"})

// 	case "PAY NOW":
// 		// Handle online payment
// 		payment.PayMeth = "Razor Pay"
// 		razorId, err := helper.ExecuteRazorpay(num, payment.Amount)

// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment gataway is not initiated"})
// 			return
// 		}
// 		payment.PayId = razorId
// 		if error := database.DB.Create(&payment); error.Error != nil {
// 			c.JSON(403, gin.H{"Message": "Payment creation failed, Try again later..."})
// 			return
// 		}

// 		c.JSON(200, gin.H{
// 			"message": "Complete the payment and place order",
// 			"data":    gin.H{"paymentId": razorId},
// 		})
// 	default:
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
// 		return
// 	}
// }
