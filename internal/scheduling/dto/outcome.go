package dto

import (
	"time"
)

type ProjectScheduling struct {
	Projects []Project `json:"projects"`
}

type ProcessScheduling struct {
	Projects  []Project           `json:"projects"`
	Processes map[int64][]Process `json:"processes" ` // key: project_id
}

type TaskScheduling struct {
	Projects    map[int64]Project      `json:"projects"` // key: project_id
	Processes   []Process              `json:"processes"`
	Milestones  map[int64][]Milestone  `json:"milestones"`  // key: process_id
	Tasks       map[int64][]Task       `json:"tasks"`       // key: process_id
	Assignments map[int64][]Assignment `json:"assignments"` // key: task_id
	Resources   []Resource             `json:"resources"`

	Timeline struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		TotalDays int       `json:"total_days"`
	} `json:"timeline"`
}
