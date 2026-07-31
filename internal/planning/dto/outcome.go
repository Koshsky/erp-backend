package dto

type ProjectPlanning struct {
	Projects []Project `json:"projects"`
}

type ProcessPlanning struct {
	Projects []DetailedProject `json:"projects"`
}

type TaskPlanning struct {
	Processes []DetailedProcess `json:"processes"`
}

type DetailedProject struct {
	Project
	Processes []Process `json:"processes"`
}

type DetailedProcess struct {
	Process
	Tasks      []DetailedTask `json:"tasks"`
	Milestones []Milestone    `json:"milestones"`
}

type DetailedTask struct {
	Task
	Resources []Resource `json:"resources"`
}
