package dto

type CreateResourceRequest struct {
	Code    string `json:"code"     example:"М"`
	Title   string `json:"title"    example:"Монтажник"`
	OwnerID *int64 `json:"owner_id" example:"3"`
}

type UpdateResourceRequest struct {
	Code    *string `json:"code"     example:"М"`
	Title   *string `json:"title"    example:"Монтажник"`
	OwnerID *int64  `json:"owner_id" example:"3"`
}
