package model

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string
	UserId  uint      `gorm:"not null"` // foreign key in relation with User
	Comment []Comment `gorm:"foreignKey:DocumentId"`
}

type Comment struct {
	gorm.Model
	Comment    string `gorm:"not null"`
	Author     string `gorm:"not null"`
	DocumentId uint   `gorm:"not null"` // foreign key in relation with Document
	UserId     uint   `gorm:"not null"` //foreign key that relates to its User
}

type User struct {
	gorm.Model
	Name     string     `gorm:"not null"`
	Email    string     `gorm:"not null"`
	Password string     `gorm:"not null"`
	Doc      []Document `gorm:"foreignKey:UserId"`
	Comment  []Comment  `gorm:"foreignKey:UserId"`
}

// Structuring Consistent Response Type
type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
