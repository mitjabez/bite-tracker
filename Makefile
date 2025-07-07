generate:
	go tool templ generate

run: generate
	 o run cmd/main.go

