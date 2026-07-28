package domain

import (
	assignmentDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
	milestoneDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
	processDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	projectDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	resourceDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/domain"
	taskDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
)

type Project projectDomain.Project
type Process processDomain.Process
type Task taskDomain.Task
type Resource resourceDomain.Resource
type Assignment assignmentDomain.Assignment
type Milestone milestoneDomain.Milestone
