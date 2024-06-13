package model

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Phone    string `json:"phone" gorm:"unique"`
	Status   string `json:"status"`
	Gender   string `json:"gender"`
	// Wallet      int    `json:"wallet" gorm:"default:0"` //not
	ReferalCode string `json:"referalcode"`
	// ReferalsId	[]uint
	Address []Address `json:"address" gorm:"foreignKey:UserId"`
	
	// IsBlocked bool   `json:"isblocked" gorm:"default:false"`
}

type Address struct {
	AddressId    uint   `json:"addressid" gorm:"primaryKey"`
	BuildingName string `json:"buildingname"`
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	Landmark     string `json:"landmark"`
	PinCode      string `json:"pincode"`
	UserId       uint   `json:"userid"`
}

type Otp struct {
	Id      uint   `json:"id"`
	Otp     string `json:"otp"`
	Email   string `json:"email"`
	Expires time.Time
}

type Cart struct {
	Id        uint `json:"id"`
	UserId    uint `json:"userid"`
	ProductId uint `json:"productid"`
	Product   Products
	Quantity  int `json:"quantity"`
}

type Orders struct {
	Id             int `gorm:"primaryKey"`
	UserId         uint
	User           Users
	AddressId      uint
	Address        Address
	CouponCode     string `json:"orderCoupon"`
	CouponId       uint
	Coupon         Coupons
	ShippingCharge int
	Total          int
	Amount         int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type OrderItem struct {
	Id           uint `gormm:"primaryKey"`
	OrderId      uint
	Order        Orders
	ProductId    uint
	Product      Products
	Quantity     int
	SubTotal     float64
	Amount       float64
	Status       string
	CancelReason string
}

type Rating struct {
	Id        uint    `gorm:"primaryKey"`
	Rating    float32 `json:"rating"`
	Review    string  `json:"review"`
	UserId    uint
	ProductId uint
	Product   Products
}

type Payment struct {
	gorm.Model
	OrderId uint
	Order   Orders
	UserId  uint
	User    Users
	Amount  int
	Status  string
	PayMeth string
	PayId   string
	TransId string
}

type Wallet struct {
	Id     uint `gorm:"primaryKey"`
	UserId uint `gorm:"unique"`
	Amount int
}

type Wishlist struct {
	Id        uint
	ProductId uint
	Product   Products
	UserId    uint
	User      Users
}
