package document

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetDocsList(w http.ResponseWriter, r *http.Request) {
	testJson, err := json.Marshal("test")
	if err != nil {
		fmt.Println(err)
	}

	w.Write(testJson)
}
