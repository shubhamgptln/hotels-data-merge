GO          ?= go

build: ## Build binaries
	$(GO) build -ldflags '$(LDFLAGS)' -o ./build/api cmd/server.go

run-api: ## Run the server app
	go run cmd/server.go