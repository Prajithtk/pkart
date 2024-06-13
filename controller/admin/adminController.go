package controller

import (
	"fmt"
	"os"
	"pkart/database"
	"pkart/middleware"
	"pkart/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminPage(c *gin.Context) {
	var totalSales []model.Orders
	var totalAmount float64
	var totalOrder int
	if err := database.DB.Find(&totalSales).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": err.Error(),
			"code":    400,
		})
	}
	for _, v := range totalSales {
		totalAmount += float64(v.Amount)
		totalOrder += 1
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Welcome admin page",
		"data":    gin.H{"sales": totalAmount, "orderCount": totalOrder},
		"code":    200,
	})
}

var Product model.Products
var Ida, _ = strconv.Atoi(os.Getenv("ID"))
var Email = os.Getenv("ADMIN")

const RoleAdmin = "Admin"

func AdminLogin(c *gin.Context) {
	var admin model.Admin
	err := c.ShouldBindJSON(&admin)
	if err != nil {
		c.JSON(400, gin.H{
			"success": "false",
			"error":   "failed to bind json"})
		return
	}
	email := os.Getenv("ADMIN")
	password := os.Getenv("ADMIN_PASSWORD")

	if email != admin.Name || password != admin.Password {
		c.JSON(404, gin.H{
			"success": "false",
			"message": "Incorrect username or password"})
		return
	} else {
		token, err := middleware.JwtToken(c, uint(Ida), Email, RoleAdmin)
		if err != nil {
			c.JSON(403, gin.H{
				"success": "false",
				"message": "failed to create token"})
			return
		}
		c.SetCookie("JwtTokenAdmin", token, int((time.Hour * 100).Seconds()), "/", "pkartz.shop", false, false)
		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Successfully Logged in!",
			"Data": gin.H{
				"Token": token,
			},
		})
	}
}

//------------------------user management---------------------------//
//------------------------------------------------------------------//

func ListUsers(c *gin.Context) {
	var usersList []model.Users
	var userInfo []gin.H

	database.DB.Order("ID asc").Find(&usersList)

	for _, val := range usersList {
		userDetails := gin.H{
			"id":     val.ID,
			"name":   val.Name,
			"email":  val.Email,
			"phone":  val.Phone,
			"status": val.Status,
		}
		userInfo = append(userInfo, userDetails)
	}
	c.JSON(200, gin.H{
		"success": "true",
		"message": "the user details are:",
		"values":  userInfo})
}
func BlockUser(c *gin.Context) {
	var user model.Users
	id := c.Param("ID")
	database.DB.First(&user, id)
	if user.Status == "Blocked" {
		database.DB.Model(&user).Update("status", "Active")
		c.JSON(200, gin.H{
			"success": "true",
			"message": fmt.Sprintf("User with ID %s has been unblocked", id)})
	} else {
		database.DB.Model(&user).Update("status", "Blocked")
		c.JSON(200, gin.H{
			"success": "true",
			"message": fmt.Sprintf("User with ID %s has been blocked", id)})
	}
}

//---------------------------category management--------------------------------//
//------------------------------------------------------------------------------//

// func ViewCategory(c *gin.Context) {
// 	var categoryList []model.Category
// 	database.DB.Order("ID asc").Find(&categoryList)

// 	for _, val := range categoryList {
// 		c.JSON(200, gin.H{
// 			"id":          val.ID,
// 			"name":        val.Name,
// 			"description": val.Description,
// 			"status":      val.Status,
// 		})
// 	}
// }
// func AddCategory(c *gin.Context) {
// 	var categoryinfo model.Category
// 	err := c.ShouldBindJSON(&categoryinfo)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
// 		return
// 	}
// 	addCategory := database.DB.Create(&model.Category{
// 		Name:        categoryinfo.Name,
// 		Description: categoryinfo.Description,
// 		Status:      categoryinfo.Status,
// 	})
// 	if addCategory.Error != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to add category"})
// 	} else {
// 		c.JSON(200, gin.H{"message": "Category added successfully"})
// 	}
// }
// func EditCategory(c *gin.Context) {
// 	var categoryinfo model.Category
// 	err := c.ShouldBindJSON(&categoryinfo)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
// 	}
// 	id := c.Param("ID")
// 	cerr := database.DB.Where("id=?", id).Updates(&categoryinfo)
// 	if cerr.Error != nil {
// 		c.JSON(404, gin.H{"error": "failed to edit category"})
// 	}
// 	c.JSON(200, gin.H{"message": "successfully editted"})
// }
// func BlockCategory(c *gin.Context) {
// 	var category model.Category
// 	id := c.Param("ID")
// 	database.DB.First(&category, id)
// 	if category.Status == "blocked" {
// 		database.DB.Model(&category).Update("status", "active")
// 		c.JSON(http.StatusOK, gin.H{"message": "Category Active"})
// 	} else {
// 		database.DB.Model(&category).Update("status", "blocked")
// 		c.JSON(http.StatusOK, gin.H{"message": "Category Blocked"})
// 	}
// }
// func DeleteCategory(c *gin.Context) {
// 	var category model.Category
// 	id := c.Param("ID")
// 	err := database.DB.Where("id=?", id).Delete(&category)
// 	if err.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"message": "Category Deleted Successfully"})
// }

