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
	Processes []domain.Process         `json:"processes"`
	Projects  map[int64]domain.Project `json:"projects"`
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
	Tasks       []domain.Task                 `json:"tasks"`
	Projects    map[int64]domain.Project      `json:"projects"`    // key: project_id
	Processes   map[int64]domain.Process      `json:"processes"`   // key: process_id
	Assignments map[int64][]domain.Assignment `json:"assignments"` // key: task_id
	Milestones  map[int64][]domain.Milestone  `json:"milestones"`  // key: process_id
	Resources   map[int64]domain.Resource     `json:"resources"`   // key: resource_id

	Timeline struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		TotalDays int       `json:"total_days"`
	} `json:"timeline"`
}
