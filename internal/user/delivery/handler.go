package delivery

import (
	"log/slog"
	"strconv"

	userservice "github.com/Koshsky/erp-backend/internal/user/service"
	"github.com/Koshsky/erp-backend/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/response"
	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/dto"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
)

// TODO: split into base CRUD and profile management (editing)

type UserHandler struct {
	logger  *slog.Logger
	service UserService
}

// NewUserHandler builds the user handler.
func NewUserHandler(logger *slog.Logger, svc *userservice.UserService) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: svc,
	}
}

// ListUsers handles the request to list all users.
//
//	@Summary		List all users
//	@Description	Returns all users in the system
//	@Tags			Users
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=[]dto.UserResponse,error=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/user [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, users)
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
//	@Router			/user/{id} [get]
func (h *UserHandler) FindUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
		return
	}

	user, err := h.service.FindUser(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, user)
}

// DeleteUser handle deleting a user.
//
//	@Tags			Users
//	@Summary		Delete a user
//	@Description	Delete a user by ID
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Success		204
//	@Failure		400	{object}	response.ErrorResponse{data=nil}
//	@Failure		500	{object}	response.ErrorResponse{data=nil}
//	@Router			/user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
		return
	}

	user, err := userctx.GetUser(c)
	if err != nil {
		response.Unauthorized(c, errors.CodeUnauthorized, "authentication required")
		return
	}
	if user.Role != domain.Admin && user.ID != id {
		response.BadRequest(c, errors.CodeBadRequest, "you can only delete your own account")
		return
	}

	if err = h.service.DeleteUser(c.Request.Context(), id); err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.NoContent(c)
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
//	@Router			/user/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid user id")
		return
	}

	user, err := userctx.GetUser(c)
	if err != nil {
		response.Unauthorized(c, errors.CodeUnauthorized, "authentication required")
		return
	}
	if user.Role != domain.Admin && user.ID != id {
		response.BadRequest(c, errors.CodeBadRequest, "you can only update your own account")
		return
	}

	body := dto.UpdateUserRequest{}
	if err = c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, err.Error())
		return
	}

	updated, err := h.service.UpdateUser(c.Request.Context(), id, body)
	if err != nil {
		response.InternalError(c, h.logger, err.Error(), err)
		return
	}
	response.OK(c, updated)
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
