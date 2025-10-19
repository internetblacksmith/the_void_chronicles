# Kamal Configuration Instructions for space_pirate

This document provides instructions for generating a complete Kamal deployment configuration for the space_pirate (Void Chronicles SSH Reader) application.

## Application Overview

**Name:** space_pirate (Void Chronicles SSH Reader)  
**Type:** Go application with SSH and HTTP servers  
**Purpose:** SSH-based book reader with 90s web disguise  
**Ports:**
- HTTP: 8080 (for web interface)
- SSH: 2222 (for book reader TUI)

**Tech Stack:**
- Go 1.23+
- Bubble Tea (TUI framework)
- Wish (SSH server)
- Built-in HTTP server
- Doppler (secret management)
- Sentry (error monitoring)

## Dependencies

### Required Services (Accessories)
**NONE** - This app is completely self-contained and does not require Redis or database.

### Shared Services
1. **Traefik** (reverse proxy)
   - Handles HTTP (port 8080) SSL termination
   - **Should already be running** from gcal-sinatra deployment
   - Do NOT redeploy Traefik

**Note:** SSH port (2222) is published directly to host, bypassing Traefik.

### External Services
1. **Doppler** (secrets management)
   - Service token required (starts with `dp.st.prd.`)
   - Minimal secrets needed (just SSH password)

2. **Sentry** (optional - error monitoring)
   - DSN required for error tracking
   - Go SDK integration

## Environment Variables (Managed via Doppler)

### Critical (Required for App to Function)
```bash
SSH_PASSWORD               # Password for SSH access (e.g., "Amigos4Life!")
HTTP_PORT                  # Should be "8080"
SSH_PORT                   # Should be "2222"
SSH_HOST                   # Should be "0.0.0.0"
KAMAL_REGISTRY_PASSWORD    # GitHub Personal Access Token (PAT) with write:packages scope
```

### Monitoring (Optional but Recommended)
```bash
SENTRY_DSN                 # Format: https://xxx@o123.ingest.sentry.io/456
SENTRY_ENVIRONMENT         # Should be "production"
```

### Doppler Integration
```bash
DOPPLER_TOKEN              # Service token from Doppler (set via kamal secrets)
```

## Kamal Configuration Structure

### config/deploy.yml

Create a `config/deploy.yml` file with the following structure:

```yaml
# Service name (used for container naming)
service: void-chronicles

# Docker image (GitHub Container Registry - free private images)
image: ghcr.io/GITHUB_USERNAME/void-chronicles

# Server configuration
servers:
  web:
    hosts:
      - SERVER_IP  # Replace with actual VPS IP address
    labels:
      # Traefik routing for HTTP interface
      traefik.http.routers.void-web.rule: Host(`vc.DOMAIN`)  # Replace DOMAIN
      traefik.http.routers.void-web.entrypoints: websecure
      traefik.http.routers.void-web.tls.certresolver: letsencrypt
      traefik.http.services.void-web.loadbalancer.server.port: 8080
    options:
      # IMPORTANT: Publish SSH port directly to host (bypasses Traefik)
      publish:
        - "2222:2222"
    # Override CMD to use Doppler for secret injection
    cmd: doppler run -- ./void-reader

# Docker registry authentication (GitHub Container Registry)
registry:
  server: ghcr.io
  username: GITHUB_USERNAME  # Replace with GitHub username (lowercase)
  password:
    - KAMAL_REGISTRY_PASSWORD  # GitHub Personal Access Token with write:packages scope

# Environment variables
env:
  # Only the Doppler token is passed to container
  # All other secrets are fetched by Doppler at runtime
  secret:
    - DOPPLER_TOKEN

# DO NOT INCLUDE accessories section (no Redis needed)
# DO NOT INCLUDE traefik section (already running from gcal-sinatra)

# Health check configuration
healthcheck:
  path: /
  port: 8080
  interval: 10s
  timeout: 5s
  max_attempts: 3

# Persistent volumes for user progress and SSH keys
volumes:
  - void-data:/app/.void_reader_data
  - void-ssh:/app/.ssh

# Deployment hooks
hooks:
  post-deploy: .kamal/hooks/post-deploy
```

### Dockerfile

The Dockerfile should use multi-stage build and include Doppler:

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder

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

# Install Doppler CLI
RUN curl -Ls --tlsv1.2 --proto "=https" --retry 3 \
    https://cli.doppler.com/install.sh | sh

# Copy binary from builder
COPY --from=builder /app/void-reader .

# Copy the book content
COPY book1_void_reavers_source ./book1_void_reavers_source

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

# Use Doppler to inject secrets at runtime
CMD ["doppler", "run", "--", "./void-reader"]
```

### .kamal/hooks/post-deploy

Create executable post-deployment hook:

```bash
#!/bin/bash
set -e

