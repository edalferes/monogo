# Simple Makefile for Monetics project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=./bin/monetics
CMD_PATH=./cmd/api

.PHONY: help build test clean run docker-build swagger

## Help
help: ## Show available commands
	@echo 'Available commands:'
	@echo '  build       - Build the application'
	@echo '  test        - Run tests'
	@echo '  run         - Run the application'
	@echo '  clean       - Clean build artifacts'
	@echo '  docker-build- Build Docker image'
	@echo '  swagger     - Generate Swagger documentation'

## Build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BINARY_NAME) $(CMD_PATH)

## Swagger
swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@$(GOCMD) run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go --parseDependency --parseInternal -o docs/openapi

## Test
test: ## Run tests
	@echo "Running tests..."
	@$(GOTEST) ./...

## Run
run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@$(GOCMD) run $(CMD_PATH)

## Clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)

## Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME) .