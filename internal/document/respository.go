package document

import (
	"fmt"

	"github.com/darkphotonKN/collabradoc/internal/db"
)

func CreateDocument(doc CreateDocumentReq) (Document, error) {
	db := db.DBCon

	newDoc := Document{
		OwnerId: doc.ID,
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

func QueryDocuments() ([]Document, error) {
	db := db.DBCon

	var documents []Document

	result := db.Find(&documents)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return documents, result.Error
	}

	return documents, nil
}
