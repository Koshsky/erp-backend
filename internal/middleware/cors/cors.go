package cors

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Config struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
}

func DefaultConfig() Config {
	return Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
}

func New(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowed := false
		for _, allowedOrigin := range config.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		if len(config.AllowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowMethods))
		}

		// Отражаем перечень заголовков, запрошенных клиентом в preflight.
		// Это надёжнее жёсткого AllowHeaders: браузер сам не формирует произвольные
		// заголовки, а сервер разрешает ровно то, что клиент реально запросил.
		// Такой подход (эхо Access-Control-Request-Headers) используют rs/cors и GoCORS.
		if requestedHeaders := c.Request.Header.Get("Access-Control-Request-Headers"); requestedHeaders != "" {
			c.Header("Access-Control-Allow-Headers", requestedHeaders)
		} else if len(config.AllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowHeaders))
		}

		if len(config.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", joinStrings(config.ExposeHeaders))
		}

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func joinStrings(strs []string) string {
	var result strings.Builder
	for i, s := range strs {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(s)
	}
	return result.String()
}

func Default() gin.HandlerFunc {
	return New(DefaultConfig())
}

func Development() gin.HandlerFunc {
	config := Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	return New(config)
}
