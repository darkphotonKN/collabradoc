package user

import "github.com/darkphotonKN/collabradoc/internal/db"

func QueryAllUsers() ([]User, error) {
	db := db.DBCon

	var users []User

	result := db.Find(&users)

	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

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
