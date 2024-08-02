package comment

import "github.com/darkphotonKN/collabradoc/internal/user"

func CreateCommentService(createCommentReq CreateCommentReq, userId uint) (Comment, error) {
	// find if user exists
	user.FindUserByIdService(userId)

	return CreateComment(createCommentReq)
}
