package controller

import (
	"net/http"
	"os"
	"pkart/database"
	"pkart/middleware"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Product model.Products
var Ida, _ = strconv.Atoi(os.Getenv("ID"))
var Email = os.Getenv("ADMIN")

const RoleAdmin = "admin"

func AdminLogin(c *gin.Context) {
	var admin model.Admin
	err := c.ShouldBindJSON(&admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}
	email := os.Getenv("ADMIN")
	password := os.Getenv("ADMIN_PASSWORD")

	if email != admin.Name || password != admin.Password {
		c.JSON(404, gin.H{"error": "Incorrect username or password"})
		return
	} else {
		middleware.JwtToken(c, uint(Ida), Email, RoleAdmin)
		c.JSON(200, gin.H{"message": "successfully loggedin"})
	}
}

//------------------------user management---------------------------//
//------------------------------------------------------------------//

func ListUsers(c *gin.Context) {
	var usersList []model.Users
	database.DB.Order("ID asc").Find(&usersList)

	for _, val := range usersList {
		c.JSON(200, gin.H{
			"id":    val.ID,
			"name":  val.Name,
			"email": val.Email,
			"phone": val.Phone,
			// "address": val.Addressid,
			"status": val.Status,
		})
	}
}
func BlockUser(c *gin.Context) {
	var user model.Users
	id := c.Param("ID")
	database.DB.First(&user, id)
	if user.Status == "blocked" {
		database.DB.Model(&user).Update("status", "active")
		c.JSON(http.StatusOK, gin.H{"message": "Unblocked User"})
	} else {
		database.DB.Model(&user).Update("status", "blocked")
		c.JSON(http.StatusOK, gin.H{"message": "Blocked User"})
	}
}

//---------------------------category management--------------------------------//
//------------------------------------------------------------------------------//

func ViewCategory(c *gin.Context) {
	var categoryList []model.Category
	database.DB.Order("ID asc").Find(&categoryList)

	for _, val := range categoryList {
		c.JSON(200, gin.H{
			"id":          val.ID,
			"name":        val.Name,
			"description": val.Description,
			"status":      val.Status,
		})
	}
}
func AddCategory(c *gin.Context) {
	var categoryinfo model.Category
	err := c.ShouldBindJSON(&categoryinfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}
	addCategory := database.DB.Create(&model.Category{
		Name:        categoryinfo.Name,
		Description: categoryinfo.Description,
		Status:      categoryinfo.Status,
	})
	if addCategory.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to add category"})
	} else {
		c.JSON(http.StatusSeeOther, gin.H{"message": "Category added successfully"})
	}
}
func EditCategory(c *gin.Context) {
	var categoryinfo model.Category
	err := c.ShouldBindJSON(&categoryinfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
	}
	id := c.Param("ID")
	cerr := database.DB.Where("id=?", id).Updates(&categoryinfo)
	if cerr.Error != nil {
		c.JSON(404, gin.H{"error": "failed to edit category"})
	}
	c.JSON(200, gin.H{"message": "successfully editted"})
}
func BlockCategory(c *gin.Context) {
	var category model.Category
	id := c.Param("ID")
	database.DB.First(&category, id)
	if category.Status == "blocked" {
		database.DB.Model(&category).Update("status", "active")
		c.JSON(http.StatusOK, gin.H{"message": "Category Active"})
	} else {
		database.DB.Model(&category).Update("status", "blocked")
		c.JSON(http.StatusOK, gin.H{"message": "Category Blocked"})
	}
}
func DeleteCategory(c *gin.Context) {
	var category model.Category
	id := c.Param("ID")
	err := database.DB.Where("id=?", id).Delete(&category)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}
	c.JSON(http.StatusSeeOther, gin.H{"message": "Category Deleted Successfully"})
}

// -----------------------------product management--------------------------------//
//--------------------------------------------------------------------------------//

func ViewProducts(c *gin.Context) {
	var productList []model.Products
	database.DB.Preload("Category").Order("ID asc").Find(&productList)

	for _, val := range productList {
		c.JSON(200, gin.H{
			"id":          val.ID,
			"name":        val.Name,
			"color":       val.Color,
			"quantity":    val.Quantity,
			"description": val.Description,
			"category":    val.Category.Name,
			"status":      val.Status,
			"images1":	   val.Image1,
			"images2":	   val.Image2,
			"images3":	   val.Image3,
		})
	}
}
func AddProducts(c *gin.Context) {
	err := c.ShouldBindJSON(&Product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}

	var products model.Products
	result := database.DB.Where("name=?", Product.Name).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product already exists!!! try edit product"})
		return
	} else {
		c.JSON(http.StatusSeeOther, gin.H{"message": "please upload images"})
	}
}
func ProductImage(c *gin.Context) {
	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch images"})
		return
	}
	files := file.File["images"]
	var imagepaths []string

	for _, val := range files {
		filepath := "./images/" + val.Filename
		if err := c.SaveUploadedFile(val, filepath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save images"})
			return
		}
		imagepaths = append(imagepaths, filepath)
	}
	Product.Image1 = imagepaths[0]
	Product.Image2 = imagepaths[1]
	Product.Image3 = imagepaths[2]

	upload := database.DB.Create(&Product)
	if upload.Error != nil {
		c.JSON(501, gin.H{"error": "Product already exist"})
		return
	}

	c.JSON(200, gin.H{"message": "Product added successfully"})
	Product = model.Products{}

}

func EditProducts(c *gin.Context) {
	var productinfo model.Products
	err := c.ShouldBindJSON(&productinfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
	}
	id := c.Param("ID")
	perr := database.DB.Where("id=?", id).Updates(&productinfo)
	if perr.Error != nil {
		c.JSON(404, gin.H{"error": "failed to edit product"})
	}
	c.JSON(200, gin.H{"message": "successfully editted"})
}
func DeleteProducts(c *gin.Context) {
	var product model.Products
	id := c.Param("ID")
	err := database.DB.Where("id=?", id).Delete(&product)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}
	c.JSON(http.StatusSeeOther, gin.H{"message": "Product Deleted Successfully"})
}
