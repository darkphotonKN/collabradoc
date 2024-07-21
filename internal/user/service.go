package user

import (
	"errors"
	"fmt"

	"github.com/darkphotonKN/collabradoc/internal/utils"
)

// Queries for all Users
func FindAllUsers() ([]UserResponse, error) {

	users, err := QueryAllUsers()

	// refactor into UserResponse slice
	var usersResponse []UserResponse

	for _, user := range users {
		serializedUser := UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		usersResponse = append(usersResponse, serializedUser)
	}

	return usersResponse, err
}

// User Sign up
func SignUp(name string, email string, password string) (User, error) {
	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		fmt.Println("Error occured while hashing password:", err)

		return User{}, err
	}
	return CreateUser(name, email, string(hashedPassword))
}

// User Login
func LoginUser(userLoginReq UserLoginRequest) (User, error) {

	// find user in database
	user, err := FindUser(userLoginReq.Email)

	if err != nil {
		return user, err
	}

	// user is now found, check password match
	authenticated := utils.CheckPasswordHash(userLoginReq.Password, user.Password)

	if !authenticated {
		return user, errors.New("Password entered was incorrect.")

	}
	return user, nil
}
