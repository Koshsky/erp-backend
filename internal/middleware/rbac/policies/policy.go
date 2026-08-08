package policies

import (
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
)

// Action is an operation on an entity.
type Action int

const (
	ActionView Action = iota
	ActionCreate
	ActionUpdate
	ActionDelete
)

// Scope is the ownership zone required to perform an action.
type Scope int

const (
	ScopeNone Scope = iota
	ScopeAll
	ScopeProjectOwner // project owner (project.owner_id == user)
	ScopeProcessOwner // process owner (process.owner_id == user)
	ScopeAnyOwner     // project or process owner
	ScopeOwned        // row owner (resource.owner_id / employee.manager_id == user)
)

// rule binds a role and the required access scope.
type rule struct {
	role  string
	scope Scope
}

// matrix is the target permission matrix. admin and worker are not listed
// explicitly: admin is ScopeAll (full access), worker is ScopeNone.
//
//nolint:gochecknoglobals // rule registry
var matrix = map[rbac.Resource]map[Action][]rule{
	rbac.ResourceProject: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeProjectOwner},
		},
		ActionCreate: {
			// rp creates a project into their own ownership (owner forced to self).
			{userdomain.ProjectManager, ScopeOwned},
		},
		ActionUpdate: {
			// dp and admin edit any project, rp — their own (including priority).
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeProjectOwner},
		},
		ActionDelete: {
			{userdomain.ProjectManager, ScopeProjectOwner},
		},
	},
	rbac.ResourceProcess: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeProjectOwner},
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionCreate: {
			{userdomain.ProjectManager, ScopeProjectOwner},
		},
		ActionUpdate: {
			{userdomain.ProjectManager, ScopeProjectOwner},
		},
		ActionDelete: {
			{userdomain.ProjectManager, ScopeProjectOwner},
		},
	},
	rbac.ResourceTask: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeProjectOwner},
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
	},
	rbac.ResourceMilestone: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeProjectOwner},
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
	},
	rbac.ResourceAssignment: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeProjectOwner},
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeProcessOwner},
		},
	},
	// === Timesheet ===
	// States: vp can view the reference (for the timesheet), only admin manages.
	rbac.ResourceState: {
		ActionView: {
			{userdomain.ProcessOwner, ScopeAll},
		},
		ActionCreate: {},
		ActionUpdate: {},
		ActionDelete: {},
	},
	// Resource categories: admin — all, vp — own; vp creates into own ownership.
	rbac.ResourceResource: {
		ActionView: {
			{userdomain.ProcessOwner, ScopeOwned},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeOwned},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeOwned},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeOwned},
		},
	},
	// Employees: admin — all, vp — their subordinates (manager_id); vp creates into own team.
	rbac.ResourceEmployee: {
		ActionView: {
			{userdomain.ProcessOwner, ScopeOwned},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeOwned},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeOwned},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeOwned},
		},
	},
}

// scopeFor returns the required access scope for (role, resource, action).
// admin gets ScopeAll (full access), unlisted roles get ScopeNone.
func scopeFor(role string, res rbac.Resource, act Action) Scope {
	if role == userdomain.Admin {
		return ScopeAll
	}
	rules, ok := matrix[res][act]
	if !ok {
		return ScopeNone
	}
	for _, r := range rules {
		if r.role == role {
			return r.scope
		}
	}
	return ScopeNone
}

// Authorize reports whether the user may perform the action on the entity
// given its owners.
func Authorize(role string, res rbac.Resource, act Action, owners rbac.Owners, userID int64) bool {
	switch scopeFor(role, res, act) {
	case ScopeAll:
		return true
	case ScopeNone:
		return false
	case ScopeProjectOwner:
		return owners.ProjectOwner != 0 && owners.ProjectOwner == userID
	case ScopeProcessOwner:
		return owners.ProcessOwner != 0 && owners.ProcessOwner == userID
	case ScopeAnyOwner:
		return (owners.ProjectOwner != 0 && owners.ProjectOwner == userID) ||
			(owners.ProcessOwner != 0 && owners.ProcessOwner == userID)
	case ScopeOwned:
		// Owned: the entity belongs to the user. Timesheet rows set Owner;
		// a project being created sets ProjectOwner.
		return userID != 0 &&
			(owners.Owner == userID || owners.ProjectOwner == userID || owners.ProcessOwner == userID)
	default:
		return false
	}
}

// Can reports whether the role can perform the action at all
// (used for a coarse check before loading lists).
func Can(role string, res rbac.Resource, act Action) bool {
	return scopeFor(role, res, act) != ScopeNone
}
