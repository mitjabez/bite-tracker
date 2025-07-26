# .PHONY: generate db-init db-start db-stop db-delete

generate:
	go tool templ generate

db-init:
	docker exec -i bite-tracker-db-1 psql -U biteapp -d bite_tracker < internal/db/schema.sql

db-start:
	docker compose up -d

db-stop:
	docker compose stop

db-delete:
	docker compose down -v

build: generate
	go build -o build/bite-tracker cmd/bitetracker/main.go

run: build
	build/bite-tracker

