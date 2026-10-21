package main

import (
	"bot/internal/actions"
	"bot/internal/agent"
	"bot/internal/api"
	"bot/internal/config"
	"bot/internal/llm"
	"bot/internal/telegram"
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
)

func main() {

	config, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := llm.NewClient(config.LLMAPIKey, config.LLMBaseURL, config.LLMModel, 10*time.Second)

	apiClient := api.NewClient(config.APIBaseURL, 10*time.Second)

	executor := actions.NewExecutor(apiClient)

	loop := agent.NewLoop(client, executor)

	location, err := time.LoadLocation(config.Timezone)
	if err != nil {
		log.Fatalf("load bot timezone: %v", err)
	}
	handler := telegram.NewHandler(loop, apiClient, location)

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.Handle),
		bot.WithDebug(),
	}

	b, err := bot.New(config.BotToken, opts...)
	if err != nil {
		log.Fatalf("create telegram bot: %v", err)
	}

	b.Start(ctx)
}
