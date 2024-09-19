package rating

import (
	"fmt"

	"github.com/darkphotonKN/collabradoc/internal/db"
	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

func QueryRatingsByDocId(documentId uint) ([]model.Rating, error) {
	db := db.DBCon

	var ratings []model.Rating

	result := db.Joins("JOIN document ON document.id = ratings.document_id", documentId).Find(&ratings)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return ratings, fmt.Errorf("No existing ratings for this document.")
	}

	fmt.Printf("\nRatings Queried:\n %+v\n\n", ratings)

	return ratings, nil
}

func CreateRating(documentId uint, userId uint, value uint) (model.Rating, error) {
	db := db.DBCon

	newRating := model.Rating{
		DocumentId: documentId,
		Value:      value,
	}

	result := db.Create(&newRating)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return newRating, fmt.Errorf("No existing ratings for this document.")
	}

	fmt.Printf("\nRating Created:\n %+v\n\n", newRating)

	return newRating, nil
}
