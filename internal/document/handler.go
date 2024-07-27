package document

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func GetDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	documents, err := GetDocuments()

	if err != nil {
		fmt.Println("Error when retrieving document list.")
		return
	}

	documentsRes := Response[[]Document]{
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
	fmt.Println("Creating Document!")

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

	newDoc, err := CreateDocumentService(req)

	fmt.Println("newDoc:", newDoc)

	if err != nil {
		fmt.Printf("Error when creating document%s\n", err)
		return
	}

	newDocRes := Response[Document]{
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
