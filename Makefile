.PHONY: help env run debug stop test docker-build-dev docker-run-dev docker-debug docker-stop docker-clean docker-logs docker-shell docker-test docker-lint docker-build-prod docker-run-prod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env: ## Create .env file from .env.example
	cp .env.example .env
	@echo ".env file created. Edit it to match your MIDI device configuration."

# Native Development Commands (for full MIDI access on macOS)

run: ## Run with Air hot reload
	@echo "🎵 Starting development with Air hot reload..."
	@echo "MIDI devices: Full access to all devices"
	@echo ""
	@if [ -f .env ]; then set -a && . ./.env && set +a; fi && air -c .air.native.toml

debug: ## Run with Air hot reload + Delve debugger (port 2346)
	@echo "🎵 Starting development with Air + Delve..."
	@echo "Debugger listening on localhost:2346"
	@echo "MIDI devices: Full access to all devices"
	@echo ""
	@mkdir -p tmp
	@printf '#!/bin/bash\nexec dlv exec --headless --listen=:2346 --api-version=2 --accept-multiclient --continue ./tmp/main\n' > tmp/debug.sh
	@chmod +x tmp/debug.sh
	@if [ -f .env ]; then set -a && . ./.env && set +a; fi && air -c .air.debug.toml

stop: ## Stop any running Delve debugger on port 2346
	@echo "Stopping any Delve processes on port 2346..."
	@pkill -f "dlv.*2346" 2>/dev/null || true
	@lsof -ti :2346 | xargs kill -9 2>/dev/null || true
	@echo "✅ Stopped"

test: ## Run tests
	go test ./...

# Docker Development Commands

docker-build-dev: ## Build development Docker image
	docker compose build reactor-dev

docker-run-dev: ## Start development environment with hot reload
	docker compose up reactor-dev

docker-debug: ## Start development environment in background and attach debugger
	docker compose up reactor-dev
	@echo "Debugger listening on localhost:2345"
	@echo "Connect your IDE to localhost:2345"

docker-stop: ## Stop all containers
	docker compose stop

docker-clean: ## Remove containers, volumes, and built binaries
	docker compose down -v
	rm -rf tmp/

docker-logs: ## Follow logs from development container
	docker compose logs -f reactor-dev

docker-shell: ## Open shell in running development container
	docker exec -it reactor-dev /bin/sh

docker-test: ## Run tests in container
	docker compose run --rm reactor-dev go test ./...

docker-lint: ## Run linter in container
	docker compose run --rm reactor-dev golangci-lint run

docker-build-prod: ## Build production Docker image
	docker build --target production -t reactor:latest .

docker-run-prod: ## Run production container
	docker run --rm reactor:latest
