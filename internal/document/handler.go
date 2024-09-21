package document

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/utils/request"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func GetDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	// get user id from context via jwt
	userId, _ := request.ExtractUserID(r.Context())

	documents, err := GetDocuments(userId)

	formattedDocuments := make([]DocumentRes, len(documents))

	// map documents to response-friendly documents
	for index, doc := range documents {
		formattedDocuments[index] = DocumentRes{
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
			Privacy: doc.Privacy,
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	documentsRes := model.Response[[]DocumentRes]{
		Status:  http.StatusOK,
		Message: "Succesfully retrieved all documents.",
		Data:    formattedDocuments,
	}

	out, err := json.Marshal(documentsRes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func GetCommunityDocsHandler(w http.ResponseWriter, r *http.Request) {
	documents, err := GetCommunityDocsService()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	documentsRes := model.Response[[]DocumentRes]{
		Status:  http.StatusOK,
		Message: "Succesfully retrieved all documents.",
		Data:    documents,
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

func ToggleDocPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := request.ExtractUserID(r.Context())

	documentIdParam := chi.URLParam(r, "documentId")

	documentId, err := strconv.ParseUint(documentIdParam, 10, 64)

	fmt.Println("documentId from param:", documentId)

	if err != nil {
		fmt.Println("Could not decode documentId from param into uint64.")
	}

	updatedDoc, err := ToggleDocPrivacyService(userId, uint(documentId))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// fmt.Printf("\n\nupdatedDoc: \n%+v\n\n", updatedDoc)

	out, err := json.Marshal(updatedDoc)

	if err != nil {
		fmt.Println("Could not marshal updated document into json.")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
