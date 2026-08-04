package dto

import "github.com/Koshsky/erp-backend/internal/common/date"

type Project struct {
	ID        int64     `json:"id"           example:"1"`
	Code      string    `json:"project_code" example:"КО_001"`
	StartDate date.Date `json:"start_date"   example:"2026-01-01" format:"date"`
	EndDate   date.Date `json:"end_date"     example:"2026-02-01" format:"date"`
	OwnerID   *int64    `json:"owner_id"     example:"1"`
	Priority  int       `json:"priority"     example:"2"`
}

type Process struct {
	ID        int64     `json:"id"         example:"1"`
	Title     string    `json:"title"      example:"Инсталляция"`
	StartDate date.Date `json:"start_date" example:"2026-01-01"  format:"date"`
	EndDate   date.Date `json:"end_date"   example:"2026-02-01"  format:"date"`
	OwnerID   *int64    `json:"owner_id"   example:"1"`
	ProjectID int64     `json:"project_id" example:"1"`
}

type Task struct {
	ID        int64     `json:"id"         example:"1"`
	Title     string    `json:"title"      example:"Пуско-наладочные работы"`
	StartDate date.Date `json:"start_date" example:"2026-01-01"              format:"date"`
	EndDate   date.Date `json:"end_date"   example:"2026-02-01"              format:"date"`
	ProcessID int64     `json:"process_id" example:"1"`
}

type Resource struct {
	ID       int64  `json:"id"       example:"1"`
	Title    string `json:"title"    example:"Монтажник"`
	Code     string `json:"code"     example:"М"`
	Quantity int    `json:"quantity" example:"7"`
}

type Milestone struct {
	ID        int64     `json:"id"         example:"1"`
	Title     string    `json:"title"      example:"Начало работ"`
	Content   string    `json:"content"    example:"Начало работ по проекту"`
	Date      date.Date `json:"date"       example:"2026-01-01"              format:"date"`
	ProcessID int64     `json:"process_id" example:"1"`
}

type Assignment struct {
	ID         int64 `json:"id"          example:"1"`
	TaskID     int64 `json:"task_id"     example:"1"`
	ResourceID int64 `json:"resource_id" example:"1"`
	Quantity   int   `json:"quantity"    example:"1"`
}
