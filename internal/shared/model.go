package model

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Content     string
	UserId      uint      `gorm:"not null"` // foreign key in relation with User
	Comment     []Comment `gorm:"foreignKey:DocumentId"`
	LiveSession LiveSession
}

type LiveSession struct {
	gorm.Model
	SessionID  string `gorm:"not null;uniqueIndex" json:"session_id"`
	DocumentID uint   `json:"document_id"`
	IsActive   bool   `json:"isActive"`
	Users      []User `gorm:"many2many:live_session_users" json:"users"`
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

// Structuring Consistent Response Type Across API Requests
type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
