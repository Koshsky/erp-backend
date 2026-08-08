package policies_test

import (
	"testing"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/middleware/rbac/policies"
)

const (
	admin  = "admin"
	dp     = "dp"
	rp     = "rp"
	vp     = "vp"
	worker = "worker"

	uAdmin = int64(1)
	uDP    = int64(2)
	uRP    = int64(3)
	uVP1   = int64(4)
	uVP2   = int64(5)
)

// Canonical owners: the project belongs to rp, the process to vp1.
//
//nolint:gochecknoglobals // rule registry
var projectOfRP = rbac.Owners{ProjectOwner: uRP}

//nolint:gochecknoglobals // rule registry
var processOfVP1 = rbac.Owners{ProjectOwner: uRP, ProcessOwner: uVP1}

func TestAuthorize_Project(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		role    string
		act     policies.Action
		owners  rbac.Owners
		userID  int64
		allowed bool
	}{
		// View
		{"admin sees any", admin, policies.ActionView, rbac.Owners{}, uAdmin, true},
		{"dp sees any", dp, policies.ActionView, rbac.Owners{}, uDP, true},
		{"rp sees own", rp, policies.ActionView, projectOfRP, uRP, true},
		{"rp does not see foreign", rp, policies.ActionView, projectOfRP, uVP1, false},
		{"vp sees none", vp, policies.ActionView, projectOfRP, uVP1, false},
		{"worker sees none", worker, policies.ActionView, projectOfRP, uAdmin, false},
		// Create
		{"admin creates", admin, policies.ActionCreate, rbac.Owners{}, uAdmin, true},
		{"rp creates own", rp, policies.ActionCreate, rbac.Owners{ProjectOwner: uRP}, uRP, true},
		{"rp cannot create for another", rp, policies.ActionCreate, rbac.Owners{ProjectOwner: uVP1}, uRP, false},
		{"dp cannot create", dp, policies.ActionCreate, rbac.Owners{}, uDP, false},
		{"vp cannot create", vp, policies.ActionCreate, rbac.Owners{}, uVP1, false},
		// Update (all fields including priority are one action)
		{"admin updates any", admin, policies.ActionUpdate, projectOfRP, uAdmin, true},
		{"dp updates any", dp, policies.ActionUpdate, projectOfRP, uDP, true},
		{"rp updates own", rp, policies.ActionUpdate, projectOfRP, uRP, true},
		{"rp updates foreign", rp, policies.ActionUpdate, projectOfRP, uVP1, false},
		{"vp cannot update", vp, policies.ActionUpdate, projectOfRP, uVP1, false},
		// Delete
		{"admin deletes", admin, policies.ActionDelete, projectOfRP, uAdmin, true},
		{"rp deletes own", rp, policies.ActionDelete, projectOfRP, uRP, true},
		{"rp deletes foreign", rp, policies.ActionDelete, projectOfRP, uVP1, false},
		{"dp cannot delete", dp, policies.ActionDelete, projectOfRP, uDP, false},
		{"vp cannot delete", vp, policies.ActionDelete, projectOfRP, uVP1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := policies.Authorize(
				tc.role,
				rbac.ResourceProject,
				tc.act,
				tc.owners,
				tc.userID,
			); got != tc.allowed {
				t.Errorf("policies.Authorize(%s, project, %v) = %v, want %v", tc.role, tc.act, got, tc.allowed)
			}
		})
	}
}

