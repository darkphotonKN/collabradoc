package document

import (
	// "fmt"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

func CreateDocumentService(doc CreateDocumentReq, userId uint) (model.Document, error) {
	return CreateDocument(doc, userId)
}

func GetDocuments(userId uint) ([]model.Document, error) {
	return QueryDocuments(userId)
}

func GetCommunityDocsService() ([]model.Document, error) {
	return QueryPublicDocuments()
}

func GetDocumentById(id uint, userId uint) (model.Document, error) {
	return QueryDocumentById(id, userId)
}

/**
* Toggles document's privacy, if it was private make it public and vice versa.
**/
func ToggleDocPrivacyService(userId uint, documentId uint) (model.Document, error) {
	// confirm document ownership and acquire document
	existingDoc, err := QueryDocumentById(documentId, userId)

	if err != nil {
		return existingDoc, err
	}

	// fmt.Printf("\n\nexistingDoc: \n%+v\n\n", existingDoc)

	// update doc
	if existingDoc.Privacy == private {
		existingDoc.Privacy = public
	} else if existingDoc.Privacy == public {
		existingDoc.Privacy = private
	}

	// update doc with new privacy setting
	updatedDoc, err := UpdateDocument(existingDoc)
	if err != nil {
		return updatedDoc, err
	}

	return updatedDoc, nil
}
