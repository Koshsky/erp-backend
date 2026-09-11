package domain

import "time"

// Execution status catalog (fixed, 3 states — see the tasks.status CHECK).
const (
	StatusNotStarted = "not_started"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

type Task struct {
	ID        int64 `json:"id"`
	ProcessID int64 `json:"process_id"`
	// ParentID links a subtask (operation) to its parent task. NULL — the task
	// is top-level in its process. Immutable after creation.
	ParentID  *int64    `json:"parent_id"`
	OwnerID   *int64    `json:"owner_id"`
	Title     string    `json:"title"`
	Color     *string   `json:"color"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	SortOrder int       `json:"order"`
}
