build: ## Build binaries
	go build -o ./build/api cmd/server.go

run-api: ## Run the server app
	go run cmd/server.go