# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git openssh-keygen curl

WORKDIR /app

# Copy go files from ssh-reader directory
COPY ssh-reader/go.mod ssh-reader/go.sum ./
RUN go mod download

COPY ssh-reader/*.go ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o void-reader

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates openssh-keygen netcat-openbsd curl bash gnupg

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/void-reader .

# Copy the book content and series metadata
COPY book1_void_reavers_source ./book1_void_reavers_source
COPY ssh-reader/series.json ./series.json

# Create necessary directories
RUN mkdir -p .ssh .void_reader_data

# Create non-root user
RUN addgroup -g 1001 -S voidreader && \
    adduser -u 1001 -S voidreader -G voidreader && \
    chown -R voidreader:voidreader /app

USER voidreader

# Expose both HTTP and SSH ports
EXPOSE 8080 2222

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD nc -z localhost 8080 || exit 1

# Run the application directly
# Secrets are injected as environment variables during deployment via Doppler (on local machine)
CMD ["./void-reader"]