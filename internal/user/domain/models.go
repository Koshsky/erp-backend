package domain

const (
	ProjectDirector string = "ДП"
	ProjectManager  string = "РП"
	ProcessOwner    string = "ВП"
)

/*
ProjectDirector (ДП) views, creates and deletes projects;
Edits: code, priority, start_date, end_date,
assigns a РП to the project (owner_id)
+views all processes and tasks

ProjectManager (РП) views their projects, processes
creates and deletes processes,
edits processes: title, start_date, end_date, owner_id
+views their processes and their tasks

ProcessOwner (ВП) views their processes, tasks, resources (shared)
Creates and deletes tasks,
Edits tasks: title, start_date, end_date
Assigns and changes resources for tasks (assignment)

Summary
PROJECTS
view
if role == "ДП" OR project.owner_id == user_id
create/edit/delete
if role == 'ДП'

PROCESSES
view
if role == 'ДП' OR project.owner_id == user_id OR process.owner_id == user_id
create/edit/delete
if project.owner_id == user_id

TASKS
view
if role == 'ДП' OR project.owner_id == user_id OR process.owner_id == user_id
create/edit/delete
if process.owner_id == user_id
ASSIGNMENTS
view
if role == 'ДП' OR project.owner_id == user_id OR process.owner_id == user_id
create/edit/delete
if process.owner_id == user_id
RESOURCES
view
if role == 'ДП' OR project.owner_id == user_id OR process.owner_id == user_id
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
