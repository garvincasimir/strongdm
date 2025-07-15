
# StrongDM Service

A Go-based HTTP service built with modern DevOps practices including automated testing, Docker containerization, and CI/CD pipeline.

## 🚀 Quick Start

### Prerequisites

- Go 1.21.0 or later
- Docker
- Make (optional, for convenience commands)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/garvincasimir/strongdm.git
   cd strongdm
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the service**
   ```bash
   export BIND_ADDR=":8080"
   go run main.go
   ```

4. **Test the service**
   ```bash
   curl http://localhost:8080
   ```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

### Test Structure

- `*_test.go` files contain unit tests
- Tests use Go's standard testing package
- Coverage reports are generated during CI/CD

## 🐳 Docker

### Building the Image

```bash
# Build for local architecture
docker build -t strongdm .

# Build for multiple platforms
docker buildx build --platform linux/amd64,linux/arm64 -t strongdm .
```

### Running with Docker

```bash
# Run the container
docker run -p 8080:8080 -e BIND_ADDR=":8080" strongdm

# Run in the background
docker run -d -p 8080:8080 -e BIND_ADDR=":8080" strongdm
```

## 🔧 Configuration

The service uses environment variables for configuration:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `BIND_ADDR` | Address and port to bind the server | - | Yes |

## 📁 Project Structure

```
.
├── .github/workflows/    # GitHub Actions CI/CD
├── scripts/             # Build and utility scripts
├── main.go             # Main application entry point
├── handler.go          # HTTP request handlers
├── handler_test.go     # Handler tests
├── counter.go          # Counter functionality
├── counter_test.go     # Counter tests
├── bucket.go           # Bucket functionality
├── bucket_test.go      # Bucket tests
├── go.mod              # Go module definition
├── .go-version         # Go version specification
├── Dockerfile          # Container image definition
├── Makefile           # Build automation
└── README.md          # This file
```

## 🚀 CI/CD Pipeline

The project uses GitHub Actions for continuous integration and deployment:

### On Pull Requests
- ✅ Run automated tests
- ✅ Generate test coverage reports
- ✅ Upload coverage to Codecov

### On Main Branch Pushes
- ✅ Run tests
- ✅ Build multiplatform Docker images (AMD64/ARM64)
- ✅ Push images to Amazon ECR

### Required Secrets

To enable the full CI/CD pipeline, configure these repository secrets:

| Secret | Description |
|--------|-------------|
| `AWS_ACCESS_KEY_ID` | AWS access key for ECR |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key for ECR |

### Required Variables

| Variable | Description |
|----------|-------------|
| `AWS_REGION` | AWS region for ECR |
| `ECR_REPOSITORY` | ECR repository name |

## 🛠️ Development Workflow

### 1. Setting Up Your Environment

```bash
# Ensure you have the correct Go version
go version  # Should match version in .go-version

# Install development dependencies
go mod tidy
```

### 2. Making Changes

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make your changes
3. Add tests for new functionality
4. Run tests locally: `go test ./...`
5. Commit your changes: `git commit -m "feat: your feature description"`
6. Push to GitHub: `git push origin feature/your-feature`
7. Create a pull request

### 3. Code Quality

- Write tests for all new functionality
- Follow Go best practices and conventions
- Use `go fmt` to format your code
- Run `go vet` to check for common errors
- Ensure all tests pass before submitting PR

## 📦 Version Management

Go version is managed centrally via `go.mod` as the single source of truth:

- **`go.mod`**: Contains the Go version (e.g., `go 1.22`)
- **`Dockerfile`**: Extracts version from `go.mod` via build argument
- **GitHub Actions**: Reads version from `go.mod`

### Updating Go Version

Simply update the `go` directive in `go.mod`:

```go
module strongdm

go 1.23  // Change this line
```

The CI/CD pipeline will automatically use the updated version for:
- Running tests
- Building Docker images
- Multi-platform builds

## 🚢 Deployment

### Running in Production

The service is designed to run in containerized environments:

```bash
# Using Docker
docker run -p 8080:8080 -e BIND_ADDR=":8080" your-registry/strongdm:latest

# Using Docker Compose
version: '3.8'
services:
  strongdm:
    image: your-registry/strongdm:latest
    ports:
      - "8080:8080"
    environment:
      - BIND_ADDR=:8080
```

### Health Checks

The service exposes HTTP endpoints that can be used for health checks in orchestration platforms like Kubernetes.

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### Getting Started
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

### Code Style
- Follow standard Go conventions
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused

### Pull Request Process
1. Ensure all tests pass
2. Update documentation if needed
3. Add a clear description of your changes
4. Link any related issues

### Reporting Issues
- Use GitHub Issues for bug reports and feature requests
- Provide clear reproduction steps for bugs
- Include relevant logs and error messages

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

- 📖 Documentation: Check this README and code comments
- 🐛 Bug Reports: Use GitHub Issues
- 💡 Feature Requests: Use GitHub Issues
- 💬 Questions: Use GitHub Discussions

---

Built with ❤️ for the StrongDM team
