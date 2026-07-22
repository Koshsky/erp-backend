package domain

import "time"

type Task struct {
	ID        int64     `db:"id" json:"id"`
	ProcessID int64     `db:"process_id" json:"process_id"`
	Title     string    `db:"title" json:"title"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
}