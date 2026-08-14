package delivery

import (
	"log/slog"
	"strconv"

	userservice "github.com/Koshsky/erp-backend/internal/user/service"
	"github.com/Koshsky/erp-backend/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/response"
	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/dto"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type UserHandler struct {
	logger  *slog.Logger
	service UserService
	mw      *rbac.Middleware
}

// NewUserHandler builds the user handler.
func NewUserHandler(logger *slog.Logger, svc *userservice.UserService, mw *rbac.Middleware) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}

// ListAllUsers handles the request to list all users (unscoped, for owner pickers).
//
//	@Summary		List all users
//	@Description	Returns all users in the system
//	@Tags			Users
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=[]dto.UserResponse,error=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/user [get]
func (h *UserHandler) ListAllUsers(c *gin.Context) {
	users, err := h.service.ListAllUsers(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, users)
}

// ListUsers handles the request to list users with role/manager filters.
//
//	@Summary		List users
//	@Description	Returns a paged list of users; admin sees all, vp sees own subordinates + self.
//	@Tags			Users
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			limit			query		int		false	"Page size (default 50, max 500)"
//	@Param			role			query		string	false	"Filter by role (e.g. worker)"
//	@Param			manager_id		query		int		false	"Filter by manager (admin)"
//	@Param			include_hash	query		bool	false	"Включить password_hash (только admin)"
//	@Param			offset			query		int		false	"Page offset"
//	@Success		200				{object}	response.SuccessResponse{data=response.Page{items=[]dto.AdminUserResponse},error=nil}
//	@Failure		400				{object}	response.ErrorResponse{data=nil}
//	@Failure		403				{object}	response.ErrorResponse{data=nil}
//	@Failure		500				{object}	response.ErrorResponse{data=nil}
//	@Router			/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, offset, perr := response.ParsePagination(c)
	if perr != nil {
		response.Error(c, h.logger, perr)
		return
	}
	user, err := userctx.GetUser(c)
	if err != nil {
		response.Unauthorized(c, errors.CodeUnauthorized, "authentication required")
		return
	}
	includeHash := false
	if raw := c.Query("include_hash"); raw != "" {
		includeHash = raw == "true" || raw == "1"
	}
	if includeHash && user.Role != domain.Admin {
		response.Error(c, h.logger, errors.ErrForbidden)
		return
	}
	items, total, err := h.service.ListUsers(
		c.Request.Context(),
		user.ID,
		user.Role,
		c.Query("role"),
		response.QueryID(c, "manager_id"),
		limit,
		offset,
	)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	out := make([]dto.AdminUserResponse, 0, len(items))
	for _, u := range items {
		hash := ""
		if includeHash {
			hash = u.PasswordHash
		}
		out = append(out, dto.AdminUserResponse{
			ID:              u.ID,
			Name:            u.Name,
			Username:        u.Username,
			Role:            u.Role,
			ManagerID:       u.ManagerID,
			Position:        u.Position,
			HireDate:        u.HireDate,
			TerminationDate: u.TerminationDate,
			PasswordHash:    hash,
		})
	}
	response.OK(c, response.Page{Items: out, Total: total, Limit: limit, Offset: offset})
}

// FindUser handles the request to get a specific user.
//
//	@Summary		User information
//	@Description	Returns information about a specific user
//	@Tags			Users
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.SuccessResponse{data=dto.UserResponse,error=nil}
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/users/{id} [get]
func (h *UserHandler) FindUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
		return
	}

	user, err := h.service.FindUser(c.Request.Context(), id)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.OK(c, user)
}

// CreateUser handles the request to create a user (worker with generated creds).
//
//	@Tags			Users
//	@Summary		Create user
//	@Description	Create a user (auto-generated username/password; password returned once)
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.CreateUserRequest	true	"User"
//	@Success		201		{object}	response.SuccessResponse{data=dto.CreateUserResult,error=nil}
//	@Failure		400		{object}	response.ErrorResponse{data=nil}
//	@Failure		500		{object}	response.ErrorResponse{data=nil}
//	@Router			/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, err.Error())
		return
	}

	user, err := userctx.GetUser(c)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	// vp создаёт рабочего в свою команду: чужой manager_id игнорируется.
	if user.Role == domain.ProcessOwner {
		req.ManagerID = &user.ID
	}

	created, err := h.service.CreateUserWithCreds(c.Request.Context(), req)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.Created(c, created)
}

// ResetPassword handles the request to reset a user's password (generated, returned once).
//
//	@Tags			Users
//	@Summary		Reset user password
//	@Description	Generate a new random password for a user and return it once
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	response.SuccessResponse{data=dto.ResetPasswordResponse,error=nil}
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/users/{id}/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
		return
	}

	res, err := h.service.ResetPassword(c.Request.Context(), id)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.OK(c, res)
}

