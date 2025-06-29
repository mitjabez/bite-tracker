generate:
	go tool templ generate
	npx tailwindcss -i views/assets/css/input.css -o ./views/assets/css/tailwind.css

run: generate
	 o run cmd/main.go

