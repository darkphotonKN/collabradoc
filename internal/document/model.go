package document

import (
	"github.com/darkphotonKN/collabradoc/internal/comment"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	Title   string
	Content string
	OwnerId uint              // custom foreign key in relation with User
	Comment []comment.Comment `gorm:"foreignKey:OwnerId"`
}

type CreateDocumentReq struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
