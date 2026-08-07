package dto

type CreateStateRequest struct {
	Code        string `json:"code"         example:"ОТП"`
	Name        string `json:"name"         example:"Отпуск"`
	IsAvailable bool   `json:"is_available" example:"false"`
}

type UpdateStateRequest struct {
	Code        *string `json:"code"         example:"ОТП"`
	Name        *string `json:"name"         example:"Отпуск"`
	IsAvailable *bool   `json:"is_available" example:"false"`
}
