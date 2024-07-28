package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/darkphotonKN/collabradoc/internal/customerrors"
	"github.com/darkphotonKN/collabradoc/internal/utils/auth"
	"github.com/go-playground/validator/v10"
)

type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := FindAllUsers()

	if err != nil {
		fmt.Println("Error when attempting to fetch all users.")
	}

	response := Response[[]UserResponse]{
		Status:  http.StatusOK,
		Message: "Success.",
		Data:    users,
	}

	out, err := json.Marshal(response)

	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validation
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sign up and create user
	newUser, err := SignUp(user.Name, user.Email, user.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// construct payload based on UserResponse type
	response := Response[UserResponse]{
		Status:  http.StatusCreated,
		Message: "Successfully created new user.",
		Data: UserResponse{
			ID:        newUser.ID,
			Name:      newUser.Name,
			Email:     newUser.Email,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
		},
	}

	out, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userLoginReq UserLoginRequest

	json.NewDecoder(r.Body).Decode(&userLoginReq)

	user, err := LoginUser(userLoginReq)

	w.Header().Set("Content-Type", "application/json")

	// if error is specifically the password incorrect error
	if err != nil {
		var status int

		switch {
		case errors.Is(err, customerrors.PasswordIncorrectErr):
			status = http.StatusUnauthorized

		default:
			status = http.StatusBadRequest
		}

		errRes := Response[error]{
			Status:  status,
			Message: err.Error(),
			Data:    err,
		}

		out, _ := json.Marshal(errRes)

		w.WriteHeader(status) // Set the HTTP status code
		w.Write(out)
		return
	}

	// construct payload based on UserResponse type

	// generate jwt based on user's Id
	jwtToken, err := auth.GenerateJWT(user.ID)

	if err != nil {
		fmt.Println("Error when attemping to generate jwt token.", jwtToken)
	}

	response := Response[UserLoginResponse]{
		Status:  http.StatusOK,
		Message: "Successfully logged in user.",
		Data: UserLoginResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			AccessToken: jwtToken,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
	}
	out, _ := json.Marshal(response)

	w.Write(out)
}
