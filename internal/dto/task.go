package dto

import "time"

type TaskResponse struct {
	ID        int64     `json:"id"`
	ProcessID int64     `json:"process_id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time `json:"end_date" time_format:"2006-01-02"`
}

type TaskDetailResponse struct {
	TaskResponse
	Assignments []AssignmentResponse `json:"assignments"`
}

type UpdateTaskRequest struct {
	Title     *string    `json:"title"`
	StartDate *time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `json:"end_date" time_format:"2006-01-02"`
}

type CreateTaskRequest struct {
	ProcessID int64     `json:"process_id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time `json:"end_date" time_format:"2006-01-02"`
}
