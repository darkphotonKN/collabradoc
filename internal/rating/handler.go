package rating

import (
	"encoding/json"
	"net/http"
	"strconv"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/utils/request"
	"github.com/go-chi/chi/v5"
)

func CreateRatingHandler(w http.ResponseWriter, r *http.Request) {
	// get user id from context via jwt
	userId, _ := request.ExtractUserID(r.Context())

	var req CreateRatingReq

	err := json.NewDecoder(r.Body).Decode(&req)

	documentIdParam := chi.URLParam(r, "documentId")

	documentId, err := strconv.ParseUint(documentIdParam, 64, 0)

	createdRating, err := CreateRatingsService(documentId, userId, req.Value)

	// map documents to response-friendly documents

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	documentsRes := model.Response[model.Rating]{
		Status:  http.StatusOK,
		Message: "Succesfully retrieved all documents.",
		Data:    createdRating,
	}

	out, err := json.Marshal(documentsRes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
