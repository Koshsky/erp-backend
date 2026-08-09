// Package timesheet wires the timesheet module's providers and routes.
package timesheet

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	calendarDelivery "github.com/Koshsky/erp-backend/internal/timesheet/calendar/delivery"
	calendarRepo "github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository"
	calendarService "github.com/Koshsky/erp-backend/internal/timesheet/calendar/service"
	employeeDelivery "github.com/Koshsky/erp-backend/internal/timesheet/employee/delivery"
	employeeRepo "github.com/Koshsky/erp-backend/internal/timesheet/employee/repository"
	employeeService "github.com/Koshsky/erp-backend/internal/timesheet/employee/service"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/timesheet/resource/delivery"
	resourceRepo "github.com/Koshsky/erp-backend/internal/timesheet/resource/repository"
	resourceService "github.com/Koshsky/erp-backend/internal/timesheet/resource/service"
	stateDelivery "github.com/Koshsky/erp-backend/internal/timesheet/state/delivery"
	stateRepo "github.com/Koshsky/erp-backend/internal/timesheet/state/repository"
	stateService "github.com/Koshsky/erp-backend/internal/timesheet/state/service"
)

// ProviderSet aggregates the timesheet module's dependencies.
var ProviderSet = wire.NewSet(
	resourceRepo.NewResourceRepository,
	resourceService.NewResourceService,
	resourceDelivery.NewResourceHandler,

	stateRepo.NewStateRepository,
	stateService.NewStateService,
	stateDelivery.NewStateHandler,

	employeeRepo.NewEmployeeRepository,
	employeeService.NewEmployeeService,
	employeeDelivery.NewEmployeeHandler,

	calendarRepo.NewCalendarRepository,
	calendarService.NewCalendarService,
	calendarDelivery.NewCalendarHandler,

	ProvideModule,
)

// Module registers the timesheet module's routes (all protected).
type Module struct {
	resource *resourceDelivery.ResourceHandler
	state    *stateDelivery.StateHandler
	employee *employeeDelivery.EmployeeHandler
	calendar *calendarDelivery.CalendarHandler
}

// ProvideModule builds the timesheet module.
func ProvideModule(
	resource *resourceDelivery.ResourceHandler,
	state *stateDelivery.StateHandler,
	employee *employeeDelivery.EmployeeHandler,
	calendar *calendarDelivery.CalendarHandler,
) Module {
	return Module{
		resource: resource,
		state:    state,
		employee: employee,
		calendar: calendar,
	}
}

// RegisterPublicRoutes is a no-op: the module has no public routes.
func (m Module) RegisterPublicRoutes(r *gin.RouterGroup) {
}

// RegisterProtectedRoutes registers the module's routes behind authentication.
func (m Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
	m.resource.RegisterRoutes(r)
	m.state.RegisterRoutes(r)
	m.employee.RegisterRoutes(r)
	m.calendar.RegisterRoutes(r)
}
