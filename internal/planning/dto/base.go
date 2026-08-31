package dto

import "github.com/Koshsky/erp-backend/pkg/date"

type Project struct {
	ID        int64     `json:"id"           example:"1"`
	Code      string    `json:"project_code" example:"КО_001"`
	Color     *string   `json:"color"        example:"#0f83c4"`
	StartDate date.Date `json:"start_date"   example:"2026-01-01" format:"date"`
	EndDate   date.Date `json:"end_date"     example:"2026-02-01" format:"date"`
	OwnerID   *int64    `json:"owner_id"     example:"1"`
	Priority  int       `json:"priority"     example:"2"`
}

type Process struct {
	ID          int64     `json:"id"           example:"1"`
	Title       string    `json:"title"        example:"Инсталляция"`
	Color       *string   `json:"color"        example:"#0f83c4"`
	StartDate   date.Date `json:"start_date"   example:"2026-01-01"  format:"date"`
	EndDate     date.Date `json:"end_date"     example:"2026-02-01"  format:"date"`
	OwnerID     *int64    `json:"owner_id"     example:"1"`
	ProjectID   int64     `json:"project_id"   example:"1"`
	ProjectCode string    `json:"project_code" example:"КО_001"`
	// Order of the process within its project (ascending display order).
	Order int `json:"order" example:"1"`
}

type Task struct {
	ID        int64     `json:"id"         example:"1"`
	Title     string    `json:"title"      example:"Пуско-наладочные работы"`
	Color     *string   `json:"color"      example:"#0f83c4"`
	StartDate date.Date `json:"start_date" example:"2026-01-01"              format:"date"`
	EndDate   date.Date `json:"end_date"   example:"2026-02-01"              format:"date"`
	ProcessID int64     `json:"process_id" example:"1"`
	OwnerID   *int64    `json:"owner_id"   example:"1"`
	// Order of the task within its process (ascending display order).
	Order int `json:"order" example:"1"`
}

type Resource struct {
	ID           int64   `json:"id"            example:"1"`
	Title        string  `json:"title"         example:"Монтажник"`
	Code         string  `json:"code"          example:"М"`
	Color        *string `json:"color"         example:"#0f83c4"`
	Quantity     int     `json:"quantity"      example:"7"`
	AssignmentID int64   `json:"assignment_id" example:"1"`
}

type Milestone struct {
	ID        int64     `json:"id"         example:"1"`
	Title     string    `json:"title"      example:"Начало работ"`
	Content   string    `json:"content"    example:"Начало работ по проекту"`
	Color     *string   `json:"color"      example:"#0f83c4"`
	Date      date.Date `json:"date"       example:"2026-01-01"              format:"date"`
	ProcessID int64     `json:"process_id" example:"1"`
}

type Assignment struct {
	ID         int64 `json:"id"          example:"1"`
	TaskID     int64 `json:"task_id"     example:"1"`
	ResourceID int64 `json:"resource_id" example:"1"`
	Quantity   int   `json:"quantity"    example:"1"`
}