echo "✅ void-chronicles deployment completed"

# Optional: Notify Sentry of new release
if [ -n "$SENTRY_AUTH_TOKEN" ] && [ -n "$SENTRY_ORG" ] && [ -n "$SENTRY_PROJECT" ]; then
  VERSION=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
  
  curl https://sentry.io/api/0/organizations/$SENTRY_ORG/releases/ \
    -X POST \
    -H "Authorization: Bearer $SENTRY_AUTH_TOKEN" \
    -H 'Content-Type: application/json' \
    -d "{
      \"version\": \"$VERSION\",
      \"projects\": [\"$SENTRY_PROJECT\"]
    }" 2>/dev/null || echo "⚠️  Sentry notification failed (non-critical)"
fi

# Test SSH connectivity
if nc -z SERVER_IP 2222; then
  echo "✅ SSH port 2222 is accessible"
else
  echo "⚠️  SSH port 2222 is not accessible (check firewall)"
fi

echo "🎉 Post-deployment tasks completed"
```

Make it executable:
```bash
chmod +x .kamal/hooks/post-deploy
```

## Configuration Checklist

### Pre-Deployment Setup

- [ ] **Doppler Configuration**
  - [ ] Create Doppler project: `doppler projects create void-chronicles`
  - [ ] Setup environment: `doppler setup` (select void-chronicles, config: prd)
  - [ ] Add all required secrets (see Environment Variables section)
  - [ ] Create service token: `doppler configs tokens create kamal-deploy --project void-chronicles --config prd`
  - [ ] Save token securely

- [ ] **Sentry Setup** (optional)
  - [ ] Create Sentry project (Go platform)
  - [ ] Copy DSN to Doppler
  - [ ] Update Go code to initialize Sentry (see Code Integration section)

- [ ] **GitHub Container Registry**
  - [ ] Create GitHub Personal Access Token (Settings → Developer settings → Personal access tokens → Tokens (classic))
  - [ ] Select scope: `write:packages`, `read:packages`, `delete:packages`
  - [ ] Add token to Doppler as `KAMAL_REGISTRY_PASSWORD`
  - [ ] Images will be private by default (free on GitHub)

### Application Files

- [ ] **Create required files**
  - [ ] `config/deploy.yml` (see above)
  - [ ] Update `Dockerfile` (see above - ensure Doppler is installed)
  - [ ] `.kamal/hooks/post-deploy`

- [ ] **Update config/deploy.yml**
  - [ ] Replace `GITHUB_USERNAME` with your GitHub username (lowercase)
  - [ ] Replace `SERVER_IP` with VPS IP address
  - [ ] Replace `DOMAIN` with your domain (e.g., example.com)
  - [ ] Verify SSH port (2222) is published to host
  - [ ] **IMPORTANT:** Do NOT include `accessories` section
  - [ ] **IMPORTANT:** Do NOT include `traefik` section

- [ ] **Set Kamal secrets**
  - [ ] `kamal secrets set DOPPLER_TOKEN="dp.st.prd.YOUR_TOKEN"`

### DNS Configuration

- [ ] **Add DNS A record**
  - [ ] Create A record: `vc.DOMAIN` → `SERVER_IP`
  - [ ] Wait for propagation (check with `dig vc.DOMAIN`)

### Firewall Configuration

- [ ] **Ensure port 2222 is open**
  - [ ] SSH to server: `ssh deploy@SERVER_IP`
  - [ ] Check UFW: `sudo ufw status`
  - [ ] If port 2222 is not allowed: `sudo ufw allow 2222/tcp`
  - [ ] Verify: `sudo ufw status | grep 2222`

## Code Integration

### Sentry Integration (Optional)

Update `ssh-reader/main.go` to initialize Sentry:

```go
package main

import (
    "log"
    "os"
    "time"
    
    "github.com/getsentry/sentry-go"
    // ... other imports
)

func main() {
    // Initialize Sentry
    if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
        err := sentry.Init(sentry.ClientOptions{
            Dsn:              dsn,
            Environment:      getEnv("SENTRY_ENVIRONMENT", "production"),
            Release:          "void-chronicles@1.0.0", // Update with actual version
            TracesSampleRate: 0.1,
            BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
                // Add custom context
                event.Tags["component"] = "ssh-reader"
                return event
            },
        })
        if err != nil {
            log.Printf("Sentry initialization failed: %v", err)
        } else {
            log.Println("Sentry initialized successfully")
            defer sentry.Flush(2 * time.Second)
        }
    }
    
    // Capture panics
    defer func() {
        if r := recover(); r != nil {
            sentry.CurrentHub().Recover(r)
            sentry.Flush(2 * time.Second)
            panic(r)
        }
    }()
    
    // Rest of your application code...
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
```

Install Sentry Go SDK:
```bash
cd ssh-reader
go get github.com/getsentry/sentry-go
go mod tidy
```

## Deployment Order (Important!)

This is the **THIRD** app deployment. Traefik should already be running.

### Verify Prerequisites

```bash
# SSH to server
ssh deploy@SERVER_IP