func TestAuthorize_Process(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		role    string
		act     policies.Action
		owners  rbac.Owners
		userID  int64
		allowed bool
	}{
		// View
		{"admin sees any", admin, policies.ActionView, processOfVP1, uAdmin, true},
		{"dp sees any", dp, policies.ActionView, processOfVP1, uDP, true},
		{"rp sees own project's process", rp, policies.ActionView, processOfVP1, uRP, true},
		{"rp does not see foreign process", rp, policies.ActionView, processOfVP1, uVP1, false},
		{"vp sees own process", vp, policies.ActionView, processOfVP1, uVP1, true},
		{"vp does not see foreign process", vp, policies.ActionView, processOfVP1, uVP2, false},
		{"worker sees none", worker, policies.ActionView, processOfVP1, uVP1, false},
		// Create (in own project)
		{"admin creates", admin, policies.ActionCreate, processOfVP1, uAdmin, true},
		{"rp creates in own project", rp, policies.ActionCreate, processOfVP1, uRP, true},
		{"rp cannot create in foreign project", rp, policies.ActionCreate, processOfVP1, uVP1, false},
		{"vp cannot create", vp, policies.ActionCreate, processOfVP1, uVP1, false},
		{"dp cannot create", dp, policies.ActionCreate, processOfVP1, uDP, false},
		// Update/delete
		{"rp updates own project's process", rp, policies.ActionUpdate, processOfVP1, uRP, true},
		{"vp cannot update process", vp, policies.ActionUpdate, processOfVP1, uVP1, false},
		{"dp cannot update process", dp, policies.ActionUpdate, processOfVP1, uDP, false},
		{"rp deletes own project's process", rp, policies.ActionDelete, processOfVP1, uRP, true},
		{"vp cannot delete process", vp, policies.ActionDelete, processOfVP1, uVP1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := policies.Authorize(
				tc.role,
				rbac.ResourceProcess,
				tc.act,
				tc.owners,
				tc.userID,
			); got != tc.allowed {
				t.Errorf("policies.Authorize(%s, process, %v) = %v, want %v", tc.role, tc.act, got, tc.allowed)
			}
		})
	}
}

func TestAuthorize_TaskMilestoneAssignment(t *testing.T) {
	t.Parallel()
	resources := []struct {
		name string
		res  rbac.Resource
	}{
		{"task", rbac.ResourceTask},
		{"milestone", rbac.ResourceMilestone},
		{"assignment", rbac.ResourceAssignment},
	}
	for _, r := range resources {
		t.Run(r.name, func(t *testing.T) {
			t.Parallel()
			cases := []struct {
				name    string
				role    string
				act     policies.Action
				owners  rbac.Owners
				userID  int64
				allowed bool
			}{
				// View
				{"admin sees any", admin, policies.ActionView, processOfVP1, uAdmin, true},
				{"dp sees any", dp, policies.ActionView, processOfVP1, uDP, true},
				{"rp sees own project's", rp, policies.ActionView, processOfVP1, uRP, true},
				{"rp view only foreign project", rp, policies.ActionView, processOfVP1, uVP1, false},
				{"vp sees own process's", vp, policies.ActionView, processOfVP1, uVP1, true},
				{"vp does not see foreign process's", vp, policies.ActionView, processOfVP1, uVP2, false},
				{"worker sees none", worker, policies.ActionView, processOfVP1, uVP1, false},
				// Create
				{"admin creates", admin, policies.ActionCreate, processOfVP1, uAdmin, true},
				{"vp creates in own process", vp, policies.ActionCreate, processOfVP1, uVP1, true},
				{"vp cannot create in foreign process", vp, policies.ActionCreate, processOfVP1, uVP2, false},
				{"rp cannot create", rp, policies.ActionCreate, processOfVP1, uRP, false},
				{"dp cannot create", dp, policies.ActionCreate, processOfVP1, uDP, false},
				// Update/delete
				{"admin updates", admin, policies.ActionUpdate, processOfVP1, uAdmin, true},
				{"vp updates own process's", vp, policies.ActionUpdate, processOfVP1, uVP1, true},
				{"rp cannot update (view only)", rp, policies.ActionUpdate, processOfVP1, uRP, false},
				{"dp cannot update", dp, policies.ActionUpdate, processOfVP1, uDP, false},
				{"admin deletes", admin, policies.ActionDelete, processOfVP1, uAdmin, true},
				{"vp deletes own process's", vp, policies.ActionDelete, processOfVP1, uVP1, true},
				{"rp cannot delete (view only)", rp, policies.ActionDelete, processOfVP1, uRP, false},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					if got := policies.Authorize(tc.role, r.res, tc.act, tc.owners, tc.userID); got != tc.allowed {
						t.Errorf(
							"policies.Authorize(%s, %s, %v) = %v, want %v",
							tc.role,
							r.name,
							tc.act,
							got,
							tc.allowed,
						)
					}
				})
			}
		})
	}
}

