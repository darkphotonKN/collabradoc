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

	result := db.Preload("LiveSession").Where("user_id = ?", userId).Find(&documents)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return documents, fmt.Errorf("No existing documents for this user.")
	}

	fmt.Printf("\nDocuments found:\n %+v\n\n", documents)

	return documents, nil
}

func QueryPublicDocuments() ([]model.Document, error) {
	db := db.DBCon

	var documents []model.Document

	result := db.Preload("LiveSession").Where("privacy = ?", "public").Find(&documents)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return documents, fmt.Errorf("No existing public documents.")
	}

	return documents, nil
}

func QueryDocumentById(id uint, userId uint) (model.Document, error) {
	db := db.DBCon

	var document model.Document

	result := db.Where("user_id = ? AND id = ?", userId, id).First(&document)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return model.Document{}, fmt.Errorf("Document does not belong to user you are attempting to create a live session with.")
	}

	return document, nil
}

func UpdateDocument(doc model.Document) (model.Document, error) {
	db := db.DBCon

	result := db.Save(&doc)

	if result.Error != nil {
		fmt.Printf("Error when attempting to update doc: %s\n\n", result.Error)
		return model.Document{}, result.Error
	}

	return doc, nil
}
