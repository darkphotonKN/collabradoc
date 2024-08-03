package model

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	Title   string
	Content string
	OwnerId uint      // custom foreign key in relation with User
	Comment []Comment `gorm:"foreignKey:OwnerId"`
}

type Comment struct {
	gorm.Model
	Comment string
	Author  string
	OwnerId uint // custom foreign key that relates to its parent Document
}

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Doc      []Document `gorm:"foreignKey:OwnerId"`
}

// Structuring Consistent Response Type
type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
