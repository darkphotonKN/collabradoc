package comment

import (
	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/user"
)

func CreateCommentService(createCommentReq CreateCommentReq, userId uint) (model.Comment, error) {
	// find if user exists
	user, err := user.FindUserByIdService(userId)

	// user doesn't exist
	if err != nil {
		return model.Comment{}, err
	}

	return CreateComment(createCommentReq, user)
}
