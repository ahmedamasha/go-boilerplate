# Variables
BINARY_NAME=main
BUILD_DIR=build
DOCKER_IMAGE=go-app
DOCKER_TAG=latest

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the application
.PHONY: build
build:
	@echo "Building application..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

# Run the application locally
.PHONY: run
run:
	@echo "Running application..."
	$(GOCMD) run ./cmd/api

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Generate Wire code
.PHONY: wire
wire:
	@echo "Generating Wire code..."
	$(GOCMD) generate ./cmd/api

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	$(GOCMD) install github.com/google/wire/cmd/wire@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Start PostgreSQL with Docker
.PHONY: postgres
postgres:
	@echo "Starting PostgreSQL container..."
	docker run --name postgres-local \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=go_microservice \
		-p 5432:5432 \
		-d postgres:15-alpine

# Stop PostgreSQL container
.PHONY: postgres-stop
postgres-stop:
	@echo "Stopping PostgreSQL container..."
	docker stop postgres-local || true
	docker rm postgres-local || true

# Run database migrations
.PHONY: migrate-up
migrate-up:
	@echo "Running database migrations..."
	goose -dir scripts/migrations postgres "host=localhost port=5432 user=user password=password  dbname=mytestdb sslmode=disable" up

# Rollback database migrations
.PHONY: migrate-down
migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path scripts/migrations -database "postgres://postgres:postgres@localhost:5432/go_microservice?sslmode=disable" down

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run with Docker Compose
.PHONY: docker-up
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build

# Stop Docker Compose
.PHONY: docker-down
docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Run application with Docker Compose in background
.PHONY: docker-up-d
docker-up-d:
	@echo "Starting services with Docker Compose in background..."
	docker-compose up --build -d

# View Docker Compose logs
.PHONY: docker-logs
docker-logs:
	docker-compose logs -f

# Clean Docker resources
.PHONY: docker-clean
docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f

# Development setup (install tools, start postgres, run migrations)
.PHONY: dev-setup
dev-setup: install-tools postgres
	@echo "Waiting for PostgreSQL to be ready..."
	sleep 5
	@echo "Running migrations..."
	$(MAKE) migrate-up
	@echo "Development setup complete!"

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  wire          - Generate Wire code"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  install-tools - Install development tools"
	@echo "  postgres      - Start PostgreSQL container"
	@echo "  postgres-stop - Stop PostgreSQL container"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Start services with Docker Compose"
	@echo "  docker-down   - Stop Docker Compose services"
	@echo "  docker-up-d   - Start services with Docker Compose in background"
	@echo "  docker-logs   - View Docker Compose logs"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  dev-setup     - Complete development setup"
	@echo "  help          - Show this help message" 