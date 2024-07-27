package document

import (
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	Title   string
	Content string
	OwnerId uint // custom foreign key in relation with User
}

type CreateDocumentReq struct {
	ID      uint   `json:"id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
