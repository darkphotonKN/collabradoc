package document

import (
	"fmt"

	"github.com/darkphotonKN/collabradoc/internal/db"
	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

func CreateDocument(doc CreateDocumentReq, userId uint) (model.Document, error) {
	db := db.DBCon

	newDoc := model.Document{
		UserId:  userId,
		Title:   doc.Title,
		Content: doc.Content,
	}

	result := db.Create(&newDoc)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return newDoc, result.Error
	}

	return newDoc, nil
}

func QueryDocuments(userId uint) ([]model.Document, error) {
	db := db.DBCon

	var documents []model.Document

	result := db.Where("user_id = ?", userId).Find(&documents)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return documents, result.Error
	}

	return documents, nil
}
