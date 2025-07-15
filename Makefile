.PHONY: build test clean run help

# Default target
all: build

# Build the application
build:
	go build -o bin/strongdm .

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Run tests with detailed coverage report
test-coverage-html:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the application
run: build
	./bin/strongdm

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build multiplatform Docker image
docker-build:
	docker buildx build --platform linux/amd64,linux/arm64 -t strongdm:latest .

# Build and load Docker image for local testing
docker-build-local:
	docker buildx build --platform linux/amd64 -t strongdm:latest --load .

# Run Docker container
docker-run: docker-build-local
	docker run -p 8080:8080 -e BIND_ADDR=:8080 strongdm:latest

# Push multiplatform Docker image (requires login)
docker-push: docker-build
	docker buildx build --platform linux/amd64,linux/arm64 -t your-dockerhub-username/strongdm:latest --push .

# Help target
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  test            - Run all tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  test-coverage-html - Generate HTML coverage report"
	@echo "  clean           - Clean build artifacts"
	@echo "  run             - Build and run the application"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  deps            - Install dependencies"
	@echo "  docker-build    - Build multiplatform Docker image"
	@echo "  docker-build-local - Build Docker image for local testing"
	@echo "  docker-run      - Build and run Docker container"
	@echo "  docker-push     - Push multiplatform Docker image to registry"
	@echo "  help            - Show this help message"
