package domain

import "time"

type Employee struct {
	ID              int64
	ResourceID      int64
	ResourceTitle   string
	Name            string
	Position        string
	ManagerID       *int64
	HireDate        *time.Time
	TerminationDate *time.Time
}

// EmployeeState — диапазон состояния сотрудника (только не-явка), [StartDate, EndDate].
type EmployeeState struct {
	ID          int64
	EmployeeID  int64
	StateID     int64
	StateCode   string
	StateName   string
	IsAvailable bool
	StartDate   time.Time
	EndDate     time.Time
}
