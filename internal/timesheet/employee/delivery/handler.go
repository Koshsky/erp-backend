package delivery

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/common/helpers"
	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/dto"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/service"
)

type EmployeeHandler struct {
	logger  *slog.Logger
	service EmployeeService
}

func NewEmployeeHandler(logger *slog.Logger, service EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		logger:  logger,
		service: service,
	}
}

// handleError маппит доменные ошибки сервиса сотрудников в HTTP-статусы.
func (h *EmployeeHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrForbidden):
		response.Forbidden(c, err.Error())
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, err.Error())
	default:
		response.InternalError(c, h.logger, err.Error(), err)
	}
}

// ListEmployeesByResource handles the request to list employees of a resource.
//
//	@Tags			TimesheetEmployees
//	@Summary		List employees
//	@Description	List concrete employees of a resource category
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"Resource ID"
//	@Success		200	{object}	response.Response{data=[]dto.EmployeeResponse}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/timesheet/resources/{id}/employees [get]
func (h *EmployeeHandler) ListEmployeesByResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid resource id")
		return
	}

	employees, err := h.service.ListEmployeesByResource(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, employees)
}

// ListEmployees handles the request to list all employees.
//
//	@Tags			TimesheetEmployees
//	@Summary		List all employees
//	@Description	List employees, optionally filtered by manager (user) id
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			manager_id	query		int	false	"Manager (user) ID filter"
//	@Success		200			{object}	response.Response{data=[]dto.EmployeeResponse}
//	@Failure		400			{object}	response.Response{data=nil}
//	@Failure		500			{object}	response.Response{data=nil}
//	@Router			/timesheet/employees [get]
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	var managerID *int64
	if raw := c.Query("manager_id"); raw != "" {
		parsedID, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			response.BadRequest(c, "invalid manager_id")
			return
		}
		managerID = &parsedID
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	employees, err := h.service.ListEmployees(c.Request.Context(), managerID, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, employees)
}

// FindEmployee handles the request to get an employee by id.
//
//	@Tags			TimesheetEmployees
//	@Summary		Get employee
//	@Description	Get an employee by id
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"Employee ID"
//	@Success		200	{object}	response.Response{data=dto.EmployeeResponse}
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/timesheet/employees/{id} [get]
func (h *EmployeeHandler) FindEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid employee id")
		return
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	employee, err := h.service.FindEmployee(c.Request.Context(), id, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, employee)
}

// CreateEmployee handles the request to create an employee for a resource.
//
//	@Tags			TimesheetEmployees
//	@Summary		Create employee
//	@Description	Create a concrete employee for a resource category
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int							true	"Resource ID"
//	@Param			employee	body		dto.CreateEmployeeRequest	true	"Employee"
//	@Success		201			{object}	response.Response{data=dto.EmployeeResponse}
//	@Failure		400			{object}	response.Response{data=nil}
//	@Failure		500			{object}	response.Response{data=nil}
//	@Router			/timesheet/resources/{id}/employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid resource id")
		return
	}

	var employee dto.CreateEmployeeRequest
	if err = c.ShouldBindJSON(&employee); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	created, err := h.service.CreateEmployee(c.Request.Context(), id, employee, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Created(c, created)
}

// UpdateEmployee handles the request to update an employee by id.
//
//	@Tags			TimesheetEmployees
//	@Summary		Update employee
//	@Description	Update an employee by id
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int							true	"Employee ID"
//	@Param			employee	body		dto.UpdateEmployeeRequest	true	"Employee"
//	@Success		200			{object}	response.Response{data=dto.EmployeeResponse}
//	@Failure		400			{object}	response.Response{data=nil}
//	@Failure		500			{object}	response.Response{data=nil}
//	@Router			/timesheet/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid employee id")
		return
	}

	body := dto.UpdateEmployeeRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	updated, err := h.service.UpdateEmployee(c.Request.Context(), id, body, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, updated)
}

