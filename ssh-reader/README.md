# SSH Reader for Void Reavers

A beautiful terminal-based book reader accessible via SSH, built with Go and the Charm libraries.

## Components

- `main.go` - SSH server and main application entry point
- `book.go` - Book loading and parsing logic
- `progress.go` - User progress tracking and bookmarks
- `void-reader` - Compiled binary (after building)

## Building

From the project root directory:
```bash
./build.sh
```

## Running

From the project root directory:
```bash
./run.sh
```

Then connect with password authentication:
```bash
ssh localhost -p 23234
# Password: Amigos4Life!
```

## Environment Variables

The SSH reader supports the following environment variables:

- `PORT` - HTTP server port (Railway provides this automatically, default: uses HTTP_PORT)
- `HTTP_PORT` - HTTP server port for local development (default: 8080)
- `SSH_PORT` - SSH server port (default: 23234 for local, 2222 for Railway)
- `SSH_HOST` - SSH server bind address (default: 0.0.0.0)
- `SSH_PASSWORD` - SSH authentication password (default: Amigos4Life!)
- `RAILWAY_ENVIRONMENT` - Set by Railway to detect cloud deployment

Example with custom settings:
```bash
SSH_PASSWORD="YourCustomPassword" SSH_PORT=2222 ./run.sh
```

## Docker Support

The application includes Docker configuration for easy deployment:
- `Dockerfile` - Multi-stage build for minimal image
- `docker-compose.yml` - Complete service configuration

## Systemd Deployment

For production deployment as a system service:
- `systemd/void-reader.service` - Service definition

## Note

This SSH reader must be run from the parent directory to access:
- Book content in `../book1_void_reavers_source/`
- SSH keys in `../.ssh/`
- User data in `../.void_reader_data/`