package rating

import model "github.com/darkphotonKN/collabradoc/internal/shared"

func GetRatingsService(documentId uint) ([]model.Rating, error) {

	return QueryRatingsByDocId(documentId)
}

func CreateRatingsService(documentId uint, userId uint, value uint) (model.Rating, error) {
	return CreateRating(documentId, userId, value)
}