// UpdateUser handles updating a user.
//
//	@Tags			Users
//	@Summary		Update a user
//	@Description	Update a user by ID
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"User ID"
//	@Param			body	body		dto.UpdateUserRequest	true	"User data"
//	@Success		200		{object}	response.SuccessResponse{data=dto.UserResponse,error=nil}
//	@Failure		400		{object}	response.ErrorResponse{data=nil}
//	@Failure		500		{object}	response.ErrorResponse{data=nil}
//	@Router			/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
		return
	}

	body := dto.UpdateUserRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, err.Error())
		return
	}

	user, err := userctx.GetUser(c)
	if err != nil {
		response.Unauthorized(c, errors.CodeUnauthorized, "authentication required")
		return
	}

	updated, err := h.service.UpdateUser(c.Request.Context(), id, body, user.Role, user.ID)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.OK(c, updated)
}

// UpdateManager handles the request to set/clear a user's manager.
//
//	@Tags			Users
//	@Summary		Update user manager
//	@Description	Explicitly set (or clear, manager_id=null) the manager of a user
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"User ID"
//	@Param			body	body		dto.UpdateManagerRequest	true	"Manager"
//	@Success		200		{object}	response.SuccessResponse{data=dto.UserResponse,error=nil}
//	@Failure		400		{object}	response.ErrorResponse{data=nil}
//	@Failure		403		{object}	response.ErrorResponse{data=nil}
//	@Failure		500		{object}	response.ErrorResponse{data=nil}
//	@Router			/users/{id}/manager [put]
func (h *UserHandler) UpdateManager(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
		return
	}

	var body dto.UpdateManagerRequest
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, err.Error())
		return
	}

	user, err := userctx.GetUser(c)
	if err != nil {
		response.Unauthorized(c, errors.CodeUnauthorized, "authentication required")
		return
	}

	updated, err := h.service.UpdateManager(c.Request.Context(), id, body.ManagerID, user.Role)
	if err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.OK(c, updated)
}

// DeleteUser handle deleting a user.
//
//	@Tags			Users
//	@Summary		Delete a user
//	@Description	Delete a user by ID (soft delete)
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Success		204
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
		return
	}

	if err = h.service.DeleteUser(c.Request.Context(), id); err != nil {
		response.Error(c, h.logger, err)
		return
	}
	response.NoContent(c)
}

// ChangePassword handles the change password request.
//
//	@Tags			Users
//	@Summary		Change Password
//	@Description	Change password (requires old password)
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Param			request	body		dto.ChangePasswordRequest	true	"Old and new password"
//	@Success		200		{object}	response.SuccessResponse{data=dto.ChangePasswordResponse,error=nil}
//	@Failure		400		{object}	response.ErrorResponse{data=nil}
//	@Router			/user/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := userctx.GetUserID(c)

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid request")
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.BadRequest(c, errors.CodeInvalidCredentials, "invalid password")
		return
	}

	response.OK(c, dto.ChangePasswordResponse{Message: "password changed"})
}

// ListDays handles the request to list calendar states of a user in a range.
//
//	@Tags			Users
//	@Summary		List worker days
//	@Description	List state ranges of a worker overlapping a date range
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id			path		int		true	"User ID"
//	@Param			start_date	query		string	true	"Start date (YYYY-MM-DD)"
//	@Param			end_date	query		string	true	"End date (YYYY-MM-DD)"
//	@Success		200			{object}	response.SuccessResponse{data=[]dto.UserStateResponse,error=nil}
//	@Failure		400			{object}	response.ErrorResponse{data=nil}
//	@Failure		500			{object}	response.ErrorResponse{data=nil}
//	@Router			/users/{id}/days [get]
func (h *UserHandler) ListDays(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
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
//	@Tags			Users
//	@Summary		Set worker days
//	@Description	Overwrite a state on a date range of a worker's calendar (splits overlapping ranges)
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int					true	"User ID"
//	@Param			body	body	dto.SetDaysRequest	true	"Days"
//	@Success		204
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/users/{id}/days [put]
func (h *UserHandler) SetDays(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
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

// DeleteDays handles the request to delete calendar days of a user.
//
//	@Tags			Users
//	@Summary		Delete worker days
//	@Description	Clear state ranges of a worker overlapping a date range (splits overlaps, optional state filter)
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id			path	int		true	"User ID"
//	@Param			start_date	query	string	true	"Start date (YYYY-MM-DD)"
//	@Param			end_date	query	string	true	"End date (YYYY-MM-DD)"
//	@Param			state_id	query	int		false	"Optional state filter"
//	@Success		204
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/users/{id}/days [delete]
func (h *UserHandler) DeleteDays(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
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
