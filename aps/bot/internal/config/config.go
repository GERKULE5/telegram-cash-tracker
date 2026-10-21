package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken   string
	APIBaseURL string
	LLMBaseURL string
	LLMAPIKey  string
	LLMModel   string
	Timezone   string
}

func Load() (Config, error) {
	if envFile := os.Getenv("ENV_FILE"); envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return Config{}, fmt.Errorf("load env file: %w", err)
		}
	}

	config := Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		APIBaseURL: os.Getenv("API_BASE_URL"),
		LLMBaseURL: os.Getenv("LLM_BASE_URL"),
		LLMAPIKey:  os.Getenv("LLM_API_KEY"),
		LLMModel:   os.Getenv("LLM_MODEL"),
		Timezone:   os.Getenv("BOT_TIMEZONE"),
	}
	if config.Timezone == "" {
		config.Timezone = "UTC"
	}

	required := []struct {
		name  string
		value string
	}{
		{"BOT_TOKEN", config.BotToken},
		{"API_BASE_URL", config.APIBaseURL},
		{"LLM_BASE_URL", config.LLMBaseURL},
		{"LLM_API_KEY", config.LLMAPIKey},
		{"LLM_MODEL", config.LLMModel},
	}

	for _, field := range required {
		if field.value == "" {
			return Config{}, fmt.Errorf("%s is required", field.name)
		}
	}

	return config, nil
}
