package domain

import (
	assignmentDomain "github.com/Koshsky/erp-backend/internal/assignment/domain"
	processDomain "github.com/Koshsky/erp-backend/internal/process/domain"
	projectDomain "github.com/Koshsky/erp-backend/internal/project/domain"
	resourceDomain "github.com/Koshsky/erp-backend/internal/resource/domain"
	taskDomain "github.com/Koshsky/erp-backend/internal/task/domain"
)

type Project projectDomain.Project
type Process processDomain.Process
type Task taskDomain.Task
type Resource resourceDomain.Resource
type Assignment assignmentDomain.Assignment

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
	Projects    map[int64]Project      `json:"projects"`
	Processes   map[int64]Process      `json:"processes"`
	Assignments map[int64][]Assignment `json:"assignments"`
	Resources   map[int64]Resource     `json:"resources"`
}
