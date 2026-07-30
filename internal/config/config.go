package config

import (
	"fmt"
	"net/url"
	"time"
)

// Default values for configuration
const (
	// JWT defaults
	defaultJWTAccessExpiry  = 15 * time.Minute
	defaultJWTRefreshExpiry = 168 * time.Hour // 7 days

	// HTTP server defaults
	defaultHTTPPort         = 8080
	defaultHTTPReadTimeout  = 15 * time.Second
	defaultHTTPWriteTimeout = 15 * time.Second
	defaultHTTPIdleTimeout  = 60 * time.Second

	// Postgres pool defaults
	defaultPostgresMaxConns          = 25
	defaultPostgresMinConns          = 5
	defaultPostgresMaxConnLifetime   = 30 * time.Minute
	defaultPostgresMaxConnIdleTime   = 5 * time.Minute
	defaultPostgresHealthCheckPeriod = 5 * time.Minute
	defaultPostgresConnectTimeout    = 10 * time.Second
)

// Config — application configuration
type Config struct {
	LogLevel   string
	LogFormat  string
	Swagger    SwaggerConfig
	JWT        JWTConfig
	HTTPServer HTTPServerConfig
	Postgres   PostgresConfig
}

type JWTConfig struct {
	SecretKey     string
	RefreshKey    string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

// HTTPServerConfig — configuration for http server
type HTTPServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// PostgresConfig — configuration for postgres database
type PostgresConfig struct {
	DSN               string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout    time.Duration
}

// SwaggerConfig — configuration for swagger
type SwaggerConfig struct {
	Enabled bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	dbURL := loadDatabaseURL()

	logLevel := getEnv("LOG_LEVEL", "info")
	logFormat := getEnv("LOG_FORMAT", "json")

	jwtCfg := JWTConfig{
		SecretKey:     getEnv("JWT_SECRET_KEY", "super-secret-key-change-in-production"),
		RefreshKey:    getEnv("JWT_REFRESH_KEY", "super-refresh-key-change-in-production"),
		AccessExpiry:  getEnvAsDuration("JWT_ACCESS_EXPIRY", defaultJWTAccessExpiry),
		RefreshExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", defaultJWTRefreshExpiry),
		Issuer:        getEnv("JWT_ISSUER", "mvs-erp"),
	}

	httpServerCfg := HTTPServerConfig{
		Port:         getEnvAsInt("HTTP_PORT", defaultHTTPPort),
		ReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", defaultHTTPReadTimeout),
		WriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", defaultHTTPWriteTimeout),
		IdleTimeout:  getEnvAsDuration("HTTP_IDLE_TIMEOUT", defaultHTTPIdleTimeout),
	}

	swaggerCfg := SwaggerConfig{
		Enabled: getEnvAsBool("SWAGGER_ENABLE", false),
	}

	pgCfg := PostgresConfig{
		DSN:               dbURL,
		MaxConns:          getEnvAsInt32("POSTGRES_MAX_CONNS", defaultPostgresMaxConns),
		MinConns:          getEnvAsInt32("POSTGRES_MIN_CONNS", defaultPostgresMinConns),
		MaxConnLifetime:   getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME", defaultPostgresMaxConnLifetime),
		MaxConnIdleTime:   getEnvAsDuration("POSTGRES_MAX_CONN_IDLE_TIME", defaultPostgresMaxConnIdleTime),
		HealthCheckPeriod: getEnvAsDuration("POSTGRES_HEALTH_CHECK_PERIOD", defaultPostgresHealthCheckPeriod),
		ConnectTimeout:    getEnvAsDuration("POSTGRES_CONNECT_TIMEOUT", defaultPostgresConnectTimeout),
	}

	return &Config{
		LogLevel:   logLevel,
		LogFormat:  logFormat,
		Swagger:    swaggerCfg,
		JWT:        jwtCfg,
		HTTPServer: httpServerCfg,
		Postgres:   pgCfg,
	}, nil
}

// loadDatabaseURL loads the database URL from the environment variables.
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
