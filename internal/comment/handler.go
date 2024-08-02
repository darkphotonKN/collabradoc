package comment

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentReq

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		fmt.Println("Error when attempting to decode payload:", err)
		return
	}

	// create comment
	newComment, err := CreateCommentService(payload)
	fmt.Println(newComment)

	if err != nil {

		fmt.Println("Error when attempting to create comment", err)
		return
	}

}
