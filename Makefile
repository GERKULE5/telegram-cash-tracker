export ENV_FILE := $(CURDIR)/.env

.PHONY: run-api
run-api:
	cd aps/api && go run ./cmd/main.go

.PHONY: run-web
run-web:
	cd aps/web && bun run dev

.PHONY: run-bot
run-bot:
	cd aps/bot && go run ./cmd/main.go

.PHONY: up
up:
	docker compose up -d