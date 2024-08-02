package comment

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Comment string
	Author  string
	OwnerId uint // custom foreign key that relates to its parent Document
}

type CreateCommentReq struct {
	Comment string `json:"comment" validate:"required"`
}
