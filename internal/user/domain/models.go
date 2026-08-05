package domain

const (
	Admin           string = "admin"
	ProjectDirector string = "dp"
	ProjectManager  string = "rp"
	ProcessOwner    string = "vp"
	Worker          string = "worker"
)

/*
Admin (admin) is the system administrator;
Views, creates and deletes projects,
Edits: code, priority, start_date, end_date,
assigns a РП to the project (owner_id),
Manages all users (create, edit, delete),
+views all processes and tasks

ProjectDirector (dp) manages the project portfolio;
views, creates and deletes projects,
Edits: code, priority, start_date, end_date,
assigns a РП to the project (owner_id)
+views all processes and tasks

ProjectManager (rp) views their projects, processes
creates and deletes processes,
edits processes: title, start_date, end_date, owner_id
+views their processes and their tasks

ProcessOwner (vp) views their processes, tasks, resources (shared)
Creates and deletes tasks,
Edits tasks: title, start_date, end_date
Assigns and changes resources for tasks (assignment)

Worker (worker) is an executor assigned to tasks.
Rights are not implemented yet.

Summary
PROJECTS
view
if role IN ('admin', 'dp') OR project.owner_id == user_id
create/edit/delete
if role IN ('admin', 'dp')

PROCESSES
view
if role IN ('admin', 'dp') OR project.owner_id == user_id OR process.owner_id == user_id
create/edit/delete
if project.owner_id == user_id

TASKS
view
if role IN ('admin', 'dp') OR project.owner_id == user_id OR process.owner_id == user_id
create/edit/delete
if process.owner_id == user_id
ASSIGNMENTS
view
if role IN ('admin', 'dp') OR project.owner_id == user_id OR process.owner_id == user_id
create/edit/delete
if process.owner_id == user_id
RESOURCES
view
if role IN ('admin', 'dp') OR project.owner_id == user_id OR process.owner_id == user_id
create/edit/delete // TODO: who?
if process.owner_id == user_id
*/

type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	Username     string `json:"username"`
	ManagerID    *int64 `json:"manager_id"`
	PasswordHash string `json:"-"`
}
