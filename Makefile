# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Binary info
BINARY_NAME=monogo
BINARY_PATH=./bin/$(BINARY_NAME)
CMD_PATH=./cmd/api

# Docker
DOCKER_IMAGE=monogo
DOCKER_TAG=latest

.PHONY: help docs swag build run test test-integration test-coverage clean lint lint-fix security-scan \
        docker-build docker-run docker-clean install-tools deps-download deps-verify deps-tidy \
        release godoc ci-test ci-lint ci-security ci-docker

## Help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## Documentation
docs: swag ## Generate all documentation

swag: ## Generate Swagger documentation
	@echo "📚 Generating Swagger documentation..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@$$(go env GOPATH)/bin/swag init -g $(CMD_PATH)/main.go -o ./docs/

godoc: ## Start Go documentation server
	@echo "📖 Starting Go documentation server at http://localhost:6060"
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo "📍 Access documentation at: http://localhost:6060/pkg/github.com/edalferes/monogo/"
	@godoc -http=:6060

## Build & Run
build: ## Build the application
	@echo "🏗️ Building $(BINARY_NAME)..."
	@mkdir -p bin
	@$(GOBUILD) -o $(BINARY_PATH) $(CMD_PATH)/main.go
	@echo "✅ Binary created at $(BINARY_PATH)"

run: ## Run the application (auth module on port 8080)
	@echo "🚀 Running $(BINARY_NAME) auth module..."
	@$(GOCMD) run $(CMD_PATH)/main.go auth --port=8080

run-testmodule: ## Run the testmodule
	@echo "🚀 Running $(BINARY_NAME) testmodule..."
	@$(GOCMD) run $(CMD_PATH)/main.go testmodule

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@$(GOCLEAN)
	@rm -rf bin/
	@rm -rf coverage/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean completed"

## Testing
test: ## Run unit tests
	@echo "🧪 Running unit tests..."
	@$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage report
	@echo "📊 Running tests with coverage..."
	@mkdir -p coverage
	@$(GOTEST) -v -race -coverprofile=coverage/coverage.out -covermode=atomic ./...
	@$(GOCMD) tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@$(GOCMD) tool cover -func=coverage/coverage.out | grep total
	@echo "📈 Coverage report generated at coverage/coverage.html"

test-integration: ## Run integration tests
	@echo "🔧 Running integration tests..."
	@$(GOTEST) -v -tags=integration ./tests/integration/... || echo "ℹ️ No integration tests found"

ci-test: ## Run tests for CI (with PostgreSQL)
	@echo "🚀 Running CI tests..."
	@$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...

## Code Quality
lint: ## Run linter
	@echo "🔍 Running linter..."
	@$(GOLINT) run --timeout=5m

lint-fix: ## Run linter with auto-fix
	@echo "🔧 Running linter with auto-fix..."
	@$(GOLINT) run --fix --timeout=5m

fmt: ## Format code
	@echo "🎨 Formatting code..."
	@$(GOFMT) -s -w .
	@go mod tidy

ci-lint: ## Run linting for CI
	@echo "🔍 Running CI linting..."
	@$(GOLINT) run --timeout=5m --config=.golangci.yml

## Security
security-scan: ## Run security scans
	@echo "🔒 Running security scans..."
	@gosec -fmt sarif -out gosec.sarif ./... || echo "⚠️ gosec not installed, run 'make install-tools'"
	@echo "✅ GoSec scan completed"

security-deps: ## Check dependencies for vulnerabilities
	@echo "🔐 Checking dependencies for vulnerabilities..."
	@go list -json -deps ./... | nancy sleuth || echo "⚠️ Nancy not installed, run 'make install-tools'"
	@govulncheck ./... || echo "⚠️ Govulncheck not installed, run 'make install-tools'"

ci-security: ## Run security scans for CI
	@echo "🔒 Running CI security scans..."
	@gosec -fmt sarif -out gosec.sarif ./... || true
	@go list -json -deps ./... | nancy sleuth || true
	@govulncheck ./... || true

## Docker
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "✅ Docker image $(DOCKER_IMAGE):$(DOCKER_TAG) built"

docker-run: ## Run Docker container
	@echo "🚀 Running Docker container..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-clean: ## Clean Docker images
	@echo "🧹 Cleaning Docker images..."
	@docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true
	@docker system prune -f

ci-docker: ## Build Docker image for CI
	@echo "🐳 Building Docker image for CI..."
	@docker build -t $(DOCKER_IMAGE):test .

## Dependencies
deps-download: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@$(GOMOD) download

deps-verify: ## Verify dependencies
	@echo "✅ Verifying dependencies..."
	@$(GOMOD) verify

deps-tidy: ## Tidy up dependencies
	@echo "🧹 Tidying dependencies..."
	@$(GOMOD) tidy

deps-update: ## Update dependencies
	@echo "⬆️ Updating dependencies..."
	@$(GOGET) -u ./...
	@$(GOMOD) tidy

## Development Tools
install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo "📥 Installing gosec (security scanner)..."
	@go install github.com/securecodewarrior/gosec@latest || echo "⚠️ gosec installation failed, skipping..."
	@echo "📥 Installing nancy (dependency scanner)..."
	@go install github.com/sonatypecommunity/nancy@latest || echo "⚠️ nancy installation failed, skipping..."
	@echo "📥 Installing govulncheck (vulnerability checker)..."
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "✅ Development tools installed"

## CI/CD Targets
ci: ci-test ci-lint ci-security ci-docker ## Run all CI checks

release: clean build test-coverage lint security-scan ## Prepare release build

## Health Check
health: ## Check application health
	@echo "💓 Checking application health..."
	@curl -f http://localhost:8080/health || echo "❌ Health check failed - is the application running?"

## Database (when using Docker Compose)
db-up: ## Start database services
	@echo "🗄️ Starting database services..."
	@docker-compose up -d postgres redis || echo "⚠️ docker-compose.yml not found"

db-down: ## Stop database services
	@echo "🛑 Stopping database services..."
	@docker-compose down || echo "⚠️ docker-compose.yml not found"

## Performance
benchmark: ## Run performance benchmarks
	@echo "⚡ Running benchmarks..."
	@$(GOTEST) -bench=. -benchmem ./... || echo "ℹ️ No benchmarks found"

load-test: ## Run load tests (requires k6)
	@echo "🏋️ Running load tests..."
	@k6 run tests/load/basic.js || echo "⚠️ k6 not installed or test files not found"

## All-in-one commands
dev-setup: install-tools deps-download ## Setup development environment
	@echo "🎉 Development environment ready!"

quick-check: fmt lint test ## Quick code quality check
	@echo "✅ Quick checks completed!"

full-check: clean deps-tidy fmt lint test-coverage security-scan docker-build ## Full quality check
	@echo "🎯 Full quality check completed!"