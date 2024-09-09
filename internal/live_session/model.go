package livesession

import (
	"encoding/json"
	"time"

	"github.com/darkphotonKN/collabradoc/internal/document"
)

type CreateLiveSessionReq struct {
	DocumentID json.Number `json:"documentId"`
}

type InviteLiveSessionReq struct {
	Email      string      `json:"email"`
	DocumentID json.Number `json:"documentId"`
}

type ExistingLiveSession struct {
	SessionID  string `json:"session_id"`
	DocumentID uint   `json:"document_id"`
}

type LiveSessionLink string

type LiveSessionInvites struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	Title      string    `json:"title"`
	DocumentID uint      `json:"documentId"`
	SessionID  string    `json:"sessionId"`
	IsActive   bool      `json:"isActive"`
}

type LiveSessionInvitesRes struct {
	ID                       uint      `json:"id"`
	CreatedAt                time.Time `json:"createdAt"`
	Title                    string    `json:"title"`
	IsActive                 bool      `json:"isActive"`
	document.LiveSessionInfo `json:"liveSession"`
}
