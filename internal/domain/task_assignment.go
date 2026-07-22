package domain

type Assignment struct {
	ID         int64 `db:"id" json:"id"`
	TaskID     int64 `db:"task_id" json:"task_id"`
	ResourceID int64 `db:"resource_id" json:"resource_id"`
	Quantity   int   `db:"quantity" json:"quantity"`
}
