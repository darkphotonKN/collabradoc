package document

func CreateDocumentService(doc CreateDocumentReq, userId uint) (Document, error) {

	return CreateDocument(doc, userId)
}

func GetDocuments(userId uint) ([]Document, error) {
	return QueryDocuments(userId)
}
