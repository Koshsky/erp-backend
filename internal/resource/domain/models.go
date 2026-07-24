package domain

type Resource struct {
	ID       int64  `db:"id" json:"id"`
	Title    string `db:"title" json:"title"`
	Code     string `db:"code" json:"code"`
	Quantity int    `db:"quantity" json:"quantity"`
}
