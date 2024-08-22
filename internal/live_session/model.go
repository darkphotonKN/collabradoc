package livesession

import "encoding/json"

type CreateLiveSessionReq struct {
	DocumentID json.Number `json:"documentId"`
}

type InviteLiveSessionReq struct {
	Email      string      `json:"email"`
	DocumentID json.Number `json:"documentId"`
}

type LiveSessionLink string
