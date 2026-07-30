package config

import (
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	DatabaseURL string
	LogLevel    string
	LogFormat   string
	Swagger     SwaggerConfig
	JWT         JWTConfig
	HTTPServer  HTTPServerConfig
}

type SwaggerConfig struct {
	Enabled bool
}

type HTTPServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type JWTConfig struct {
	SecretKey     string
	RefreshKey    string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

func Load() (*Config, error) {
	dbURL := loadDatabaseURL()

	logLevel := getEnv("LOG_LEVEL", "info")
	logFormat := getEnv("LOG_FORMAT", "json")

	jwtCfg := JWTConfig{
		SecretKey:     getEnv("JWT_SECRET_KEY", "super-secret-key-change-in-production"),
		RefreshKey:    getEnv("JWT_REFRESH_KEY", "super-refresh-key-change-in-production"),
		AccessExpiry:  getEnvAsDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
		RefreshExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 168*time.Hour),
		Issuer:        getEnv("JWT_ISSUER", "mvs-erp"),
	}

	httpServerCfg := HTTPServerConfig{
		Port:         getEnvAsInt("HTTP_PORT", 8080),
		ReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:  getEnvAsDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
	}

	swaggerCfg := SwaggerConfig{
		Enabled: getEnvAsBool("SWAGGER_ENABLE", false),
	}

	return &Config{
		DatabaseURL: dbURL,
		LogLevel:    logLevel,
		LogFormat:   logFormat,
		Swagger:     swaggerCfg,
		JWT:         jwtCfg,
		HTTPServer:  httpServerCfg,
	}, nil
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
