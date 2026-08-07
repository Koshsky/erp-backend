package dto

type CreateResourceRequest struct {
	Code  string `json:"code"  example:"М"`
	Title string `json:"title" example:"Монтажник"`
}

type UpdateResourceRequest struct {
	Code  *string `json:"code"  example:"М"`
	Title *string `json:"title" example:"Монтажник"`
}
