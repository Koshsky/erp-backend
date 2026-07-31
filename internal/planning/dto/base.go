package dto

import "time"

type Project struct {
	ID        int64     `json:"id"           example:"1"`
	Code      string    `json:"project_code" example:"КО_001"`
	StartDate time.Time `json:"start_date"   example:"2026-01-01T00:00:00Z"`
	EndDate   time.Time `json:"end_date"     example:"2026-02-01T00:00:00Z"`
	OwnerID   int64     `json:"owner_id"     example:"1"`
	Priority  int       `json:"priority"     example:"2"`
}

type Process struct {
	ID        int64     `json:"id"         example:"1"`
	Title     string    `json:"title"      example:"Инсталляция"`
	StartDate time.Time `json:"start_date" example:"2026-01-01T00:00:00Z"`
	EndDate   time.Time `json:"end_date"   example:"2026-02-01T00:00:00Z"`
	OwnerID   int64     `json:"owner_id"   example:"1"`
	ProjectID int64     `json:"project_id" example:"1"`
}

type Task struct {
	ID        int64     `json:"id"         example:"1"`
	Title     string    `json:"title"      example:"Пуско-наладочные работы"`
	StartDate time.Time `json:"start_date" example:"2026-01-01T00:00:00Z"`
	EndDate   time.Time `json:"end_date"   example:"2026-02-01T00:00:00Z"`
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
	Title     string    `json:"title"                                        examaple:"Начало работ"`
	Content   string    `json:"content"    example:"Начало работ по проекту"`
	Date      time.Time `json:"date"       example:"2026-01-01T00:00:00Z"`
	ProcessID int64     `json:"process_id" example:"1"`
}

type Assignment struct {
	ID         int64 `json:"id"          example:"1"`
	TaskID     int64 `json:"task_id"     example:"1"`
	ResourceID int64 `json:"resource_id" example:"1"`
	Quantity   int   `json:"quantity"    example:"1"`
}
