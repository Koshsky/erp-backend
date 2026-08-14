package domain

import (
	"strings"
	"time"
)

const (
	Admin           string = "admin"
	ProjectDirector string = "dp"
	ProjectManager  string = "rp"
	ProcessOwner    string = "vp"
	Worker          string = "worker"
)

/*
Roles:
Admin (admin) — system administrator: full access to all entities and user management.
ProjectDirector (dp) — project portfolio director: sees all projects, processes,
tasks, milestones and assignments; may change only project priorities.
ProjectManager (rp) — project manager: sees own projects and the processes of own
projects (and their tasks/milestones/assignments — view only); creates own projects
and processes in own projects, edits/deletes them; does not change milestones/tasks/assignments.
ProcessOwner (vp) — process owner: sees own processes and their tasks/milestones/
assignments; creates, edits and deletes tasks, milestones and assignments in own
processes. Does not create or edit processes.
Worker (worker) — no rights yet.

Permission matrix (admin — everything; worker — nothing):
Projects:
  view: dp — all, rp — own
  create: rp (into own ownership)
  edit fields (code/dates): rp — own (does not change the owner)
  edit priority: dp — all
  delete: rp — own
Processes:
  view: dp — all, rp — of own projects, vp — own
  create/edit/delete: rp — in own projects
Tasks / Milestones / Assignments:
  view: dp — all, rp — of own projects, vp — of own processes
  create/edit/delete: vp — in own processes

Timesheet:
Resources:
  view/edit/delete: admin — all, vp — own
  create: admin, vp (into own ownership)
Employees:
  view/edit/delete: admin — all, vp — own subordinates
  create: admin, vp (into own team)
States:
  view: admin and vp (reference for the timesheet)
  create/edit/delete: admin

Implementation: internal/middleware/rbac (single matrix + owner chains) plus
middleware on the routes.
*/

type User struct {
	ID           int64
	LastName     string
	FirstName    string
	MiddleName   *string
	Role         string
	Username     string
	PasswordHash string
	// Руководитель рабочего (user с ролью vp); для остальных ролей — nil.
	ManagerID       *int64
	Position        string
	HireDate        *time.Time
	TerminationDate *time.Time
}

// FullName возвращает полное ФИО «Фамилия Имя Отчество» (без пустых частей).
func (u *User) FullName() string {
	parts := []string{}
	if u.LastName != "" {
		parts = append(parts, u.LastName)
	}
	if u.FirstName != "" {
		parts = append(parts, u.FirstName)
	}
	if u.MiddleName != nil && *u.MiddleName != "" {
		parts = append(parts, *u.MiddleName)
	}
	return strings.Join(parts, " ")
}

// UserState is a worker state range (non-presence only), [StartDate, EndDate].
type UserState struct {
	ID          int64
	UserID      int64
	StateID     int64
	StateCode   string
	StateName   string
	IsAvailable bool
	StartDate   time.Time
	EndDate     time.Time
}
