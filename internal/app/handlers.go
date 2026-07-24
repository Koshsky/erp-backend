package app

import (
	"github.com/gin-gonic/gin"

	assignmentDelivery "github.com/Koshsky/erp/api/internal/assignment/delivery"
	milestoneDelivery "github.com/Koshsky/erp/api/internal/milestone/delivery"
	processDelivery "github.com/Koshsky/erp/api/internal/process/delivery"
	projectDelivery "github.com/Koshsky/erp/api/internal/project/delivery"
	resourceDelivery "github.com/Koshsky/erp/api/internal/resource/delivery"
	schedulingDelivery "github.com/Koshsky/erp/api/internal/scheduling/delivery"
	taskDelivery "github.com/Koshsky/erp/api/internal/task/delivery"
	userDelivery "github.com/Koshsky/erp/api/internal/user/delivery"
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
