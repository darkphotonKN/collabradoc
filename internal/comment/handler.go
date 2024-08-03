package comment

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/utils/request"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := request.ExtractUserID(r.Context())

	var payload CreateCommentReq

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		fmt.Println("Error when attempting to decode payload:", err)
		return
	}

	// create comment
	newComment, err := CreateCommentService(payload, userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	response := model.Response[model.Comment]{
		Status:  http.StatusCreated,
		Message: fmt.Sprintf("Successfully created comment for document %d", payload.ID),
		Data:    newComment,
	}

	responseJson, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(responseJson)

}
