package document

import (
	"time"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

type CreateDocumentReq struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type DocumentRes struct {
	ID          uint              `json:"id"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	UserId      uint              `json:"userId"`
	LiveSession model.LiveSession `json:"liveSession"`
	Comment     []model.Comment   `json:"comment"`
}
