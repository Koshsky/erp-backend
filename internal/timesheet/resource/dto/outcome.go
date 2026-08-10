package dto

type ResourceResponse struct {
	ID             int64  `json:"id"              example:"1"`
	Code           string `json:"code"            example:"М"`
	Title          string `json:"title"           example:"Монтажник"`
	OwnerID        *int64 `json:"owner_id"        example:"3"`
	EmployeesCount int    `json:"employees_count" example:"4"`
}
