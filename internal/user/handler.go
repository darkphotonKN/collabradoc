package user

import (
	"encoding/json"
	"fmt"
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
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sign up and create user
	var response Response[User]

	newUser, err := CreateUser(user.Name, user.Email, user.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response = Response[User]{
		Status:  http.StatusCreated,
		Message: "Created New User",
		Data:    newUser,
	}

	out, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
