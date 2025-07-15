# Build stage
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

# Install ca-certificates for SSL/TLS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build args for cross-compilation
ARG TARGETOS
ARG TARGETARCH

# Build the application for the target platform
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -installsuffix cgo -ldflags='-w -s' -o strongdm .

# Final stage
FROM scratch

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/strongdm /strongdm

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/strongdm"]
