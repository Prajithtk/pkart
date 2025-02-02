package database

import (
	"fmt"
	"os"
	"pkart/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbConnect() {
	var coupon model.Coupons

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	DB = db
	fmt.Println("Connected to database")

	DB.AutoMigrate(&model.Users{})
	DB.AutoMigrate(&model.Admin{})
	DB.AutoMigrate(&model.Address{})
	DB.AutoMigrate(&model.Products{})
	DB.AutoMigrate(&model.Category{})
	DB.AutoMigrate(&model.Cart{})
	DB.AutoMigrate(&model.Orders{})
	DB.AutoMigrate(&model.OrderItem{})
	DB.AutoMigrate(&model.Rating{})
	DB.AutoMigrate(&model.Coupons{})
	DB.AutoMigrate(&model.Payment{})
	DB.AutoMigrate(&model.Wishlist{})
	DB.AutoMigrate(&model.Otp{})
	DB.AutoMigrate(&model.Wallet{})

	coupon.Value = 0
	coupon.Name = "NO COUPON"
	coupon.Desc = "NO COUPON"
	coupon.Min = 0
	coupon.Code = "NO COUPON"
	if ERROR := DB.Create(&coupon); ERROR.Error != nil {
		fmt.Println("Error this coupon is already exist")
	}
}
