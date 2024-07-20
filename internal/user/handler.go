package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	testJson, err := json.Marshal(users)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(testJson)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sign up and create user

	log.Println("Creating user with password:", user.Password)

	newUser, err := CreateUser(user.Name, user.Email, user.Password)

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
	w.Write(out)
}
