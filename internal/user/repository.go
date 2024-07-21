package user

import (
	"github.com/darkphotonKN/collabradoc/internal/db"
)

// Queries for all Current Users in the DB
func QueryAllUsers() ([]User, error) {
	db := db.DBCon

	var users []User

	result := db.Find(&users)

	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

// Creates a Single User
func CreateUser(name string, email string, password string) (User, error) {
	db := db.DBCon

	newUser := User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	result := db.Create(&newUser)

	if result.Error != nil {
		return newUser, result.Error
	}

	return newUser, nil
}

// Queries for a Single User based on ID from the DB
func FindUser(email string) (User, error) {
	db := db.DBCon

	var user User
	result := db.First(&user, "email = ?", email)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
