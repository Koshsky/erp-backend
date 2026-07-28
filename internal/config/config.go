package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Koshsky/erp-backend/internal/security/jwt"
)

type Config struct {
	DatabaseURL   string
	LogLevel      string
	LogFormat     string
	SwaggerEnable bool
	JWT           jwt.JWTConfig
}

func Load() (*Config, error) {
	dbURL := loadDatabaseURL()

	logLevel := getEnv("LOG_LEVEL", "info")
	logFormat := getEnv("LOG_FORMAT", "json")

	swaggerEnable := os.Getenv("SWAGGER_ENABLE") == "true"

	jwtCfg := jwt.JWTConfig{
		SecretKey:     getEnv("JWT_SECRET_KEY", "super-secret-key-change-in-production"),
		RefreshKey:    getEnv("JWT_REFRESH_KEY", "super-refresh-key-change-in-production"),
		AccessExpiry:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"), 15*time.Minute),
		RefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"), 168*time.Hour),
		Issuer:        getEnv("JWT_ISSUER", "mvs-erp"),
	}

	return &Config{
		DatabaseURL:   dbURL,
		LogLevel:      logLevel,
		LogFormat:     logFormat,
		SwaggerEnable: swaggerEnable,
		JWT:           jwtCfg,
	}, nil
}

func getEnv(key, fallback string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return fallback
}

func parseDuration(val string, fallback time.Duration) time.Duration {
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}

func loadDatabaseURL() string {
	user := getEnv("DATABASE_USER", "postgres")
	password := getEnv("DATABASE_PASSWORD", "postgres")
	host := getEnv("DATABASE_HOST", "localhost")
	port := getEnv("DATABASE_PORT", "5432")
	name := getEnv("DATABASE_NAME", "mvs-erp")
	sslMode := getEnv("DATABASE_SSL_MODE", "disable")

	built := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   name,
	}

	query := built.Query()
	query.Set("sslmode", sslMode)
	built.RawQuery = query.Encode()

	return built.String()
}
