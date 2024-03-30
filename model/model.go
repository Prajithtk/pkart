package model

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password"`
}

type Users struct {
	gorm.Model
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"password"`
	Phone    string    `json:"phone" gorm:"unique"`
	Status   string    `json:"status"`
	Gender   string    `json:"gender"`
	Wallet   uint      `json:"wallet" gorm:"default:0"`
	Address  []Address `json:"address" gorm:"foreignKey:UserId"`
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

type Products struct {
	gorm.Model
	Name        string   `json:"name"`
	Price       uint     `json:"price"`
	Color       string   `json:"color"`
	Quantity    uint     `json:"quantity"`
	Description string   `json:"description"`
	CategoryId  uint     `json:"categoryid"`
	Status      string   `json:"status"`
	Category    Category `json:"category"`
	Image1      string
	Image2      string
	Image3      string
}

type Category struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
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
	Quantity  uint `json:"quantity"`
}

type Orders struct {
	Id        uint `gorm:"primaryKey"`
	UserId    uint
	User      Users
	AddressId uint
	CouponId  uint
	Coupon    Coupons
	Total     int
	Amount    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	Id        uint `gormm:"primaryKey"`
	OrderId   uint
	Order     Orders
	ProductId uint
	Product   Products
	Quantity  int
	Status    string
}

type Coupons struct {
	Id        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not mull; unique" json:"name"`
	Desc      string `gorm:"not mull" json:"desc"`
	Code      string `gorm:"not null" json:"code"`
	Condition int
	Value     int `gorm:"not null" json:"value"`
	Exp       time.Time
}

// type Payment struct {
// 	Id        uint
// 	OrderId   uint
// 	UserId    uint
// 	Amount    uint
// 	Status    bool
// 	PayMeth   string
// 	TransId   uint
// 	CreatedAt time.Time
// }

// type Wishlist struct {
// 	Id        uint
// 	ProductId uint
// 	UserId    uint
// 	Quantity  uint
// }

// type Banner struct {
// 	Id  uint
// 	Url string
// 	Img uint
// }
