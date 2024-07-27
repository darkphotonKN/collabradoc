package document

import "github.com/darkphotonKN/collabradoc/internal/db"

func CreateDocument(doc CreateDocumentReq) (Document, error) {
	db := db.DBCon

	newDoc := Document{
		OwnerId: doc.ID,
		Title:   doc.Title,
		Content: doc.Content,
	}

	result := db.Create(&newDoc)

	if result.Error != nil {
		return newDoc, result.Error
	}

	return newDoc, nil
}
