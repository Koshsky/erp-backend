package domain

import "time"

type Resource struct {
	ID             int64
	Title          string
	Code           string
	OwnerID        *int64
	EmployeesCount int
}

// ResourceMember is a user attached to a resource.
type ResourceMember struct {
	ID              int64
	Name            string
	Role            string
	Position        string
	ManagerID       *int64
	HireDate        *time.Time
	TerminationDate *time.Time
}
