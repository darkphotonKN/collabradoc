package user

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := GetUser()

	if err != nil {
		fmt.Println("Error when attempting to fetch all users.")
	}

	testJson, err := json.Marshal(users)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(testJson)
}
