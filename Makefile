.PHONY: help build-dev run-dev debug stop clean logs shell test lint build-prod run-prod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build-dev: ## Build development Docker image
	docker-compose build reactor-dev

run-dev: ## Start development environment with hot reload
	docker-compose up reactor-dev

debug: ## Start development environment in background and attach debugger
	docker-compose up -d reactor-dev
	@echo "Debugger listening on localhost:2345"
	@echo "Connect your IDE to localhost:2345"

stop: ## Stop all containers
	docker-compose down

clean: ## Remove containers, volumes, and built binaries
	docker-compose down -v
	rm -rf tmp/

logs: ## Follow logs from development container
	docker-compose logs -f reactor-dev

shell: ## Open shell in running development container
	docker exec -it reactor-dev /bin/sh

test: ## Run tests in container
	docker-compose run --rm reactor-dev go test ./...

lint: ## Run linter in container
	docker-compose run --rm reactor-dev golangci-lint run

build-prod: ## Build production Docker image
	docker build --target production -t reactor:latest .

run-prod: ## Run production container
	docker run --rm reactor:latest

env: ## Create .env file from .env.example
	cp .env.example .env
	@echo ".env file created. Edit it to match your MIDI device configuration."
