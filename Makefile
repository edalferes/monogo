# Simple Makefile for Monetics project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=./bin/monetics
CMD_PATH=./cmd/api

.PHONY: help build test clean run docker-build swagger lint lint-fix

## Help
help: ## Show available commands
	@echo 'Available commands:'
	@echo '  build         - Build the application'
	@echo '  test          - Run tests'
	@echo '  test-coverage - Run tests with coverage'
	@echo '  test-verbose  - Run tests with verbose output'
	@echo '  run           - Run the application'
	@echo '  clean         - Clean build artifacts'
	@echo '  docker-build  - Build Docker image'
	@echo '  swagger       - Generate Swagger documentation'
	@echo '  lint          - Run golangci-lint'
	@echo '  lint-fix      - Fix auto-fixable lint issues'

## Build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p ./bin
	@$(GOBUILD) -o $(BINARY_NAME) $(CMD_PATH)

## Swagger
swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@$(GOCMD) run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go --parseDependency --parseInternal -o docs/openapi

## Test
test: ## Run tests
	@echo "Running tests..."
	@$(GOTEST) ./...

## Test with coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@$(GOTEST) -coverprofile=coverage.out -covermode=atomic ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## Test with JUnit output
test-junit: ## Run tests with JUnit output (requires gotestsum)
	@echo "Running tests with JUnit output..."
	@gotestsum --junitfile tests.xml --format testname -- -coverprofile=coverage.out -covermode=atomic ./...

## Test verbose
test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@$(GOTEST) -v ./...

## Lint
lint: ## Run golangci-lint
	@echo "Running linter..."
	@$(shell go env GOPATH)/bin/golangci-lint run --config=.golangci.yml

## Fix lint issues
lint-fix: ## Fix auto-fixable lint issues
	@echo "Fixing lint issues..."
	@$(shell go env GOPATH)/bin/golangci-lint run --fix --config=.golangci.yml

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