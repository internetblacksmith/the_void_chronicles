FROM golang:1.24-alpine AS builder

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

# Copy the actual book content
COPY book1_void_reavers_source ./book1_void_reavers_source

# Create necessary directories (but don't generate SSH key - it will be stored in persistent volume)
RUN mkdir -p .ssh .void_reader_data

# Create non-root user
RUN addgroup -g 1001 -S voidreader && \
    adduser -u 1001 -S voidreader -G voidreader && \
    chown -R voidreader:voidreader /app

USER voidreader

EXPOSE 8080 2222

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD nc -z localhost 8080 || exit 1

CMD ["./void-reader"]