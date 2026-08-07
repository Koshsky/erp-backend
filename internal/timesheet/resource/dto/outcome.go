package dto

type ResourceResponse struct {
	ID             int64  `json:"id"              example:"1"`
	Code           string `json:"code"            example:"М"`
	Title          string `json:"title"           example:"Монтажник"`
	EmployeesCount int    `json:"employees_count" example:"4"`
}
