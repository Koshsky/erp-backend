package dto

type CreateResourceRequest struct {
	Code     string `json:"code"     example:"М"`
	Title    string `json:"title"    example:"Монтажник"`
	Quantity int    `json:"quantity" example:"7"`
}

type UpdateResourceRequest struct {
	Code     *string `json:"code"     example:"М"`
	Title    *string `json:"title"    example:"Монтажник"`
	Quantity *int    `json:"quantity" example:"7"`
}
