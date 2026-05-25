.PHONY: run build tidy sqlc migrate-up migrate-down migrate-create test

APP_NAME=vide-be
MIGRATE_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

run:
	go run ./cmd/api

# Ex: make build APP_NAME=vide-be
build: 
	go build -o bin/$(APP_NAME) ./cmd/api

tidy:
	go mod tidy

sqlc:
	sqlc generate

migrate-up:
	migrate -path db/migrations -database "$(MIGRATE_DSN)" up

migrate-down:
	migrate -path db/migrations -database "$(MIGRATE_DSN)" down 1

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

test:
	go test ./... -v -cover
