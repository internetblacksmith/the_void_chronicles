FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git openssh-keygen

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o void-reader

FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates openssh-keygen netcat-openbsd

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/void-reader .

# Copy book content (will be added via build context)
COPY book1_void_reavers_source ./book1_void_reavers_source

# Create necessary directories
RUN mkdir -p .ssh .void_reader_data

# Generate SSH host key
RUN ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N "" -C "void-reader-docker-host"

# Create non-root user
RUN addgroup -g 1001 -S voidreader && \
    adduser -u 1001 -S voidreader -G voidreader

# Change ownership of directories
RUN chown -R voidreader:voidreader /app

# Switch to non-root user
USER voidreader

# Expose ports
EXPOSE 8080 2222

# Health check for HTTP
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD nc -z localhost 8080 || exit 1

# Start the application
CMD ["./void-reader"]