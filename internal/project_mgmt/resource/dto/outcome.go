package dto

type ResourceResponse struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Title    string `json:"title"`
	Quantity int    `json:"quantity"`
}
