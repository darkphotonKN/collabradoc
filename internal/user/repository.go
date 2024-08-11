package user

import (
	"fmt"

	"github.com/darkphotonKN/collabradoc/internal/db"
	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

// Queries for all Current Users in the DB
func QueryAllUsers() ([]model.User, error) {
	db := db.DBCon

	var users []model.User

	result := db.Find(&users)

	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

// Creates a Single User
func CreateUser(name string, email string, password string) (model.User, error) {
	db := db.DBCon

	newUser := model.User{
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
func FindUser(email string) (model.User, error) {
	db := db.DBCon

	var user model.User
	result := db.First(&user, "email = ?", email)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

// Queries for a single user based on ID
func FindUserById(id uint) (model.User, error) {
	db := db.DBCon

	fmt.Println("attempt to find user with id:", id)

	var user model.User
	result := db.First(&user, "id = ?", id)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
