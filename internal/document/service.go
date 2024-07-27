package document

func CreateDocumentService(doc CreateDocumentReq) (Document, error) {

	return CreateDocument(doc)
}

func GetDocuments() ([]Document, error) {
	return QueryDocuments()
}
