.PHONY: docs swag

docs: swag

swag:
	rm -rf docs
	go get github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/api/main.go
run:
	go run cmd/api/main.go
test:
	go test ./... -cover