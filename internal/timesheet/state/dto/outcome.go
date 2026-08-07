package dto

type StateResponse struct {
	ID          int64  `json:"id"           example:"4"`
	Code        string `json:"code"         example:"vacation"`
	Name        string `json:"name"         example:"Отпуск"`
	IsAvailable bool   `json:"is_available" example:"false"`
}
