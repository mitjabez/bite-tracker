generate:
	go tool templ generate

db-init:
	docker exec -i bite-tracker-db-1 psql -U biteapp -d bite_tracker < db/schema.sql

db-start:
	docker compose up -d

db-stop:
	docker compose stop

db-delete:
	docker compose down -v

run: generate
	go run cmd/main.go

