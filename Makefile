# Makefile for URL Shortener Service

.PHONY: build run test clean fmt deps

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application locally
run:
	go run ./cmd/server

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Format code
fmt:
	go fmt ./...

# Install dependencies
deps:
	go mod download
	go mod tidy
