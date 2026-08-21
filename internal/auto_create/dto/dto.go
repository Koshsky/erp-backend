package dto

// AutoCreateConfig — конфигурация автосоздания процессов/задач при вставке
// проекта. Используется и как запрос (PUT), и как ответ (GET).
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
