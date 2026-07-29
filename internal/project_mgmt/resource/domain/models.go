package domain

type Resource struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}
