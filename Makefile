.PHONY: build run test clean deps fmt lint

# Build the application
build:
	go build -o bin/sales-insights cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run ./...

# Run database migrations
migrate:
	go run cmd/api/main.go migrate 