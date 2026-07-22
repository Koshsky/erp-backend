package domain

import "time"

type Milestone struct {
	ID        int64     `db:"id" json:"id"`
	ProcessID int64     `db:"process_id" json:"process_id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	Date      time.Time `db:"date" json:"date"`
}
