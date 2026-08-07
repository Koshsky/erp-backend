package domain

type State struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
}
