build: ## Build binaries
	go build -a -installsuffix cgo -o ./build/api cmd/server.go

run-api: ## Run the server app
	go run cmd/server.go