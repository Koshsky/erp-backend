package delivery

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Koshsky/erp/api/internal/user/dto"
	"github.com/gin-gonic/gin"
)

// Response structures
type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeBindError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Warn("invalid request payload", "error", err)
	c.JSON(http.StatusBadRequest, response{Error: "invalid request payload"})
}

type UserHandler struct {
	logger  *slog.Logger
	service UserService
}

func NewUserHandler(logger *slog.Logger, service UserService) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: service,
	}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: users})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid user id"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: user})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var body dto.CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}
	created, err := h.service.CreateUser(c.Request.Context(), body)
	if err != nil {
		h.logger.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response{Data: created})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid user id"})
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "invalid user id"})
		return
	}

	body := dto.UpdateUserRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		writeBindError(c, h.logger, err)
		return
	}

	updated, err := h.service.UpdateUser(c.Request.Context(), id, body)
	if err != nil {
		h.logger.Error("failed to update user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Data: updated})
}

// func (h *UserHandler) UpdateUserRole(c *gin.Context) {
// 	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, response{Error: "invalid user id"})
// 		return
// 	}

// 	var body struct {
// 		Role domain.UserRole `json:"role"`
// 	}
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, response{Error: err.Error()})
// 		return
// 	}

// 	user, err := h.service.UpdateUserRole(c.Request.Context(), id, body.Role)
// 	if err != nil {
// 		if isNotFoundError(err) {
// 			c.JSON(http.StatusNotFound, response{Error: "user not found"})
// 			return
// 		}
// 		h.logger.Error("failed to update user role", "id", id, "error", err)
// 		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, response{Data: user})
// }

// func (h *UserHandler) UpdateUserPassword(c *gin.Context) {
// 	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, response{Error: "invalid user id"})
// 		return
// 	}

// 	var body struct {
// 		Password string `json:"password"`
// 	}
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, response{Error: err.Error()})
// 		return
// 	}
// 	if body.Password == "" {
// 		c.JSON(http.StatusBadRequest, response{Error: "password is required"})
// 		return
// 	}

// 	user, err := h.service.UpdateUserPassword(c.Request.Context(), id, body.Password)
// 	if err != nil {
// 		if isNotFoundError(err) {
// 			c.JSON(http.StatusNotFound, response{Error: "user not found"})
// 			return
// 		}
// 		h.logger.Error("failed to update user password", "id", id, "error", err)
// 		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, response{Data: user})
// }
