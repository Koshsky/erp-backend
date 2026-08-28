package delivery

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	authservice "github.com/Koshsky/erp-backend/internal/auth/service"
	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
	"github.com/Koshsky/erp-backend/internal/response"
)

// refreshCookieName — the HttpOnly cookie holding the opaque refresh token (AD-05).
const refreshCookieName = "mvs_refresh"

// refreshCookiePath restricts sending the cookie to auth endpoints only.
const refreshCookiePath = "/api/v1/auth"

type AuthHandler struct {
	service AuthService
	logger  *slog.Logger
	cfg     config.JWTConfig
}

// NewAuthHandler builds the auth handler.
func NewAuthHandler(logger *slog.Logger, svc *authservice.AuthService, cfg config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: svc,
		cfg:     cfg,
	}
}

// isHTTPS reports whether the request came over real https: direct TLS or
// X-Forwarded-Proto/Scheme from the reverse proxy (nginx on /api/ sets
// X-Forwarded-Proto $scheme, overwriting the client header).
func isHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	proto := c.GetHeader("X-Forwarded-Proto")
	if proto == "" {
		proto = c.GetHeader("X-Forwarded-Scheme")
	}
	return strings.EqualFold(proto, "https")
}

// setRefreshCookie sets the HttpOnly/SameSite=Strict refresh cookie. The Secure
// flag is applied only behind real https (cfg.RefreshCookieSecure — the "strict
// https requirement"): over http the browser will not persist a Secure cookie,
// /auth/refresh becomes unavailable and the web version keeps failing with
// "Session expired".
//
//nolint:gosec // Secure is driven by refresh_cookie_secure + the actual protocol
func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   h.cfg.RefreshCookieSecure && isHTTPS(c),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(time.Duration(h.cfg.RefreshExpiry).Seconds()),
	})
}

// clearRefreshCookie removes the refresh cookie (logout). Secure must match
// how the cookie was set (the browser identifies a cookie by its attributes).
//
//nolint:gosec // Secure likewise from refresh_cookie_secure; the other flags are static.
func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   h.cfg.RefreshCookieSecure && isHTTPS(c),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

// Login handles the login request.
//
//	@Tags			Auth
//	@Summary		Login
//	@Description	Authenticate user; the refresh token goes into an HttpOnly cookie
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	response.SuccessResponse{data=dto.AuthResponse,error=nil}
//	@Failure		400		{object}	response.ErrorResponse{data=nil}
//	@Failure		401		{object}	response.ErrorResponse{data=nil}
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, errors.CodeBadRequest, "invalid request")
		return
	}

	res, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, errors.CodeInvalidCredentials, "invalid credentials")
		return
	}

	h.setRefreshCookie(c, res.RefreshToken)
	response.OK(c, res.Auth)
}

// RefreshToken handles the request to refresh access token using the HttpOnly cookie.
//
//	@Tags			Auth
//	@Summary		Refresh Token
//	@Description	Rotate the refresh session from the HttpOnly cookie; returns a new access token
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=dto.AuthResponse,error=nil}
//	@Failure		401	{object}	response.ErrorResponse{data=nil}
//	@Router			/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	token, err := c.Cookie(refreshCookieName)
	if err != nil || token == "" {
		response.Unauthorized(c, errors.CodeInvalidToken, "invalid refresh token")
		return
	}

	res, err := h.service.RefreshToken(c.Request.Context(), token)
	if err != nil {
		response.Unauthorized(c, errors.CodeInvalidToken, "invalid refresh token")
		return
	}

	h.setRefreshCookie(c, res.RefreshToken)
	response.OK(c, res.Auth)
}

// Logout revokes the refresh session and clears the cookie.
//
//	@Tags			Auth
//	@Summary		Logout
//	@Description	Revoke the refresh session and clear the cookie (idempotent)
//	@Produce		json
//	@Success		200	{object}	response.SuccessResponse{data=map[string]string,error=nil}
//	@Router			/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	if token, err := c.Cookie(refreshCookieName); err == nil && token != "" {
		_ = h.service.Logout(c.Request.Context(), token)
	}
	h.clearRefreshCookie(c)
	response.OK(c, map[string]string{"message": "logged out"})
}
