package dto

type AssignmentResponse struct {
	ID         int64 `json:"id"`
	TaskID     int64 `json:"task_id"`
	ResourceID int64 `json:"resource_id"`
	Quantity   int   `json:"quantity"`
}

type UpdateAssignmentRequest struct {
	Quantity *int `json:"quantity"`
}

type CreateAssignmentRequest struct {
	TaskID     int64 `json:"task_id"`
	ResourceID int64 `json:"resource_id"`
	Quantity   int   `json:"quantity"`
}
