package configs

import (
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	ExternalAPI string
	LogLevel    string
}

func LoadConfig() (*Config, error) {
	serverPort := getEnv("SERVER_PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "")
	externalAPI := getEnv("EXTERNAL_API", "")
	logLevel := getEnv("LOG_LEVEL", "info")

	if databaseURL == "" {
		return nil, ErrMissingDatabaseURL
	}

	if externalAPI == "" {
		return nil, ErrMissingExternalAPI
	}

	cfg := &Config{
		ServerPort:  serverPort,
		DatabaseURL: databaseURL,
		ExternalAPI: externalAPI,
		LogLevel:    logLevel,
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

var (
	ErrMissingDatabaseURL = &ConfigError{"DATABASE_URL is required but not set"}
	ErrMissingExternalAPI = &ConfigError{"EXTERNAL_API is required but not set"}
)

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}