//------------------------------coupons-------------------------------------------//
//--------------------------------------------------------------------------------//

// type Couponst struct {
// 	model.Coupons
// 	AddDay int `json:"day"`
// }

// func AddCoupon(c *gin.Context) {

// 	var addCoupon Couponst
// 	if err := c.ShouldBindJSON(&addCoupon); err != nil {
// 		c.JSON(404, gin.H{"error": "failed to bind json"})
// 		return
// 	}
// 	var coupon model.Coupons
// 	coupon.Name = addCoupon.Name
// 	coupon.Desc = addCoupon.Desc
// 	coupon.Code = addCoupon.Code
// 	coupon.Value = addCoupon.Value
// 	coupon.Min = addCoupon.Min
// 	coupon.Exp = time.Now().AddDate(0, 0, addCoupon.AddDay)
// 	err := database.DB.Create(&coupon)
// 	if err.Error != nil {
// 		c.JSON(409, gin.H{"message": "Coupon name or code already exist, please try to edit"})
// 	} else {
// 		c.JSON(200, gin.H{"message": "Coupon added successfully"})
// 	}
// }

// func ViewCoupon(c *gin.Context) {
// 	var couponlist []model.Coupons
// 	var couponinfo []gin.H
// 	database.DB.Order("ID asc").Find(&couponlist)
// 	// fmt.Println(coupons)
// 	for _, val := range couponlist {
// 		coupondetails := gin.H{
// 			"ID":          val.Id,
// 			"NAME":        val.Name,
// 			"DESCRIPTION": val.Desc,
// 			"CODE":        val.Code,
// 			"VALUE":       val.Value,
// 			"EXPIRY":      val.Exp,
// 		}
// 		couponinfo = append(couponinfo, coupondetails)
// 	}
// 	c.JSON(200, gin.H{
// 		"status":  "True",
// 		"message": "The coupon details are :",
// 		"values":  couponinfo})
// }

// func EditCoupon(c *gin.Context) {
// 	Id, _ := strconv.Atoi(c.Param("Id"))
// 	var edit model.Coupons
// 	var coupon model.Coupons
// 	c.BindJSON(&edit)
// 	database.DB.First(&coupon, Id)
// 	database.DB.Model(&model.Coupons{}).Where("Id=?", Id).Updates(edit)
// 	if coupon.Id == 0 {
// 		c.JSON(404, gin.H{"error": "Coupon not found."})
// 	} else {
// 		c.JSON(200, gin.H{"message": "Coupon edited succesfully."})
// 	}
// }

// func DeleteCoupon(c *gin.Context) {

// 	Id, _ := strconv.Atoi(c.Param("Id"))

// 	var coupon model.Coupons

// 	database.DB.First(&coupon, Id)

// 	if coupon.Id == 0 {
// 		c.JSON(404, gin.H{
// 			"success": "False",
// 			"message": "Coupon not found.",
// 			"data":    "{ }"})
// 	} else {
// 		database.DB.Delete(&coupon)
// 		c.JSON(200, gin.H{
// 			"success": "True",
// 			"message": "Coupon deleted successfully",
// 			"data":    "{ }"})
// 	}
// }

// -----------------------------product management--------------------------------//
//--------------------------------------------------------------------------------//

