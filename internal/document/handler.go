package document

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/utils/request"
	"github.com/go-playground/validator/v10"
)

func GetDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	// get user id from context via jwt
	userId, _ := request.ExtractUserID(r.Context())

	documents, err := GetDocuments(userId)

	var resDocuments []DocumentRes

	// map documents to response-friendly documents
	for _, doc := range documents {
		resDocuments = append(resDocuments, DocumentRes{
			ID:        doc.ID,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
			Title:     doc.Title,
			Content:   doc.Content,
			UserId:    doc.UserId,
			LiveSessionInfo: LiveSessionInfo{
				SessionID: doc.LiveSession.SessionID,
			},
			Comment: doc.Comment,
		})
	}

	if err != nil {
		fmt.Println("Error when retrieving document list.")
		return
	}

	documentsRes := model.Response[[]DocumentRes]{
		Status:  http.StatusOK,
		Message: "Succesfully retrieved all documents.",
		Data:    resDocuments,
	}

	out, err := json.Marshal(documentsRes)

	if err != nil {
		fmt.Println("Error occured when encoding into json.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func CreateDocHandler(w http.ResponseWriter, r *http.Request) {
	// get user id from context via jwt
	userId, _ := request.ExtractUserID(r.Context())

	var req CreateDocumentReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		fmt.Printf("Error when decoding json %s\n", err)
	}

	// validation
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newDoc, err := CreateDocumentService(req, userId)

	fmt.Println("newDoc:", newDoc)

	if err != nil {
		fmt.Printf("Error when creating document%s\n", err)
		return
	}

	newDocRes := model.Response[model.Document]{
		Status:  http.StatusCreated,
		Message: "Successfully created new document.",
		Data:    newDoc,
	}

	out, err := json.Marshal(newDocRes)

	if err != nil {
		fmt.Printf("Error when encoding created document response: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
