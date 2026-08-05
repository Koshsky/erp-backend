// Package config — application configuration.
//
// Non-secret settings are read from config.yaml (the path is set by the
// CONFIG_PATH environment variable, defaulting to ./config.yaml). Secrets and
// the DB URL come only from environment variables: DATABASE_URL, JWT_SECRET_KEY,
// JWT_REFRESH_KEY.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration is a Go-formatted duration (e.g., "15s", "30m") from YAML.
type Duration time.Duration

// UnmarshalYAML parses a duration string into [Duration].
func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	parsed, err := time.ParseDuration(node.Value)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", node.Value, err)
	}
	*d = Duration(parsed)
	return nil
}

// Config is the root application configuration.
type Config struct {
	HTTPServer HTTPServerConfig `yaml:"server"`
	Postgres   PostgresConfig   `yaml:"postgres"`
	JWT        JWTConfig        `yaml:"jwt"`
	Logging    LoggingConfig    `yaml:"logging"`
	Swagger    SwaggerConfig    `yaml:"swagger"`
	CORS       CORSConfig       `yaml:"cors"`
	RateLimit  RateLimitConfig  `yaml:"rate_limiting"`
	Profiling  ProfilingConfig  `yaml:"profiling"`
}

// HTTPServerConfig is the HTTP server settings.
type HTTPServerConfig struct {
	Port         int      `yaml:"port"`
	ReadTimeout  Duration `yaml:"read_timeout"`
	WriteTimeout Duration `yaml:"write_timeout"`
	IdleTimeout  Duration `yaml:"idle_timeout"`
}

// PostgresConfig is the Postgres connection pool settings.
//
// URL comes from the DATABASE_URL environment variable (yaml:"-").
type PostgresConfig struct {
	URL               string   `yaml:"-"`
	MaxConns          int32    `yaml:"max_conns"`
	MinConns          int32    `yaml:"min_conns"`
	MaxConnLifetime   Duration `yaml:"max_conn_lifetime"`
	MaxConnIdleTime   Duration `yaml:"max_conn_idle_time"`
	HealthCheckPeriod Duration `yaml:"health_check_period"`
	ConnectTimeout    Duration `yaml:"connect_timeout"`
}

// JWTConfig is the JWT settings.
//
// Secrets come from environment variables (yaml:"-").
type JWTConfig struct {
	SecretKey     string   `yaml:"-"`
	RefreshKey    string   `yaml:"-"`
	AccessExpiry  Duration `yaml:"access_expiry"`
	RefreshExpiry Duration `yaml:"refresh_expiry"`
	Issuer        string   `yaml:"issuer"`
}

// LoggingConfig is the logging settings.
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// SwaggerConfig is the embedded Swagger UI settings.
type SwaggerConfig struct {
	Enabled bool `yaml:"enabled"`
}

// CORSConfig is the Cross-Origin Resource Sharing settings.
type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           Duration `yaml:"max_age"`
}

// ProfilingConfig is the pprof server settings.
type ProfilingConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
}

// RateLimitConfig is the per-client rate limiting settings.
type RateLimitConfig struct {
	Enabled           bool     `yaml:"enabled"`
	RequestsPerSecond float64  `yaml:"requests_per_second"`
	Burst             int      `yaml:"burst"`
	CleanupInterval   Duration `yaml:"cleanup_interval"`
	Expiration        Duration `yaml:"expiration"`
}

// Load loads the configuration: config.yaml plus secrets from environment variables.
func Load() (*Config, error) {
	path := getEnv("CONFIG_PATH", "config.yaml")

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}

	var cfg Config
	if parseErr := yaml.Unmarshal(raw, &cfg); parseErr != nil {
		return nil, fmt.Errorf("parse config file %s: %w", path, parseErr)
	}

	if envErr := applyEnv(&cfg); envErr != nil {
		return nil, envErr
	}

	return &cfg, nil
}

// applyEnv applies secrets and the DB URL from environment variables to the configuration.
func applyEnv(cfg *Config) error {
	cfg.Postgres.URL = getEnv("DATABASE_URL", "")
	if cfg.Postgres.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	cfg.JWT.SecretKey = getEnv("JWT_SECRET_KEY", "")
	if cfg.JWT.SecretKey == "" {
		return fmt.Errorf("JWT_SECRET_KEY is required")
	}

	cfg.JWT.RefreshKey = getEnv("JWT_REFRESH_KEY", "")
	if cfg.JWT.RefreshKey == "" {
		return fmt.Errorf("JWT_REFRESH_KEY is required")
	}

	return nil
}
