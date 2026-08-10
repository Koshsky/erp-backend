package delivery

import (
	"log/slog"
	"strconv"

	"github.com/Koshsky/erp-backend/internal/timesheet/employee/service"

	"github.com/gin-gonic/gin"

	userctx "github.com/Koshsky/erp-backend/internal/userctx"
	"github.com/Koshsky/erp-backend/pkg/date"
	"github.com/Koshsky/erp-backend/pkg/errors"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/response"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/dto"
)

type EmployeeHandler struct {
	logger  *slog.Logger
	service EmployeeService
	mw      *rbac.Middleware
}

// NewEmployeeHandler builds the EmployeeHandler handler.
func NewEmployeeHandler(logger *slog.Logger, svc *service.EmployeeService, mw *rbac.Middleware) *EmployeeHandler {
	return &EmployeeHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
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
//	@Success		200	{object}	response.SuccessResponse{data=[]dto.EmployeeResponse,error=nil}
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/resources/{id}/employees [get]
func (h *EmployeeHandler) ListEmployeesByResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid resource id")
		return
	}

	employees, err := h.service.ListEmployeesByResource(c.Request.Context(), id)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.OK(c, employees)
}

// ListEmployees handles the request to list all employees.
//
//	@Tags			TimesheetEmployees
//	@Summary		List all employees
//	@Description	List all employees
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=[]dto.EmployeeResponse,error=nil}
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/employees [get]
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	employees, err := h.service.ListEmployees(c.Request.Context())
	if err != nil {
		response.Error(c, h.logger, err)
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
//	@Success		200	{object}	response.SuccessResponse{data=dto.EmployeeResponse,error=nil}
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/employees/{id} [get]
func (h *EmployeeHandler) FindEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid employee id")
		return
	}

	employee, err := h.service.FindEmployee(c.Request.Context(), id)
	if err != nil {
		response.Error(c, h.logger, err)
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
//	@Success		201			{object}	response.SuccessResponse{data=dto.EmployeeResponse,error=nil}
//	@Failure		400			{object}	response.ErrorResponse{data=nil}
//	@Failure		500			{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/resources/{id}/employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid resource id")
		return
	}

	var employee dto.CreateEmployeeRequest
	if err = c.ShouldBindJSON(&employee); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, err.Error())
		return
	}

	user, err := userctx.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}

	created, err := h.service.CreateEmployee(c.Request.Context(), id, employee, user.ID, user.Role)
	if err != nil {
		response.Error(c, h.logger, err)
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
//	@Success		200			{object}	response.SuccessResponse{data=dto.EmployeeResponse,error=nil}
//	@Failure		400			{object}	response.ErrorResponse{data=nil}
//	@Failure		500			{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid employee id")
		return
	}

	body := dto.UpdateEmployeeRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, err.Error())
		return
	}

	updated, err := h.service.UpdateEmployee(c.Request.Context(), id, body)
	if err != nil {
		response.Error(c, h.logger, err)
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
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid employee id")
		return
	}

	if err = h.service.DeleteEmployee(c.Request.Context(), id); err != nil {
		response.Error(c, h.logger, err)
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
//	@Success		200			{object}	response.SuccessResponse{data=[]dto.EmployeeStateResponse,error=nil}
//	@Failure		400			{object}	response.ErrorResponse{data=nil}
//	@Failure		500			{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/employees/{id}/days [get]
func (h *EmployeeHandler) ListDays(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid employee id")
		return
	}

	start, err := date.Parse(c.Query("start_date"))
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid start_date")
		return
	}
	end, err := date.Parse(c.Query("end_date"))
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid end_date")
		return
	}

	states, err := h.service.ListStates(c.Request.Context(), id, start, end)
	if err != nil {
		response.Error(c, h.logger, err)
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
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/employees/{id}/days [put]
func (h *EmployeeHandler) SetDays(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid employee id")
		return
	}

	body := dto.SetDaysRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, err.Error())
		return
	}

	if err = h.service.SetDays(c.Request.Context(), id, body); err != nil {
		response.Error(c, h.logger, err)
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
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/timesheet/employees/{id}/days [delete]
func (h *EmployeeHandler) DeleteDays(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid employee id")
		return
	}

	start, err := date.Parse(c.Query("start_date"))
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid start_date")
		return
	}
	end, err := date.Parse(c.Query("end_date"))
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid end_date")
		return
	}

	var stateID *int64
	if raw := c.Query("state_id"); raw != "" {
		parsedID, parseErr := strconv.ParseInt(raw, 10, 64)
		if parseErr != nil {
			response.BadRequest(c, errors.CodeBadRequest, "invalid state_id")
			return
		}
		stateID = &parsedID
	}

	if err = h.service.DeleteDays(c.Request.Context(), id, start, end, stateID); err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.NoContent(c)
}
