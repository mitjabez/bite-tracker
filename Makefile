.PHONY: generate db-up db-down db-start db-stop build

DB_URL = postgres://biteapp:superburrito@localhost:5432/bite_tracker?sslmode=disable
DB_PATH = internal/db/migrations
BIN = build/bite-tracker

generate:
	templ generate
	sqlc generate

db-up:
	migrate -database $(DB_URL) -path $(DB_PATH) up

db-down:
	migrate -database $(DB_URL) -path $(DB_PATH) down -all

db-start:
	docker compose up -d

db-stop:
	docker compose stop

build:
	go build -o $(BIN) cmd/bitetracker/main.go

run: build
	$(BIN)
