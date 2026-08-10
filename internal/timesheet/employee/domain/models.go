package domain

import "time"

type Employee struct {
	ID              int64
	ResourceID      int64
	Name            string
	Position        string
	ManagerID       *int64
	HireDate        *time.Time
	TerminationDate *time.Time
}

// EmployeeState is an employee state range (non-presence only), [StartDate, EndDate].
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
