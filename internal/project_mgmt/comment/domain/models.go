package domain

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	AuthorID  int64     `json:"author_id"`
	ParentID  *int64    `json:"parent_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
