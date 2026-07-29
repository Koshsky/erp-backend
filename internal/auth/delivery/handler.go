package delivery

import (
	"log/slog"
	"net/http"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
	"github.com/Koshsky/erp-backend/internal/common/ctx"
	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service AuthService
	logger  *slog.Logger
}

func NewAuthHandler(logger *slog.Logger, service AuthService) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: service,
	}
}

// @Tags Auth
// @Summary Login
// @Description Authenticate user and return JWT token
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, "invalid credentials")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Tags Auth
// @Summary Register
// @Description Create user and return JWT token
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register credentials"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req.Name, req.Username, req.Password)
	if err != nil {
		response.InternalError(c, h.logger, "failed to register user", err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Tags Auth
// @Summary Change Password
// @Description Change password (requires old password)
// @Security ApiKeyAuth
// @Accept json
// @Param request body dto.ChangePasswordRequest true "Old and new password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := ctx.GetUserID(c.Request.Context())

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.BadRequest(c, "invalid password")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// @Tags Auth
// @Summary Refresh Token
// @Description Refresh access token using refresh token, returns new pair
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "Refresh token"
// @Success 200 {object} dto.RefreshResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "refresh_token required")
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid refresh token")
		return
	}

	c.JSON(http.StatusOK, resp)
}
