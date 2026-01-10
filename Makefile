.PHONY: help build-dev run-dev debug stop clean logs shell test lint build-prod run-prod run-native debug-native stop-native test-native

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build-dev: ## Build development Docker image
	docker compose build reactor-dev

run-dev: ## Start development environment with hot reload
	docker compose up reactor-dev

debug: ## Start development environment in background and attach debugger
	docker compose up reactor-dev
	@echo "Debugger listening on localhost:2345"
	@echo "Connect your IDE to localhost:2345"

stop: ## Stop all containers
	docker compose stop

clean: ## Remove containers, volumes, and built binaries
	docker compose down -v
	rm -rf tmp/

logs: ## Follow logs from development container
	docker compose logs -f reactor-dev

shell: ## Open shell in running development container
	docker exec -it reactor-dev /bin/sh

test: ## Run tests in container
	docker compose run --rm reactor-dev go test ./...

lint: ## Run linter in container
	docker compose run --rm reactor-dev golangci-lint run

build-prod: ## Build production Docker image
	docker build --target production -t reactor:latest .

run-prod: ## Run production container
	docker run --rm reactor:latest

env: ## Create .env file from .env.example
	cp .env.example .env
	@echo ".env file created. Edit it to match your MIDI device configuration."

# Native Development Commands (for full MIDI access on macOS)

run-native: ## Run natively with Air hot reload
	@echo "🎵 Starting native development with Air hot reload..."
	@echo "MIDI devices: Full access to all devices"
	@echo ""
	air -c .air.native.toml

debug-native: ## Run natively with Air hot reload + Delve debugger (port 2346)
	@echo "🎵 Starting native development with Air + Delve..."
	@echo "Debugger listening on localhost:2346"
	@echo "MIDI devices: Full access to all devices"
	@echo ""
	@mkdir -p tmp
	@printf '#!/bin/bash\nexec dlv exec --headless --listen=:2346 --api-version=2 --accept-multiclient --continue ./tmp/main\n' > tmp/debug.sh
	@chmod +x tmp/debug.sh
	air -c .air.debug.toml

stop-native: ## Stop any running Delve debugger on port 2346
	@echo "Stopping any Delve processes on port 2346..."
	@pkill -f "dlv.*2346" 2>/dev/null || true
	@lsof -ti :2346 | xargs kill -9 2>/dev/null || true
	@echo "✅ Stopped"

test-native: ## Run tests natively
	go test ./...
