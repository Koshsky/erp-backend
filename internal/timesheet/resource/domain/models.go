package domain

import "time"

type Resource struct {
	ID             int64
	Title          string
	Code           string
	Color          *string
	OwnerID        *int64
	EmployeesCount int
}

// ResourceMember is a user attached to a resource.
type ResourceMember struct {
	ID              int64
	Name            string
	Preset          *string
	Position        string
	ManagerID       *int64
	HireDate        *time.Time
	TerminationDate *time.Time
}

// ResourceAbsence is a member's absence range with the state reason.
type ResourceAbsence struct {
	UserID    int64
	UserName  string
	StateID   int64
	StateCode string
	StateName string
	StartDate time.Time
	EndDate   time.Time
}
