
# StrongDM Service

A Go HTTP service with rate limiting using a leaky bucket algorithm.

## Quick Start

### Prerequisites
- Go 1.22+
- Docker (optional)

### Run Locally
```bash
git clone https://github.com/garvincasimir/strongdm.git
cd strongdm
go mod download
export BIND_ADDR=":8080"
go run main.go
```

### Test
```bash
curl http://localhost:8080
```

## Docker

```bash
# Build
docker build -t strongdm .

# Run
docker run -p 8080:8080 -e BIND_ADDR=":8080" strongdm
```

## Development

```bash
# Run tests
go test ./...

# Build
go build -o bin/strongdm .

# Format & lint
make quality
```

## API

**GET /** - Rate limited endpoint (120 requests/minute per IP)

Returns JSON with rate limit information:
```json
{
  "bucket": "192.168.1.1",
  "resetAt": "2025-07-15T10:30:00Z",
  "bucketSize": 2,
  "remaining": 1,
  "allowed": true
}
```

## CI/CD

- **Pull Requests**: Run tests
- **Main Branch**: Build and push Docker images to ECR
