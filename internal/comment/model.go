package comment

type CreateCommentReq struct {
	ID      uint   `json:"id" validate:"required"`
	Comment string `json:"comment" validate:"required"`
}
