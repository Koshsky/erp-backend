package dto

type CreateResourceRequest struct {
	Code     string `json:"code"`
	Title    string `json:"title"`
	Quantity int    `json:"quantity"`
}

type UpdateResourceRequest struct {
	Code     *string `json:"code"`
	Title    *string `json:"title"`
	Quantity *int    `json:"quantity"`
}
