package dto

type CreateCommentRequest struct {
	Content string `json:"content" example:"Перенести сроки?"`
	// Ответ на другой комментарий той же задачи; пусто — корневой комментарий.
	ParentID *int64 `json:"parent_id" example:"3"`
}
