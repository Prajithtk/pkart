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

type Products struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null; unique"`
	Price       int    `json:"price"`
	Offer       int
	Color       string `json:"color"`
	Quantity    int    `json:"quantity"`
	Description string `json:"description"`
	CategoryId  uint   `json:"categoryid"`
	Status      string `json:"status"`
	AvrgRating  float32
	Category    Category ` gorm:"foreignKey:CategoryId"`
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

type Coupons struct {
	Id    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not mull; unique" json:"name"`
	Desc  string `gorm:"not mull" json:"desc"`
	Code  string `gorm:"not null" json:"code"`
	Value int    `gorm:"not null" json:"value"`
	Min   int
	Exp   time.Time
}

// type Offer struct {
// 	Id           uint
// 	ProductId    int       `json:"productid"`
// 	SpecialOffer string    `json:"offer"`
// 	Discount     float64   `json:"discount"`
// 	ValidFrom    time.Time `json:"valid_from"`
// 	ValidTo      time.Time `json:"valid_to"`
// }

// type Banner struct {
// 	Id  uint
// 	Url string
// 	Img uint
// }
