package dto

// AutoCreateConfig — configuration for auto-creating processes/tasks when a
// project is inserted. Used both as a request (PUT) and as a response (GET).
type AutoCreateConfig struct {
	Enabled   bool              `json:"enabled"`
	Processes []ProcessTemplate `json:"processes"`
}

type ProcessTemplate struct {
	Title   string         `json:"title"`
	OwnerID *int64         `json:"owner_id"`
	Tasks   []TaskTemplate `json:"tasks"`
}

type TaskTemplate struct {
	Title     string            `json:"title"`
	Resources []ResourceBinding `json:"resources"`
}

type ResourceBinding struct {
	ResourceID int64 `json:"resource_id"`
	Quantity   int   `json:"quantity"`
}
