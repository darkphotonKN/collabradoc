package rating

type CreateRatingReq struct {
	Value      uint `json:"value"`
	DocumentID uint `json:"documentId"`
}
