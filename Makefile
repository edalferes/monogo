.PHONY: docs swag godoc help

docs: swag

swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	$$(go env GOPATH)/bin/swag init -g cmd/api/main.go

godoc:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo "Starting godoc server on http://localhost:6060"
	@echo "Access the documentation at: http://localhost:6060/pkg/github.com/edalferes/monogo/"
	@$$(go env GOPATH)/bin/godoc -http=:6060

build:
	go build -o bin/app cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test ./... -cover

help:
	@echo "Available commands:"
	@echo "  make docs  - Generate Swagger documentation"
	@echo "  make godoc - Start godoc server on :6060"
	@echo "  make build - Build the application"
	@echo "  make run   - Run the application"
	@echo "  make test  - Run tests with coverage"