func TestOwnersSharesOwner(t *testing.T) {
	t.Parallel()
	taskOfVP1 := rbac.Owners{ProjectOwner: uRP, ProcessOwner: uVP1}
	taskOfVP2 := rbac.Owners{ProjectOwner: uRP, ProcessOwner: uVP2}
	resourceOfVP1 := rbac.Owners{Owner: uVP1}
	resourceOfVP2 := rbac.Owners{Owner: uVP2}
	noOwner := rbac.Owners{}

	cases := []struct {
		name string
		a    rbac.Owners
		b    rbac.Owners
		want bool
	}{
		{"shares process owner", taskOfVP1, resourceOfVP1, true},
		{"shares process owner (vp2)", taskOfVP2, resourceOfVP2, true},
		{"different owners", taskOfVP1, resourceOfVP2, false},
		{"shares project owner", rbac.Owners{ProjectOwner: uRP}, rbac.Owners{ProjectOwner: uRP}, true},
		{"no owners", noOwner, noOwner, false},
		{"empty chain on one side", noOwner, resourceOfVP1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.a.SharesOwner(tc.b); got != tc.want {
				t.Errorf("SharesOwner(%+v, %+v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestCan(t *testing.T) {
	t.Parallel()
	cases := []struct {
		role string
		res  rbac.Resource
		act  policies.Action
		want bool
	}{
		{rp, rbac.ResourceTask, policies.ActionUpdate, false},    // rp does not change tasks
		{vp, rbac.ResourceTask, policies.ActionUpdate, true},     // vp changes tasks of own processes
		{rp, rbac.ResourceProcess, policies.ActionCreate, true},  // rp creates processes
		{vp, rbac.ResourceProcess, policies.ActionCreate, false}, // vp does not create processes
		{dp, rbac.ResourceProject, policies.ActionUpdate, true},  // dp edits projects
		{vp, rbac.ResourceProject, policies.ActionUpdate, false}, // vp does not edit projects
		{worker, rbac.ResourceTask, policies.ActionView, false},  // worker sees nothing
	}
	for _, tc := range cases {
		if got := policies.Can(tc.role, tc.res, tc.act); got != tc.want {
			t.Errorf("policies.Can(%s, %v, %v) = %v, want %v", tc.role, tc.res, tc.act, got, tc.want)
		}
	}
}

func TestAuthorize_Timesheet(t *testing.T) {
	t.Parallel()
	resourceOfVP := rbac.Owners{Owner: uVP1}
	empOfVP := rbac.Owners{Owner: uVP1}

	cases := []struct {
		name    string
		res     rbac.Resource
		act     policies.Action
		role    string
		owners  rbac.Owners
		userID  int64
		allowed bool
	}{
		// === States ===
		{"admin views states", rbac.ResourceState, policies.ActionView, admin, rbac.Owners{}, uAdmin, true},
		{
			"vp views states (reference for the timesheet)",
			rbac.ResourceState,
			policies.ActionView,
			vp,
			rbac.Owners{},
			uVP1,
			true,
		},
		{"rp does not view states", rbac.ResourceState, policies.ActionView, rp, rbac.Owners{}, uRP, false},
		{"dp does not view states", rbac.ResourceState, policies.ActionView, dp, rbac.Owners{}, uDP, false},
		{"admin creates state", rbac.ResourceState, policies.ActionCreate, admin, rbac.Owners{}, uAdmin, true},
		{"vp cannot create state", rbac.ResourceState, policies.ActionCreate, vp, rbac.Owners{}, uVP1, false},
		{"vp cannot update state", rbac.ResourceState, policies.ActionUpdate, vp, rbac.Owners{}, uVP1, false},
		{"vp cannot delete state", rbac.ResourceState, policies.ActionDelete, vp, rbac.Owners{}, uVP1, false},
		// === Timesheet resources ===
		{"admin views all resources", rbac.ResourceResource, policies.ActionView, admin, rbac.Owners{}, uAdmin, true},
		{"vp views own resource", rbac.ResourceResource, policies.ActionView, vp, resourceOfVP, uVP1, true},
		{
			"vp does not view foreign resource",
			rbac.ResourceResource,
			policies.ActionView,
			vp,
			resourceOfVP,
			uVP2,
			false,
		},
		{"rp does not view resources", rbac.ResourceResource, policies.ActionView, rp, resourceOfVP, uRP, false},
		{"dp does not view resources", rbac.ResourceResource, policies.ActionView, dp, resourceOfVP, uDP, false},
		{"vp creates own resource", rbac.ResourceResource, policies.ActionCreate, vp, resourceOfVP, uVP1, true},
		{
			"vp cannot create resource for another",
			rbac.ResourceResource,
			policies.ActionCreate,
			vp,
			rbac.Owners{Owner: uVP2},
			uVP1,
			false,
		},
		{
			"rp cannot create resource",
			rbac.ResourceResource,
			policies.ActionCreate,
			rp,
			rbac.Owners{Owner: uRP},
			uRP,
			false,
		},
		{"vp updates own resource", rbac.ResourceResource, policies.ActionUpdate, vp, resourceOfVP, uVP1, true},
		{
			"vp does not update foreign resource",
			rbac.ResourceResource,
			policies.ActionUpdate,
			vp,
			resourceOfVP,
			uVP2,
			false,
		},
		{"vp deletes own resource", rbac.ResourceResource, policies.ActionDelete, vp, resourceOfVP, uVP1, true},
		// === Employees ===
		{"admin views all employees", rbac.ResourceEmployee, policies.ActionView, admin, rbac.Owners{}, uAdmin, true},
		{"vp views own employee", rbac.ResourceEmployee, policies.ActionView, vp, empOfVP, uVP1, true},
		{"vp does not view foreign employee", rbac.ResourceEmployee, policies.ActionView, vp, empOfVP, uVP2, false},
		{"rp does not view employees", rbac.ResourceEmployee, policies.ActionView, rp, empOfVP, uRP, false},
		{"vp creates own employee", rbac.ResourceEmployee, policies.ActionCreate, vp, empOfVP, uVP1, true},
		{
			"vp cannot create employee for another",
			rbac.ResourceEmployee,
			policies.ActionCreate,
			vp,
			rbac.Owners{Owner: uVP2},
			uVP1,
			false,
		},
		{
			"rp cannot create employee",
			rbac.ResourceEmployee,
			policies.ActionCreate,
			rp,
			rbac.Owners{Owner: uRP},
			uRP,
			false,
		},
		{"vp updates own employee", rbac.ResourceEmployee, policies.ActionUpdate, vp, empOfVP, uVP1, true},
		{"vp does not update foreign employee", rbac.ResourceEmployee, policies.ActionUpdate, vp, empOfVP, uVP2, false},
		{"vp deletes own employee", rbac.ResourceEmployee, policies.ActionDelete, vp, empOfVP, uVP1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := policies.Authorize(tc.role, tc.res, tc.act, tc.owners, tc.userID); got != tc.allowed {
				t.Errorf("policies.Authorize(%s, %v, %v) = %v, want %v", tc.role, tc.res, tc.act, got, tc.allowed)
			}
		})
	}
}
