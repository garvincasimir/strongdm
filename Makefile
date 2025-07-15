.PHONY: build test clean run help deps install-tools setup fmt lint staticcheck quality docker-build docker-build-local docker-run

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

# Run linter (installs golangci-lint if not present)
lint: install-tools
	$(GOBIN)/golangci-lint run

# Run staticcheck (installs staticcheck if not present)
staticcheck: install-tools
	$(GOBIN)/staticcheck ./...

# Run all quality checks
quality: fmt lint staticcheck test
	@echo "All quality checks passed!"

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@test -f $(GOBIN)/golangci-lint || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@test -f $(GOBIN)/staticcheck || (echo "Installing staticcheck..." && go install honnef.co/go/tools/cmd/staticcheck@latest)
	@echo "Development tools installed successfully!"

# Setup development environment
setup: deps install-tools
	@echo "Development environment setup complete!"

# Go binary path
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

# Extract Go version from go.mod
GO_VERSION := $(shell grep '^go ' go.mod | awk '{print $$2}')

# Build multiplatform Docker image
docker-build:
	docker buildx build --build-arg GO_VERSION=$(GO_VERSION) --platform linux/amd64,linux/arm64 -t strongdm:latest .

# Build and load Docker image for local testing
docker-build-local:
	docker buildx build --build-arg GO_VERSION=$(GO_VERSION) -t strongdm:latest --load .

# Run Docker container
docker-run: docker-build-local
	docker run -p 8080:8080 -e BIND_ADDR=:8080 strongdm:latest


# Help target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build & Run:"
	@echo "  build           - Build the application"
	@echo "  run             - Build and run the application"
	@echo "  clean           - Clean build artifacts"
	@echo ""
	@echo "Testing & Quality:"
	@echo "  test            - Run all tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  test-coverage-html - Generate HTML coverage report"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter (installs golangci-lint if needed)"
	@echo "  staticcheck     - Run staticcheck (installs staticcheck if needed)"
	@echo "  quality         - Run all quality checks (fmt, lint, staticcheck, test)"
	@echo ""
	@echo "Dependencies & Setup:"
	@echo "  deps            - Install Go dependencies"
	@echo "  install-tools   - Install development tools (golangci-lint, staticcheck)"
	@echo "  setup           - Setup complete development environment"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build    - Build multiplatform Docker image (uses Go version from go.mod)"
	@echo "  docker-build-local - Build Docker image for local testing (uses Go version from go.mod)"
	@echo "  docker-run      - Build and run Docker container"
	@echo "  docker-push     - Push multiplatform Docker image to registry"
	@echo ""
	@echo "  help            - Show this help message"
