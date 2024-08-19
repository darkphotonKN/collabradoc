package livesession

import "encoding/json"

type CreateLiveSessionReq struct {
	DocumentID json.Number `json:"documentId"`
}

type LiveSessionLink string
