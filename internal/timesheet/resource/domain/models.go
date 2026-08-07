package domain

type Resource struct {
	ID             int64
	Title          string
	Code           string
	OwnerID        *int64
	EmployeesCount int
}
