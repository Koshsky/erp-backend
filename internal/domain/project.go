package domain

import "time"

type Project struct {
	ID        int64     `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
	Priority  int       `db:"priority" json:"priority"`
}