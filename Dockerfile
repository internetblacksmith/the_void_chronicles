FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git openssh-keygen

WORKDIR /app

# Copy go files from ssh-reader directory
COPY ssh-reader/go.mod ssh-reader/go.sum ./
RUN go mod download

COPY ssh-reader/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o void-reader

FROM alpine:latest

RUN apk --no-cache add ca-certificates openssh-keygen netcat-openbsd

WORKDIR /app

COPY --from=builder /app/void-reader .

# Create necessary directories and placeholder content
RUN mkdir -p .ssh .void_reader_data book1_void_reavers_source/chapters && \
    echo "# Welcome to The Void Chronicles\n\nBook content will be loaded here." > book1_void_reavers_source/chapters/chapter-01.md && \
    ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N "" -C "void-chronicles-host"

# Create non-root user
RUN addgroup -g 1001 -S voidreader && \
    adduser -u 1001 -S voidreader -G voidreader && \
    chown -R voidreader:voidreader /app

USER voidreader

EXPOSE 8080 2222

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD nc -z localhost 8080 || exit 1

CMD ["./void-reader"]