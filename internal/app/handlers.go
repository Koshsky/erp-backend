package app

import (
	"github.com/gin-gonic/gin"

	assignmentDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	milestoneDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/delivery"
	processDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/process/delivery"
	projectDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/project/delivery"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/delivery"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	schedulingDelivery "github.com/Koshsky/erp-backend/internal/scheduling/delivery"
	userDelivery "github.com/Koshsky/erp-backend/internal/user/delivery"
)

type RouteRegistrar interface {
	RegisterRoutes(router *gin.RouterGroup)
}

var (
	_ RouteRegistrar = (*schedulingDelivery.SchedulingHandler)(nil)
	_ RouteRegistrar = (*userDelivery.UserHandler)(nil)
	_ RouteRegistrar = (*taskDelivery.TaskHandler)(nil)
	_ RouteRegistrar = (*resourceDelivery.ResourceHandler)(nil)
	_ RouteRegistrar = (*projectDelivery.ProjectHandler)(nil)
	_ RouteRegistrar = (*processDelivery.ProcessHandler)(nil)
	_ RouteRegistrar = (*milestoneDelivery.MilestoneHandler)(nil)
	_ RouteRegistrar = (*assignmentDelivery.AssignmentHandler)(nil)
)
