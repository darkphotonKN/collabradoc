package comment

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/utils/request"
	"github.com/go-playground/validator/v10"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := request.ExtractUserID(r.Context())

	var payload CreateCommentReq

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		fmt.Println("Error when attempting to decode payload:", err)
		return
	}

	// validation
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// create comment
	newComment, err := CreateCommentService(payload, userId)

	if err != nil {

		errResponse := model.Response[model.Comment]{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    newComment,
		}

		errResponseJson, err := json.Marshal(errResponse)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(errResponseJson)
		return
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

	w.Write(responseJson)

}