# Check Traefik is running
docker ps | grep traefik
# Should show: traefik container running

# Check port 2222 is open in firewall
sudo ufw status | grep 2222
# Should show: 2222/tcp ALLOW Anywhere
```

## Deployment Commands

### Initial Deployment

```bash
# Navigate to project root (not ssh-reader subdirectory)
cd /path/to/space_pirate

# 1. Build Docker image for GitHub Container Registry
docker build -f Dockerfile -t ghcr.io/GITHUB_USERNAME/void-chronicles .

# 2. Login to GitHub Container Registry
echo $KAMAL_REGISTRY_PASSWORD | docker login ghcr.io -u GITHUB_USERNAME --password-stdin

# 3. Push to GitHub Container Registry
docker push ghcr.io/GITHUB_USERNAME/void-chronicles

# 3. Initialize Kamal (if not already done)
kamal init

# 4. Validate configuration
kamal config

# 5. Deploy application
kamal deploy

# 6. Check logs
kamal app logs -f

# 7. Verify HTTP deployment
curl -I https://vc.DOMAIN

# 8. Verify SSH deployment
ssh -p 2222 SERVER_IP
# Password: (from SSH_PASSWORD in Doppler)
```

### Update Deployment

```bash
# 1. Make code changes
git commit -am "Update feature"

# 2. Build and push new image to GitHub Container Registry
docker build -f Dockerfile -t ghcr.io/GITHUB_USERNAME/void-chronicles .
docker push ghcr.io/GITHUB_USERNAME/void-chronicles

# 3. Deploy update (zero-downtime for HTTP, brief interruption for SSH)
kamal deploy

# 4. Check logs
kamal app logs
```

### Rollback

```bash
# Rollback to previous version
kamal app rollback

# Check current version
kamal app version
```

## Verification Steps

After deployment, verify:

### 1. Container Health
```bash
kamal app containers
# Should show: void-chronicles running

kamal app logs | grep -i error
# Should have no critical errors
```

### 2. Doppler Integration
```bash
kamal app exec --interactive sh
doppler secrets
# Should list all secrets (values hidden)
exit
```

### 3. HTTP Access
```bash
curl -I https://vc.DOMAIN
# Should return: HTTP/2 200

curl https://vc.DOMAIN
# Should return HTML for "Bob's Personal Homepage" (90s disguise)
```

### 4. SSL Certificate
```bash
echo | openssl s_client -connect vc.DOMAIN:443 -servername vc.DOMAIN 2>/dev/null | openssl x509 -noout -dates
# Should show valid dates
```

### 5. SSH Access
```bash
# Test SSH connection
ssh -p 2222 SERVER_IP
# Enter password from Doppler (SSH_PASSWORD)

# Should show:
# - Book selection menu
# - 10 books in the series
# - Ability to navigate with arrow keys

# Try reading:
# - Press Enter to select Book 1
# - Should show chapter list
# - Navigate with h/l for chapters
# - Press q to quit
```

### 6. SSH Port Published
```bash
# Verify port 2222 is listening on host
ssh deploy@SERVER_IP 'sudo netstat -tulpn | grep 2222'
# Should show: docker-proxy listening on 0.0.0.0:2222

# Test from external machine
nc -zv SERVER_IP 2222
# Should show: Connection to SERVER_IP 2222 port [tcp/*] succeeded!
```

### 7. User Progress Tracking
```bash
# SSH to app, read some chapters
ssh -p 2222 SERVER_IP

# Exit and reconnect
# Progress should be saved

# Check data persistence
kamal app exec 'ls -la .void_reader_data'
# Should show progress files
```

### 8. Sentry Integration (if enabled)
```bash
kamal app logs | grep -i sentry
# Should see: "Sentry initialized successfully"

# Trigger an error in the app (try invalid input)
# Check Sentry dashboard for error
```

## Troubleshooting

### Doppler Secrets Not Loading

```bash
# Check token is set
kamal app exec 'env | grep DOPPLER'

# Test Doppler manually
kamal app exec --interactive sh
doppler secrets --token=$DOPPLER_TOKEN
```

### SSH Port Not Accessible

```bash
# Check firewall on server
ssh deploy@SERVER_IP 'sudo ufw status | grep 2222'

# If not allowed, add rule
ssh deploy@SERVER_IP 'sudo ufw allow 2222/tcp'

