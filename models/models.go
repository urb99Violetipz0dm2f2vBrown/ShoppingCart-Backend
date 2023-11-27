package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string `json:"title" gorm:"text;not null; default:null"`
	Author      string `json:"author" gorm:"text;not null; default:null"`
	Description string `json:"description" gorm:"text;not null; default:null"`
	Price       int64  `json:"price" gorm:"int;not null; default:null"`
	Genre       string `json:"genre" gorm:"text;not null; default:null"`
}

// cart has many books, BookID is the foreign key
type Cart struct {
	gorm.Model
	Books []Book `json:"books" gorm:"many2many:cart_books;"`
}
