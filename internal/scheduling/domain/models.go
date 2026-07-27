package domain

import (
	"time"

	assignmentDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
	milestoneDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
	processDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	projectDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	resourceDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/domain"
	taskDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
)

type Project projectDomain.Project
type Process processDomain.Process
type Task taskDomain.Task
type Resource resourceDomain.Resource
type Assignment assignmentDomain.Assignment
type Milestone milestoneDomain.Milestone

/*
Простой просмотр всех проектов
*/
type ProjectScheduling struct {
	Projects []Project `json:"projects"`
}

/*
Просмотр всех процессов и связанных проектов
Проекты важны - code, start_date, end_date использует фронтенд для отображения
*/
type ProcessScheduling struct {
	Processes []Process         `json:"processes"`
	Projects  map[int64]Project `json:"projects"`
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
	Tasks       []Task                 `json:"tasks"`
	Projects    map[int64]Project      `json:"projects"`    // key: project_id
	Processes   map[int64]Process      `json:"processes"`   // key: process_id
	Assignments map[int64][]Assignment `json:"assignments"` // key: task_id
	Milestones  map[int64][]Milestone  `json:"milestones"`  // key: process_id
	Resources   map[int64]Resource     `json:"resources"`   // key: resource_id

	Timeline struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		TotalDays int       `json:"total_days"`
	} `json:"timeline"`
}
