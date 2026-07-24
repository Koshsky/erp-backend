package domain

import "time"

type Process struct {
	ID        int64     `db:"id" json:"id"`
	OwnerID   int64     `db:"owner_id" json:"owner_id"`
	ProjectID int64     `db:"project_id" json:"project_id"`
	Title     string    `db:"title" json:"title"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
}