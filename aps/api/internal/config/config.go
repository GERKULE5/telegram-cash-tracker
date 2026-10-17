package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	PostgresHost     string
	PostgresPort     string
	PostgresUsername string
	PostgresPassword string
	PostgresDB       string
}

func Load() (Config, error) {
	if envFile := os.Getenv("ENV_FILE"); envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return Config{}, fmt.Errorf("failed load env %v", err)
		}
	}

	config := Config{
		Port:             os.Getenv("API_PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUsername: os.Getenv("POSTGRES_USERNAME"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
	}

	required := []struct {
		name  string
		value string
	}{
		{"API_PORT", config.Port},
		{"POSTGRES_HOST", config.PostgresHost},
		{"POSTGRES_PORT", config.PostgresPort},
		{"POSTGRES_USERNAME", config.PostgresUsername},
		{"POSTGRES_PASSWORD", config.PostgresPassword},
		{"POSTGRES_DB", config.PostgresDB},
	}

	for _, field := range required {
		if field.value == "" {
			return Config{}, fmt.Errorf("%s is required", field.name)
		}
	}

	return config, nil
}
