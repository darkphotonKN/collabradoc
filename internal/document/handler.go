package document

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func GetDocsList(w http.ResponseWriter, r *http.Request) {
	testJson, err := json.Marshal("test")
	if err != nil {
		fmt.Println(err)
	}

	w.Write(testJson)
}

func CreateDocHandler(w http.ResponseWriter, r *http.Request) {

	var req CreateDocumentReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		fmt.Errorf("Error when decoding json %s", err)
	}

	newDoc, err := CreateDocumentService(req)

	if err != nil {
		fmt.Errorf("Error when creating document%s", err)
	}

	newDocRes := Response[Document]{
		Status:  http.StatusCreated,
		Message: "Successfully created new document.",
		Data:    newDoc,
	}

	out, err := json.Marshal(newDocRes)

	if err != nil {
		fmt.Errorf("Error when encoding created document response: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