// func ViewProducts(c *gin.Context) {
// 	var productList []model.Products
// 	var productinfo []gin.H
// 	database.DB.Preload("Category").Order("ID asc").Find(&productList)
// 	for _, val := range productList {
// 		productdetails := gin.H{
// 			"id":          val.ID,
// 			"name":        val.Name,
// 			"color":       val.Color,
// 			"quantity":    val.Quantity,
// 			"description": val.Description,
// 			"category":    val.Category.Name,
// 			"status":      val.Status,
// 			"images1":     val.Image1,
// 			"images2":     val.Image2,
// 			"images3":     val.Image3,
// 		}
// 		productinfo = append(productinfo, productdetails)
// 	}
// 	c.JSON(200, gin.H{
// 		"status":  "true",
// 		"message": "product details are:",
// 		"values":  productinfo})
// }

// func AddProducts(c *gin.Context) {
// 	// var Product model.Products
// 	if err := c.ShouldBindJSON(&Product); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
// 		return
// 	}
// 	// Checking if a product with the same name already exists
// 	var existingProduct model.Products
// 	if result := database.DB.Where("name=?", Product.Name).First(&existingProduct); result.Error == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": "false",
// 			"message": "product already exists!!! try edit product"})
// 		return
// 	}
// 	// If no existing product found, proceed with adding the new product
// 	c.JSON(http.StatusSeeOther, gin.H{
// 		"success": "true",
// 		"message": "please upload images"})
// }

// func ProductImage(c *gin.Context) {
// 	file, err := c.MultipartForm()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": "false",
// 			"message": "Failed to fetch images"})
// 		return
// 	}
// 	files := file.File["images"]
// 	var imagePaths []string
// 	for i, val := range files {
// 		filePath := "./images/" + strconv.Itoa(i) + "_" + val.Filename
// 		if err := c.SaveUploadedFile(val, filePath); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"success": "false",
// 				"message": "Failed to save images"})
// 			return
// 		}
// 		imagePaths = append(imagePaths, filePath)
// 	}
// 	Product.Image1 = imagePaths[0]
// 	Product.Image2 = imagePaths[1]
// 	Product.Image3 = imagePaths[2]
// 	if err := database.DB.Create(&Product).Error; err != nil {
// 		c.JSON(501, gin.H{
// 			"success": "false",
// 			"message": "Failed to add product to database"})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"success": "true",
// 		"message": "Product added successfully"})
// 	Product = model.Products{}
// }

// func EditProducts(c *gin.Context) {
// 	var productinfo model.Products
// 	if err := c.ShouldBindJSON(&productinfo); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": "false",
// 			"message": "failed to bind json"})
// 		return
// 	}
// 	id := c.Param("ID")
// 	if err := database.DB.Where("id=?", id).Updates(&productinfo); err.Error != nil {
// 		c.JSON(404, gin.H{
// 			"success": "false",
// 			"message": "failed to edit product"})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"success": "true",
// 		"message": "successfully editted"})
// }
// func DeleteProducts(c *gin.Context) {
// 	var product model.Products
// 	id := c.Param("ID")
// 	err := database.DB.Where("id=?", id).Delete(&product)
// 	if err.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": "false",
// 			"message": "failed to delete product"})
// 		return
// 	}
// 	c.JSON(http.StatusSeeOther, gin.H{
// 		"success": "true",
// 		"message": "Product Deleted Successfully"})
// }

// ------------------------------order management---------------------------------//
//--------------------------------------------------------------------------------//

// func ShowOrders(c *gin.Context) {
// 	var order []model.OrderItem
// 	var show []gin.H
// 	err := database.DB.Preload("Order").Preload("Product").Preload("Order.User").Find(&order).Error
// 	if err != nil {
// 		c.JSON(404, gin.H{"Message": "No orders found!"})
// 		return
// 	}
// 	for _, val := range order {
// 		show = append(show, gin.H{
// 			"Id":           val.Id,
// 			"OrderId":      val.OrderId,
// 			"Username":     val.Order.User.Name,
// 			"User_Email":   val.Order.User.Email,
// 			"Product_Name": val.Product.Name,
// 			"Image1":       val.Product.Image1,
// 			"Image2":       val.Product.Image2,
// 			"Image3":       val.Product.Image3,
// 			"Quantity":     val.Quantity,
// 			"Status":       val.Status,
// 		})
// 	}
// 	c.JSON(200, gin.H{
// 		"success": "true",
// 		"message": "Order items are:",
// 		"values":  show})
// }

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
