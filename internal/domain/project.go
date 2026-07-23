package domain

import "time"

type Project struct {
	ID        int64     `db:"id" json:"id"`
	OwnerID   int64     `db:"owner_id" json:"owner_id"`
	Code      string    `db:"code" json:"code"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
	Priority  int       `db:"priority" json:"priority"`
}

type DetailedProject struct {
	Project
	Processes []DetailedProcess `json:"processes"`
}