# Check port is published in Docker
kamal app exec --interactive sh
netstat -tulpn | grep 2222
```

### Can't Connect via SSH

```bash
# Check SSH_PASSWORD is set
kamal app exec 'env | grep SSH_PASSWORD'

# Check SSH server is running
kamal app logs | grep -i "ssh server"

# Test with verbose SSH
ssh -v -p 2222 SERVER_IP
```

### HTTP Works but SSH Doesn't

```bash
# Verify port publishing in deploy.yml
# Should have:
#   options:
#     publish:
#       - "2222:2222"

# Redeploy if needed
kamal deploy
```

### Progress Not Saving

```bash
# Check volume is mounted
kamal app exec 'ls -la .void_reader_data'

# Check permissions
kamal app exec 'ls -la / | grep void_reader_data'

# Check volume exists on host
ssh deploy@SERVER_IP 'docker volume ls | grep void'
```

## SSH Usage Guide

### Basic Controls

Once connected via SSH:

```
Navigation:
- ↑/↓ or j/k  : Navigate menu/scroll
- ←/→ or h/l  : Previous/Next chapter
- Enter       : Select
- b           : Set bookmark
- q           : Back/Quit
- ?           : Show help
```

### Example Session

```bash
# Connect to book reader
ssh -p 2222 SERVER_IP
# Enter password: Amigos4Life!

# You'll see book list (10 books in series)
# Currently only Book 1 is available

# Select Book 1: Void Reavers
# Press Enter

# Chapter list appears
# Use arrow keys to navigate
# Press Enter to read a chapter

# Navigate chapters:
# Press l (lowercase L) for next chapter
# Press h for previous chapter

# Quit:
# Press q to go back
# Press q again to exit
```

## Performance Optimization

### Resource Usage

Go app is extremely lightweight:
- Memory: ~30-50 MB
- CPU: Minimal (only during active SSH sessions)
- Disk: ~100 MB (including book content)

### Concurrent SSH Sessions

The app supports multiple simultaneous SSH connections:
- Each user gets independent session
- Progress tracked per session/user
- No performance degradation with multiple users

## Security Notes

1. **SSH Password** - Change default password in Doppler
2. **SSH Key** - Generated on first run, persisted in volume
3. **No sensitive data** - Just book content and reading progress
4. **User isolation** - Each SSH session is isolated
5. **HTTPS enforced** - For web interface

## Maintenance

### Update Secrets

```bash
# Change SSH password
doppler secrets set SSH_PASSWORD="NewSecurePassword123!"
kamal app restart
```

### View Logs

```bash
# Application logs
kamal app logs -f

# Filter for SSH connections
kamal app logs | grep -i "ssh"

# Filter for errors
kamal app logs | grep -i error
```

### Backup User Progress

```bash
# SSH to server
ssh deploy@SERVER_IP

# Backup progress data
sudo tar -czf void-progress-backup-$(date +%Y%m%d).tar.gz \
  /var/lib/docker/volumes/void-data/_data

# Download backup
scp deploy@SERVER_IP:~/void-progress-backup-*.tar.gz ./
```

### Clear User Progress (if needed)

```bash
kamal app exec --interactive sh
rm -rf .void_reader_data/*
# Users will start fresh on next connection
```

## Environment-Specific Configurations

### Staging Environment

```bash
doppler setup --config stg
doppler secrets set SSH_PASSWORD="staging_password"
doppler secrets set SENTRY_DSN="staging_sentry_dsn"

# Deploy to staging
kamal deploy --destination staging
```

### Development Environment

For local testing:

```bash
# Install Doppler locally
brew install dopplerhq/cli/doppler

# Setup dev config
cd space_pirate
doppler setup --config dev
doppler secrets set SSH_PASSWORD="dev_password"
doppler secrets set HTTP_PORT="8080"
doppler secrets set SSH_PORT="2222"

# Run locally with Doppler
cd ssh-reader
doppler run -- go run .

# Connect locally
ssh -p 2222 localhost
```

## Additional Resources

- **Bubble Tea Docs**: https://github.com/charmbracelet/bubbletea
- **Wish (SSH) Docs**: https://github.com/charmbracelet/wish
- **Kamal Docs**: https://kamal-deploy.org
- **Doppler Docs**: https://docs.doppler.com
- **Sentry Go SDK**: https://docs.sentry.io/platforms/go/
- **Project Guide**: `/home/paolo/projects/kamal_config/DOPPLER_SENTRY_POSTHOG_INTEGRATION.md`

---

**Generated for:** space_pirate (Void Chronicles SSH Reader)  
**Last Updated:** 2025-01-10  
**Kamal Version:** 2.x  
**Go Version:** 1.23+  
**IMPORTANT:** This is the THIRD app deployment. Traefik must already be running. SSH port 2222 must be open in firewall.
