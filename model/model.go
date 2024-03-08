package model

import (
	"time"
	"gorm.io/gorm"
)

type Admin struct {
	Name     string `json:"name"`
	Email    string `gorm:"not null" json:"email"`
	Password string `json:"password"`
}

type Users struct {
	gorm.Model
	Name      string `json:"name"`
	Email     string `gorm:"unique" json:"email"`
	Password  string `json:"password"`
	Phone     uint   `gorm:"unique" json:"phone"`
	Addressid uint   `json:"addressid"`
	Status    string `json:"status"`
	Gender    string `json:"gender"`
}

// type Address struct {
// 	Id       uint
// 	Name     string
// 	Phone    uint
// 	PinCode  uint
// 	City     string
// 	State    string
// 	Landmark string
// 	Address  string
// }

type Products struct {
	gorm.Model
	Name        string `json:"name"`
	Price       uint   `json:"price"`
	Color       string `json:"color"`
	Quantity    uint   `json:"quantity"`
	Description string `json:"description"`
	CategoryId  uint   `json:"categoryid"`
	Status      string `json:"status"`
	Category    Category
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
	Id      uint
	Otp     string
	Email   string
	Expires time.Time
}

// type Orders struct {
// 	Id        uint
// 	UserId    uint
// 	ProductId uint
// 	CouponId  uint
// 	Amount    uint
// 	Status    string
// }

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

// type Coupons struct {
// 	Id    uint
// 	Name  string
// 	Desc  string
// 	Code  string
// 	Value string
// 	Start time.Time
// 	Exp   time.Time
// }

// type Cart struct {
// 	Id        uint
// 	UserId    uint
// 	ProductId uint
// 	Quantity  uint
// 	Subtotal  uint
// }

// type Banner struct {
// 	Id  uint
// 	Url string
// 	Img uint
// }
