package user

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUserList(w http.ResponseWriter, r *http.Request) {
	testJson, err := json.Marshal("test")
	if err != nil {
		fmt.Println(err)
	}

	w.Write(testJson)
}
