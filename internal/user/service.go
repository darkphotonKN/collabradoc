package user

import (
	"fmt"

	"github.com/darkphotonKN/collabradoc/internal/customerrors"
	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/utils/auth"
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

// User Sign Up
func SignUp(name string, email string, password string) (model.User, error) {
	// check if user already exists
	existingUser, err := FindUser(email)

	// if user already exists, throw an error
	if err == nil {
		// return existing user and the err
		return existingUser, customerrors.UserExistsErr
	}

	hashedPassword, err := auth.HashPassword(password)

	if err != nil {
		fmt.Println("Error occured while hashing password:", err)

		return model.User{}, err
	}
	return CreateUser(name, email, string(hashedPassword))
}

// User Login
func LoginUser(userLoginReq UserLoginRequest) (model.User, error) {

	// find user in database
	user, err := FindUser(userLoginReq.Email)

	if err != nil {
		return user, err
	}

	// user is now found, check password match
	authenticated := auth.CheckPasswordHash(userLoginReq.Password, user.Password)

	if !authenticated {
		return user, customerrors.PasswordIncorrectErr

	}
	return user, nil
}

func FindUserByIdService(id uint) (model.User, error) {
	return FindUserById(id)
}
