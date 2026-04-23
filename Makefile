# Makefile for Smart Outgoing Demo

.PHONY: help build run test lint clean fmt docker-build docker-run dev setup

# Default target
help:
	@echo "Available commands:"
	@echo "  setup      - Install dependencies and setup development environment"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  dev        - Run in development mode with hot reload"
	@echo "  test       - Run all tests"
	@echo "  test-cover - Run tests with coverage"
	@echo "  lint       - Run linter"
	@echo "  fmt        - Format code"
	@echo "  clean      - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run application in Docker"

# Variables
BINARY_NAME=smart-outgoing-demo
DOCKER_IMAGE=smart-outgoing-demo:latest
GO_FILES=$(shell find . -name "*.go" -type f)
PORT?=8080

# Setup development environment
setup:
	@echo "Setting up development environment..."
	go mod download
	go install github.com/cweill/gotests/gotests@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from example"; fi

# Build the application
build: $(BINARY_NAME)

$(BINARY_NAME): $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) cmd/server/main.go

# Run the application
run: build
	@echo "Starting $(BINARY_NAME) on port $(PORT)..."
	./$(BINARY_NAME)

# Development mode with hot reload
dev:
	@echo "Starting development server with hot reload..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean -cache

# Generate tests
generate-tests:
	@echo "Generating unit tests..."
	gotests -all -w ./...

# Security scan
security:
	@echo "Running security scan..."
	gosec ./...

# Vulnerability check
vuln-check:
	@echo "Checking for vulnerabilities..."
	govulncheck ./...

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p $(PORT):8080 --env-file .env $(DOCKER_IMAGE)

# All checks (run before committing)
checks: fmt lint test security
	@echo "All checks completed successfully!"

# CI pipeline
ci: setup checks test-cover
	@echo "CI pipeline completed successfully!"

# Development workflow
dev-workflow: setup dev
	@echo "Development environment ready!"

# Production build
prod-build:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) cmd/server/main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Check for outdated dependencies
deps-outdated:
	@echo "Checking for outdated dependencies..."
	go list -u -m all

# Run benchmark tests
benchmark:
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem ./...

# Profile the application
profile:
	@echo "Profiling application..."
	go build -o $(BINARY_NAME) cmd/server/main.go
	./$(BINARY_NAME) -cpuprofile=cpu.prof -memprofile=mem.prof
	go tool pprof cpu.prof

# Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation available at http://localhost:6060"

# Database migrations (if applicable)
migrate-up:
	@echo "Running database migrations..."
	# Add migration commands here

migrate-down:
	@echo "Rolling back database migrations..."
	# Add rollback commands here

# Health check
health:
	@echo "Checking application health..."
	curl -f http://localhost:$(PORT)/api/config || echo "Health check failed"

# Load test
load-test:
	@echo "Running load tests..."
	# Add load test commands here

# Backup data
backup:
	@echo "Backing up data..."
	# Add backup commands here

# Restore data
restore:
	@echo "Restoring data..."
	# Add restore commands here
