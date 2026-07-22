package domain

type Resource struct {
	ID       int64  `db:"id" json:"id"`
	Title    string `db:"title" json:"title"`
	Code     string `db:"code" json:"code"`
	Quantity int    `db:"quantity" json:"quantity"`
}

type ResourceUsage struct {
	ID            int64  `db:"id" json:"id"`
	Title         string `db:"title" json:"title"`
	Code          string `db:"code" json:"code"`
	TotalQuantity int64  `db:"total_quantity" json:"total_quantity"`
	UsedQuantity  int64  `db:"used_quantity" json:"used_quantity"`
	Available     int64  `db:"available_quantity" json:"available_quantity"`
}
