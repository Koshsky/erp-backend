package dto

import "time"

// TimelineProcessResponse — полный ответ на GET /api/v1/timeline/processes/:id
type TimelineResponse struct {
	Processes []TimelineProcessItem `json:"processes"`
	Resources []TimelineResource    `json:"resources"`
}

type TimelineProcessItem struct {
	ID          int64                   `json:"id"`
	Title       string                  `json:"title"`
	StartDate   time.Time               `json:"start_date" time_format:"2006-01-02"`
	EndDate     time.Time               `json:"end_date" time_format:"2006-01-02"`
	ProjectCode string                  `json:"project_code"`
	Tasks       []TimelineTaskWithResources `json:"tasks"`
}

type TimelineTaskWithResources struct {
	ID        int64                    `json:"id"`
	Title     string                   `json:"title"`
	StartDate time.Time                `json:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time                `json:"end_date" time_format:"2006-01-02"`
	Resources []TimelineTaskResource   `json:"resources"`
}

type TimelineTaskResource struct {
	ResourceID int64  `json:"resource_id"`
	Quantity   int    `json:"quantity"`
	Code       string `json:"code"`
}

type TimelineResource struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Title    string `json:"title"`
	Quantity int    `json:"quantity"`
}
