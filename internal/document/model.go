package document

import (
	"time"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

const (
	private = "private"
	public  = "public"
)

type CreateDocumentReq struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type LiveSessionInfo struct {
	SessionID string `json:"sessionId"`
}

type DocumentRes struct {
	ID              uint            `json:"id"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	Title           string          `json:"title"`
	Content         string          `json:"content"`
	UserId          uint            `json:"userId"`
	LiveSessionInfo LiveSessionInfo `json:"liveSession"`
	Comment         []model.Comment `json:"comment,omitempty"`
	Privacy         string          `json:"privacy"`
	AverageRating   float32         `json:"rating"`
}
