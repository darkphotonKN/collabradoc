package document

import model "github.com/darkphotonKN/collabradoc/internal/shared"

func CreateDocumentService(doc CreateDocumentReq, userId uint) (model.Document, error) {

	return CreateDocument(doc, userId)
}

func GetDocuments(userId uint) ([]model.Document, error) {
	return QueryDocuments(userId)
}
