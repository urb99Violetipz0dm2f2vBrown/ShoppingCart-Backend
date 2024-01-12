package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string `json:"title" gorm:"text;not null; default:null"`
	Author      string `json:"author" gorm:"text;not null; default:null"`
	Description string `json:"description" gorm:"text;not null; default:null"`
	Price       int64  `json:"price" gorm:"int;not null; default:null"`
	Genre       string `json:"genre" gorm:"text;not null; default:null"`
	CartID      uint   `json:"cart_id"` // Foreign key linking Book to Cart
}

// Cart model
type Cart struct {
	gorm.Model
	Books []Book `json:"books" gorm:"foreignKey:CartID;references:ID"`
}
