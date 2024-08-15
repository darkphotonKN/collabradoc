package livesession

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/utils/request"
	"github.com/go-playground/validator/v10"
)

func CreateLiveSessionHandler(w http.ResponseWriter, r *http.Request) {

	userId, _ := request.ExtractUserID(r.Context())

	var req CreateLiveSessionReq
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validation
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create live session for the specific document and user
	newLiveSession, err := CreateLiveSessionService(userId, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newLiveSessionRes := model.Response[model.LiveSession]{
		Status:  http.StatusCreated,
		Message: fmt.Sprintf("Successfully created new live session %s for user %d", newLiveSession.SessionID, userId),
		Data:    newLiveSession,
	}

	out, err := json.Marshal(newLiveSessionRes)

	if err != nil {
		fmt.Printf("Error when encoding created live session response: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func GetLiveSessionHandler(w http.ResponseWriter, r *http.Request) {

}
