package rating

import (
	"encoding/json"
	"net/http"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

func CreateRatingHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateRatingReq

	err := json.NewDecoder(r.Body).Decode(&req)

	createdRating, err := CreateRatingsService(req.DocumentID, req.Value)

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
