package rating

import model "github.com/darkphotonKN/collabradoc/internal/shared"

func GetRatingsService(documentId uint) ([]model.Rating, error) {

	return QueryRatingsByDocId(documentId)
}

func CreateRatingsService(documentId uint, value uint) (model.Rating, error) {
	return CreateRating(documentId, value)
}

/**
* Finds the ratings of a document and calculates the average.
**/
func CountRatingsAvg(documentId uint) (float32, error) {

	ratings, err := GetRatingsService(documentId)

	if err != nil {
		return 0, err
	}

	// calculate average
	var length = float32(len(ratings))
	var sum float32 = 0

	for _, rating := range ratings {
		sum += float32(rating.Value)
	}

	return sum / length, nil

}
