package dto

type ResourceResponse struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Title    string `json:"title"`
	Quantity int    `json:"quantity"`
}

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

type ResourceUsageResponse struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Code          string `json:"code"`
	UsedQuantity  int64  `json:"used_quantity"`
	Available     int64  `json:"available"`
	TotalQuantity int64  `json:"total_quantity"`
}
