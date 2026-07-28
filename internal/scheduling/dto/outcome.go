package dto

import (
	"time"

	"github.com/Koshsky/erp-backend/internal/scheduling/domain"
)

/*
Простой просмотр всех проектов
*/
type ProjectScheduling struct {
	Projects []domain.Project `json:"projects"`
}

/*
Просмотр всех процессов и связанных проектов
Проекты важны - code, start_date, end_date использует фронтенд для отображения
*/
type ProcessScheduling struct {
	Projects  []domain.Project           `json:"projects"`
	Processes map[int64][]domain.Process `json:"processes"` // key: project_id
}

/*
Детализированные процессы (процессы + таски + ассигны + ресурсы)

Получаем все таски владельца процессов, order by project.priority ASC
Projects в ответе для кода.
TODO: после денормализации БД можно убрать projects

Process для start_date, end_date

Resources - все ресурсы, включая неиспользуемые.
*/
type TaskScheduling struct {
	Projects    map[int64]domain.Project      `json:"projects"` // key: project_id
	Processes   []domain.Process              `json:"processes"`
	Milestones  map[int64][]domain.Milestone  `json:"milestones"`  // key: process_id
	Tasks       map[int64][]domain.Task       `json:"tasks"`       // key: process_id
	Assignments map[int64][]domain.Assignment `json:"assignments"` // key: task_id
	Resources   []domain.Resource             `json:"resources"`

	Timeline struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		TotalDays int       `json:"total_days"`
	} `json:"timeline"`
}