// DeleteEmployee handles the request to delete an employee by id.
//
//	@Tags			TimesheetEmployees
//	@Summary		Delete employee
//	@Description	Delete an employee by id (soft delete)
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path	int	true	"Employee ID"
//	@Success		204
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/timesheet/employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid employee id")
		return
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	if err = h.service.DeleteEmployee(c.Request.Context(), id, user.ID, user.Role); err != nil {
		h.handleError(c, err)
		return
	}
	response.NoContent(c)
}

// ListDays handles the request to list calendar states of an employee in a range.
//
//	@Tags			TimesheetEmployees
//	@Summary		List employee days
//	@Description	List state ranges of an employee overlapping a date range
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id			path		int		true	"Employee ID"
//	@Param			start_date	query		string	true	"Start date (YYYY-MM-DD)"
//	@Param			end_date	query		string	true	"End date (YYYY-MM-DD)"
//	@Success		200			{object}	response.Response{data=[]dto.EmployeeStateResponse}
//	@Failure		400			{object}	response.Response{data=nil}
//	@Failure		500			{object}	response.Response{data=nil}
//	@Router			/timesheet/employees/{id}/days [get]
func (h *EmployeeHandler) ListDays(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid employee id")
		return
	}

	start, err := date.Parse(c.Query("start_date"))
	if err != nil {
		response.BadRequest(c, "invalid start_date")
		return
	}
	end, err := date.Parse(c.Query("end_date"))
	if err != nil {
		response.BadRequest(c, "invalid end_date")
		return
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	states, err := h.service.ListStates(c.Request.Context(), id, start, end, user.ID, user.Role)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, states)
}

// SetDays handles the request to set a state for a range of days.
//
//	@Tags			TimesheetEmployees
//	@Summary		Set employee days
//	@Description	Overwrite a state on a date range of an employee's calendar (splits overlapping ranges)
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int					true	"Employee ID"
//	@Param			body	body	dto.SetDaysRequest	true	"Days"
//	@Success		204
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/timesheet/employees/{id}/days [put]
func (h *EmployeeHandler) SetDays(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid employee id")
		return
	}

	body := dto.SetDaysRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	if err = h.service.SetDays(c.Request.Context(), id, body, user.ID, user.Role); err != nil {
		h.handleError(c, err)
		return
	}
	response.NoContent(c)
}

// DeleteDays handles the request to delete calendar days of an employee.
//
//	@Tags			TimesheetEmployees
//	@Summary		Delete employee days
//	@Description	Clear state ranges of an employee overlapping a date range (splits overlaps, optional state filter)
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id			path	int		true	"Employee ID"
//	@Param			start_date	query	string	true	"Start date (YYYY-MM-DD)"
//	@Param			end_date	query	string	true	"End date (YYYY-MM-DD)"
//	@Param			state_id	query	int		false	"Optional state filter"
//	@Success		204
//	@Failure		400	{object}	response.Response{data=nil}
//	@Failure		500	{object}	response.Response{data=nil}
//	@Router			/timesheet/employees/{id}/days [delete]
func (h *EmployeeHandler) DeleteDays(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid employee id")
		return
	}

	start, err := date.Parse(c.Query("start_date"))
	if err != nil {
		response.BadRequest(c, "invalid start_date")
		return
	}
	end, err := date.Parse(c.Query("end_date"))
	if err != nil {
		response.BadRequest(c, "invalid end_date")
		return
	}

	var stateID *int64
	if raw := c.Query("state_id"); raw != "" {
		parsedID, parseErr := strconv.ParseInt(raw, 10, 64)
		if parseErr != nil {
			response.BadRequest(c, "invalid state_id")
			return
		}
		stateID = &parsedID
	}

	user, err := helpers.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	if err = h.service.DeleteDays(c.Request.Context(), id, start, end, stateID, user.ID, user.Role); err != nil {
		h.handleError(c, err)
		return
	}
	response.NoContent(c)
}
