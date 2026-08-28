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
	// UserRateLimit is the per authenticated user limit for protected routes.
	UserRateLimit RateLimitConfig   `yaml:"user_rate_limiting"`
	Profiling     ProfilingConfig   `yaml:"profiling"`
	Maintenance   MaintenanceConfig `yaml:"maintenance"`
	Tracing       TracingConfig     `yaml:"tracing"`
	RBAC          RBACConfig        `yaml:"rbac"`
}

// RBACConfig — runtime-RBAC settings (rules reload from the DB).
type RBACConfig struct {
	// RefreshInterval — the background rule reload period from Postgres
	// (eventual consistency between instances; local edits
	// apply immediately).
	RefreshInterval Duration `yaml:"refresh_interval"`
}

// HTTPServerConfig is the HTTP server settings.
type HTTPServerConfig struct {
	Port         int      `yaml:"port"`
	ReadTimeout  Duration `yaml:"read_timeout"`
	WriteTimeout Duration `yaml:"write_timeout"`
	IdleTimeout  Duration `yaml:"idle_timeout"`
	// TrustedProxies is the CIDR/client networks allowed to forward client IP
	// headers (X-Forwarded-For / X-Real-IP) to the backend. It MUST list only
	// the reverse proxy (nginx) networks; trusting all networks lets remote
	// clients spoof their IP and bypass per-IP rate limiting.
	TrustedProxies []string `yaml:"trusted_proxies"`
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
	SecretKey string `yaml:"-"`
	// RefreshKey remains legacy-compatible: refresh tokens are opaque and stored
	// in the DB (AD-06); the key is only kept in the environment.
	RefreshKey          string   `yaml:"-"`
	AccessExpiry        Duration `yaml:"access_expiry"`
	RefreshExpiry       Duration `yaml:"refresh_expiry"`
	Issuer              string   `yaml:"issuer"`
	RefreshCookieSecure bool     `yaml:"refresh_cookie_secure"`
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

// MaintenanceConfig — background data normalization (periodic run of
// fn_normalize_employee_states for employee_states).
type MaintenanceConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Interval Duration `yaml:"interval"`
}

// TracingConfig — OpenTelemetry distributed tracing settings (OTLP/gRPC
// exporter, by default the in-stack Jaeger collector).
type TracingConfig struct {
	Enabled          bool    `yaml:"enabled"`
	ExporterEndpoint string  `yaml:"exporter_endpoint"`
	ServiceName      string  `yaml:"service_name"`
	SamplerRatio     float64 `yaml:"sampler_ratio"`
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

	// TRACING_ENDPOINT overrides the OTLP exporter endpoint (dev: host-run air
	// reaches the in-docker Jaeger collector via its published port; kept empty
	// in the full-stack docker run, which uses the in-network "jaeger:4317").
	if endpoint := getEnv("TRACING_ENDPOINT", ""); endpoint != "" {
		cfg.Tracing.ExporterEndpoint = endpoint
	}

	return nil
}
