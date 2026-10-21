export ENV_FILE := $(CURDIR)/.env
name ?= 

.PHONY: run-api
run-api:
	cd aps/api && go run ./cmd/main.go

.PHONY: run-web
run-web:
	cd aps/web && bun --env-file=../../.env run dev

.PHONY: run-bot
run-bot:
	cd aps/bot && go run ./cmd/main.go

.PHONY: up
up:
	docker compose up -d

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir ./aps/api/migrations ${name}

.PHONY: migrate-up
migrate-up:
	docker compose run --rm migrate up

.PHONY: migrate-down
migrate-down:
	docker compose run --rm migrate down 1

.PHONY: swagger
swagger:
	cd aps/api && go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g main.go -d cmd,internal -o docs --parseInternal
