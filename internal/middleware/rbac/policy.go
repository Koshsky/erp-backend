// Package rbac is the policy engine. The rules themselves (the role matrix,
// condition builders and route policies) live in the policies subpackage.
package rbac

// Resource is an entity access rights apply to.
type Resource int

const (
	ResourceProject Resource = iota
	ResourceProcess
	ResourceTask
	ResourceMilestone
	ResourceAssignment

	// ResourceState is an employee state reference.
	ResourceState
	// ResourceResource is a timesheet resource category (specialization).
	ResourceResource
	// ResourceWorker is a worker (user with role worker).
	ResourceWorker
	// ResourceComment is a task comment (threaded discussion on a task).
	ResourceComment
)

// Owners is an entity's chain of owners: the project owner and the process owner.
// For a project, ProcessOwner is empty; for a process, the owners of its parent
// and of itself; for a task/milestone/assignment, the owners of their process and
// project. For timesheet entities (resource/employee), the Owner field holds the
// row owner.
type Owners struct {
	ProjectOwner int64
	ProcessOwner int64
	Owner        int64
}

// SharesOwner reports whether the two entities share a common owner.
// It is used by business rules such as "a resource can only be assigned to a
// task of its own owner".
func (o Owners) SharesOwner(other Owners) bool {
	return o.shares(other.ProjectOwner) || o.shares(other.ProcessOwner) || o.shares(other.Owner)
}

// shares reports whether one of o's owners equals v.
func (o Owners) shares(v int64) bool {
	return v != 0 && (o.ProjectOwner == v || o.ProcessOwner == v || o.Owner == v)
}
