.PHONY: docs swag build run test help

docs: swag

swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	$$(go env GOPATH)/bin/swag init -g cmd/api/main.go


build:
	go build -o bin/app cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test ./... -cover

help:
	@echo "Available commands:"
	@echo "  make docs  - Generate Swagger documentation"
	@echo "  make build - Build the application"
	@echo "  make run   - Run the application"
	@echo "  make test  - Run tests with coverage"