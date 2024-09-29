version = $(shell go env GOVERSION)

build: ## Build binaries
	go build -a -installsuffix cgo -o ./build/api cmd/server.go

unit-tests:
	go test -v ./...


run-api: ## Run the server app
	go run cmd/server.go