build:
	go tool templ generate
	npx tailwindcss -i views/assets/css/input.css -o ./views/assets/css/tailwind.css

run: build
	go run cmd/main.go

