# Makefile for URL Shortener Service

.PHONY: build run clean docker-build docker-run docker-run-detached docker-stop docker-logs redis-cli fmt deps

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application locally
run:
	go run ./cmd/server

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	/opt/homebrew/bin/docker compose-build

# Run with Docker Compose
docker-run:
	docker compose up

# Run with Docker Compose in background
docker-run-detached:
	docker compose up -d

# Stop Docker services
docker-stop:
	docker compose down

# View logs
docker-logs:
	docker compose logs -f

# Access Redis CLI
redis-cli:
	docker compose exec db redis-cli

# Format code
fmt:
	go fmt ./...

# Install dependencies
deps:
	go mod download
	go mod tidy
