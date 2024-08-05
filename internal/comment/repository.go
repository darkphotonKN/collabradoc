package comment

import (
	"fmt"

	"github.com/darkphotonKN/collabradoc/internal/db"
	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

func CreateComment(comment CreateCommentReq, author model.User) (model.Comment, error) {
	db := db.DBCon

	newComment := model.Comment{
		Author:     author.Name,
		Comment:    comment.Comment,
		UserId:     author.ID,  // user commenting on document
		DocumentId: comment.ID, // document commented on
	}

	result := db.Create(&newComment)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return newComment, result.Error
	}

	return newComment, nil
}