package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL   string
	LogLevel      string
	LogFormat     string
	SwaggerEnable bool
}

func Load() (*Config, error) {
	dbURL, err := loadDatabaseURL()
	if err != nil {
		return nil, err
	}

	logLevel, err := requiredEnv("LOG_LEVEL")
	if err != nil {
		return nil, err
	}

	logFormat, err := requiredEnv("LOG_FORMAT")
	if err != nil {
		return nil, err
	}

	swaggerEnable := os.Getenv("SWAGGER_ENABLE") == "true"

	return &Config{
		DatabaseURL:   dbURL,
		LogLevel:      logLevel,
		LogFormat:     logFormat,
		SwaggerEnable: swaggerEnable,
	}, nil
}

func loadDatabaseURL() (string, error) {
	user, err := requiredEnv("DATABASE_USER")
	if err != nil {
		return "", err
	}

	password, err := requiredEnv("DATABASE_PASSWORD")
	if err != nil {
		return "", err
	}

	host, err := requiredEnv("DATABASE_HOST")
	if err != nil {
		return "", err
	}

	port, err := requiredEnv("DATABASE_PORT")
	if err != nil {
		return "", err
	}

	name, err := requiredEnv("DATABASE_NAME")
	if err != nil {
		return "", err
	}

	sslMode, err := requiredEnv("DATABASE_SSL_MODE")
	if err != nil {
		return "", err
	}

	built := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   name,
	}

	query := built.Query()
	query.Set("sslmode", sslMode)
	built.RawQuery = query.Encode()

	return built.String(), nil
}

func requiredEnv(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", fmt.Errorf("%s environment variable is required", key)
	}

	return value, nil
}
