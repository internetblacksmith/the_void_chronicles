# Void Chronicles Deployment Guide

Complete step-by-step guide for deploying and managing the Void Chronicles SSH Reader application using Kamal.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
- [First Deployment](#first-deployment)
- [Updating the Application](#updating-the-application)
- [Managing Secrets](#managing-secrets)
- [Troubleshooting](#troubleshooting)
- [Maintenance Tasks](#maintenance-tasks)

## Overview

**What You're Deploying:**
- SSH-based book reader (port 2222)
- HTTP web interface (port 8080, SSL via Traefik)
- Self-contained Go application (no database/Redis needed)
- Private Docker images on GitHub Container Registry (free)
- Secrets managed via Doppler

**Architecture:**
```
User → SSH (2222) → VPS → Docker Container → Void Reader App
User → HTTPS → Traefik → Docker Container → HTTP Interface
```

## Prerequisites

### Local Machine Requirements

1. **Docker** - For building images
   ```bash
   docker --version  # Should be 20.10+
   ```

2. **Kamal** - Deployment tool
   ```bash
   gem install kamal
   kamal version  # Should be 2.0+
   ```

3. **Doppler CLI** - Secret management
   ```bash
   doppler --version  # Should be 3.0+
   ```

4. **Git** - Version control
   ```bash
   git --version
   ```

### VPS Requirements

1. **Ubuntu/Debian server** with:
   - Docker installed
   - SSH access (port 22)
   - Sudo privileges for deploy user
   - Traefik already running (from previous deployments)

2. **Firewall rules:**
   ```bash
   # SSH to your VPS
   ssh deploy@YOUR_VPS_IP
   
   # Check firewall status
   sudo ufw status
   
   # Should show these ports open:
   # 22/tcp    - SSH
   # 80/tcp    - HTTP (Traefik)
   # 443/tcp   - HTTPS (Traefik)
   # 2222/tcp  - Void Reader SSH
   
   # If 2222 is not open, add it:
   sudo ufw allow 2222/tcp
   sudo ufw reload
   ```

3. **Verify Traefik is running:**
   ```bash
   docker ps | grep traefik
   # Should show traefik container running
   ```

### Domain Requirements

- Domain name with DNS access
- A record pointing to your VPS IP
- Example: `vc.yourdomain.com` → `YOUR_VPS_IP`

## Initial Setup

### Step 1: GitHub Personal Access Token

Create a token for GitHub Container Registry (private images, free):

1. Go to: https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Settings:
   - **Note:** `Kamal Deployment - Void Chronicles`
   - **Expiration:** No expiration (or 1 year)
   - **Scopes:** Check these boxes:
     - `write:packages`
     - `read:packages`
     - `delete:packages`
4. Click "Generate token"
5. **SAVE THE TOKEN** - You won't see it again!
   - Format: `ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`

### Step 2: Configure Doppler Secrets

The Doppler project `void-reader` with `prd` config already exists. Add the registry password:

```bash
# Navigate to project directory
cd /home/paolo/books/space_pirate

# Add GitHub token to Doppler
doppler secrets set KAMAL_REGISTRY_PASSWORD="ghp_YOUR_GITHUB_TOKEN_HERE" \
  --project void-reader --config prd

# Verify all secrets are set
doppler secrets --project void-reader --config prd

# Should show:
# - SSH_PASSWORD: Amigos4Life!
# - HTTP_PORT: 8080
# - SSH_PORT: 2222
# - SSH_HOST: 0.0.0.0
# - KAMAL_REGISTRY_PASSWORD: ghp_...
```

**Current Doppler Service Token:**
```
dp.st.prd.AuWbT9LZ8bkUyvnQzYDI5kYZZJyvbXSEHbqKBkWJDpp
```

### Step 3: Update config/deploy.yml

Edit `/home/paolo/books/space_pirate/config/deploy.yml` and replace placeholders:

```bash
# Open the file
nano config/deploy.yml

# Replace these values:
GITHUB_USERNAME  → Your GitHub username (lowercase, e.g., "paolobrasolin")
SERVER_IP        → Your VPS IP address (e.g., "123.45.67.89")
DOMAIN           → Your domain (e.g., "blacksmith.tech")
```

**Example:**
```yaml
image: ghcr.io/paolobrasolin/void-chronicles

servers:
  web:
    hosts:
      - 123.45.67.89
    labels:
      traefik.http.routers.void-web.rule: Host(`vc.blacksmith.tech`)

registry:
  server: ghcr.io
  username: paolobrasolin
```

Save and exit (Ctrl+X, then Y, then Enter).

### Step 4: Set Kamal Secrets

Kamal needs the Doppler token to fetch secrets at runtime:

```bash
# Set the Doppler token for Kamal
kamal secrets set DOPPLER_TOKEN="dp.st.prd.AuWbT9LZ8bkUyvnQzYDI5kYZZJyvbXSEHbqKBkWJDpp"

# Verify it was set
kamal secrets
# Should show DOPPLER_TOKEN (value hidden)
```

### Step 5: Configure DNS

1. Log in to your DNS provider (Cloudflare, Route53, etc.)
2. Add an A record:
   - **Name:** `vc` (or whatever subdomain you chose)
   - **Type:** A
   - **Value:** Your VPS IP address
   - **TTL:** 300 (5 minutes)
3. Wait for DNS propagation (test with: `dig vc.yourdomain.com`)

### Step 6: Verify Prerequisites Checklist

Before deploying, verify:

```bash
# 1. Doppler secrets configured
doppler secrets --project void-reader --config prd | grep -E "SSH_PASSWORD|KAMAL_REGISTRY_PASSWORD"

# 2. Kamal secrets configured
kamal secrets | grep DOPPLER_TOKEN

# 3. config/deploy.yml has no placeholders
grep -E "GITHUB_USERNAME|SERVER_IP|DOMAIN" config/deploy.yml
# Should return nothing if all placeholders are replaced

# 4. DNS resolves correctly
dig vc.yourdomain.com +short
# Should return your VPS IP

# 5. Firewall allows port 2222
ssh deploy@YOUR_VPS_IP 'sudo ufw status | grep 2222'
# Should show: 2222/tcp ALLOW Anywhere

# 6. Traefik is running
ssh deploy@YOUR_VPS_IP 'docker ps | grep traefik'
# Should show traefik container
```

## First Deployment

### Step 1: Build and Push Docker Image

```bash
# Navigate to project root
cd /home/paolo/books/space_pirate

# Login to GitHub Container Registry
echo "ghp_YOUR_GITHUB_TOKEN" | docker login ghcr.io -u GITHUB_USERNAME --password-stdin

# Build the Docker image
docker build -f Dockerfile -t ghcr.io/GITHUB_USERNAME/void-chronicles:latest .

# This will:
# - Compile the Go application
# - Install Doppler CLI
# - Copy book content
# - Create multi-stage optimized image

# Push to GitHub Container Registry
docker push ghcr.io/GITHUB_USERNAME/void-chronicles:latest

# Verify image was pushed
# Go to: https://github.com/GITHUB_USERNAME?tab=packages
# You should see "void-chronicles" package (private)
```

### Step 2: Validate Kamal Configuration

```bash
# Test configuration syntax
kamal config

# Should output your full config with no errors
# Verify:
# - Image: ghcr.io/GITHUB_USERNAME/void-chronicles
# - Server IP is correct
# - Traefik labels have correct domain
# - SSH port 2222 is published
# - DOPPLER_TOKEN is in env.secret
```

### Step 3: Deploy Application

```bash
# Deploy to production
kamal deploy

# This will:
# 1. Pull image from ghcr.io
# 2. Create volumes (void-data, void-ssh)
# 3. Start container with Doppler integration
# 4. Configure Traefik routing
# 5. Run health checks
# 6. Execute post-deploy hook

# Expected output:
# ...
# Acquiring the deploy lock
# Finished all in X seconds
# ✅ void-chronicles deployment completed
# ✅ SSH port 2222 is accessible
# 🎉 Post-deployment tasks completed
```

### Step 4: Verify Deployment

```bash
# 1. Check container is running
kamal app containers
# Should show: void-chronicles-web running

# 2. Check logs
kamal app logs --tail 50
# Look for:
# - "Doppler initialized successfully" (or similar)
# - "HTTP server listening on :8080"
# - "SSH server listening on :2222"
# - No errors

# 3. Test HTTP endpoint
curl -I https://vc.yourdomain.com
# Should return: HTTP/2 200

curl https://vc.yourdomain.com
# Should return HTML for "Bob's Personal Homepage"

# 4. Test SSL certificate
echo | openssl s_client -connect vc.yourdomain.com:443 -servername vc.yourdomain.com 2>/dev/null | openssl x509 -noout -dates
# Should show valid dates from Let's Encrypt

# 5. Test SSH access
ssh -p 2222 YOUR_VPS_IP
# Enter password: Amigos4Life!
# Should show book selection menu

# 6. Verify Doppler secrets are loaded
kamal app exec 'env | grep -E "SSH_PASSWORD|HTTP_PORT|SSH_PORT"'
# Should show all environment variables
```

### Step 5: Test Full User Flow

```bash
# Connect via SSH
ssh -p 2222 YOUR_VPS_IP

# You should see:
# ┌─────────────────────────────────────┐
# │     VOID CHRONICLES LIBRARY         │
# │                                     │
# │  1. Book 1: Void Reavers           │
# │  2. Book 2: Coming Soon...         │
# │  ...                                │
# └─────────────────────────────────────┘

# Navigate with arrow keys
# Press Enter on "Book 1: Void Reavers"

# Chapter list appears:
# Chapter 1 - The Title
# Chapter 2 - Another Title
# ...

# Press Enter to read a chapter
# Test navigation:
# - l (lowercase L) = next chapter
# - h = previous chapter
# - q = quit/back

# Exit and reconnect
# Your progress should be saved
```

## Updating the Application

### When to Update

Update when you:
- Add new chapters to the book
- Fix bugs in the Go code
- Update dependencies
- Change configuration

### Update Process

```bash
# 1. Make your changes
cd /home/paolo/books/space_pirate

# Example: Add a new chapter
nano book1_void_reavers_source/chapters/chapter-21.md
# Add content...
# Save and exit

# 2. Commit changes (optional but recommended)
git add .
git commit -m "Add chapter 21"

# 3. Rebuild Docker image
docker build -f Dockerfile -t ghcr.io/GITHUB_USERNAME/void-chronicles:latest .

# 4. Push new image
docker push ghcr.io/GITHUB_USERNAME/void-chronicles:latest

# 5. Deploy update
kamal deploy

# Kamal will:
# - Pull new image
# - Perform zero-downtime rolling update for HTTP
# - Brief interruption for SSH (users will be disconnected)
# - Preserve user progress data (volumes persist)

# 6. Verify update
kamal app logs --tail 20

# 7. Test SSH connection
ssh -p 2222 YOUR_VPS_IP
# Verify new chapter appears
```

### Quick Update Commands

```bash
# Full update in one go
docker build -f Dockerfile -t ghcr.io/GITHUB_USERNAME/void-chronicles:latest . && \
docker push ghcr.io/GITHUB_USERNAME/void-chronicles:latest && \
kamal deploy
```

### Rollback if Something Goes Wrong

```bash
# Rollback to previous version
kamal app rollback

# Check current version
kamal app version

# View rollback history
kamal app containers --version
```

## Managing Secrets

### View Current Secrets

```bash
# View Doppler secrets (production)
doppler secrets --project void-reader --config prd

# View Kamal secrets
kamal secrets
```

### Update SSH Password

```bash
# Change password in Doppler
doppler secrets set SSH_PASSWORD="NewSecurePassword123!" \
  --project void-reader --config prd

# Restart app to pick up new password
kamal app restart

# Test new password
ssh -p 2222 YOUR_VPS_IP
# Use new password
```

### Update Other Secrets

```bash
# Update any secret in Doppler
doppler secrets set SECRET_NAME="value" --project void-reader --config prd

# Restart app
kamal app restart
```

### Rotate GitHub Token

```bash
# 1. Create new GitHub token (see Step 1 in Initial Setup)

# 2. Update in Doppler
doppler secrets set KAMAL_REGISTRY_PASSWORD="ghp_NEW_TOKEN" \
  --project void-reader --config prd

# 3. Update in Kamal
kamal secrets set KAMAL_REGISTRY_PASSWORD="ghp_NEW_TOKEN"

# 4. Next deployment will use new token
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs for errors
kamal app logs --tail 100

# Common issues:
# - Doppler token invalid: Update with kamal secrets
# - Port conflict: Check if port 2222 is already in use
# - Missing environment variables: Check doppler secrets

# Restart container
kamal app restart
```

### SSH Connection Refused

```bash
# 1. Verify container is running
kamal app containers

# 2. Check if SSH server started
kamal app logs | grep -i "ssh"

# 3. Verify port 2222 is published
ssh deploy@YOUR_VPS_IP 'docker ps --format "table {{.Names}}\t{{.Ports}}"'
# Should show: 0.0.0.0:2222->2222/tcp

# 4. Check firewall
ssh deploy@YOUR_VPS_IP 'sudo ufw status | grep 2222'

# 5. Test from VPS directly
ssh deploy@YOUR_VPS_IP 'nc -zv localhost 2222'
# Should succeed
```

### HTTP 502 Bad Gateway

```bash
# 1. Check container health
kamal app containers

# 2. Check HTTP server started
kamal app logs | grep -i "http"

# 3. Test health endpoint directly
kamal app exec 'curl -I http://localhost:8080'

# 4. Check Traefik labels
kamal config | grep -A 10 labels

# 5. Restart app
kamal app restart
```

### Doppler Secrets Not Loading

```bash
# 1. Verify DOPPLER_TOKEN is set
kamal app exec 'env | grep DOPPLER_TOKEN'

# 2. Test Doppler CLI inside container
kamal app exec --interactive sh
doppler secrets
exit

# 3. Check Doppler token is valid
doppler me --token "YOUR_DOPPLER_TOKEN"

# 4. Recreate token if needed
doppler configs tokens create kamal-deploy-new \
  --project void-reader --config prd --max-age 0

# 5. Update Kamal secret
kamal secrets set DOPPLER_TOKEN="NEW_TOKEN"

# 6. Redeploy
kamal deploy
```

### User Progress Lost

```bash
# Check if volume exists
ssh deploy@YOUR_VPS_IP 'docker volume ls | grep void-data'

# Check volume is mounted
kamal app exec 'ls -la .void_reader_data'

# If empty, progress was lost (volumes were recreated)
# To prevent: Never run "kamal remove" unless intentional
```

### SSL Certificate Issues

```bash
# Traefik manages certificates automatically
# If certificate not issued:

# 1. Check Traefik logs
ssh deploy@YOUR_VPS_IP 'docker logs traefik 2>&1 | grep -i letsencrypt'

# 2. Verify DNS is correct
dig vc.yourdomain.com +short
# Should return VPS IP

# 3. Check Traefik labels
kamal config | grep -A 5 "traefik.http.routers.void-web"

# 4. Restart Traefik (if needed)
ssh deploy@YOUR_VPS_IP 'docker restart traefik'

# Wait 1-2 minutes for certificate to be issued
```

## Maintenance Tasks

### View Application Logs

```bash
# Live logs
kamal app logs --follow

# Last 100 lines
kamal app logs --tail 100

# Search for errors
kamal app logs | grep -i error

# Filter for SSH connections
kamal app logs | grep -i "ssh"
```

### Check Resource Usage

```bash
# Container stats
ssh deploy@YOUR_VPS_IP 'docker stats void-chronicles-web --no-stream'

# Disk usage
ssh deploy@YOUR_VPS_IP 'docker system df'

# Volume sizes
ssh deploy@YOUR_VPS_IP 'docker system df -v | grep void'
```

### Backup User Progress

```bash
# SSH to VPS
ssh deploy@YOUR_VPS_IP

# Create backup
sudo tar -czf ~/void-progress-$(date +%Y%m%d).tar.gz \
  /var/lib/docker/volumes/void-data/_data

# Download to local machine
exit
scp deploy@YOUR_VPS_IP:~/void-progress-*.tar.gz ./
```

### Restore User Progress

```bash
# Upload backup to VPS
scp void-progress-20251017.tar.gz deploy@YOUR_VPS_IP:~/

# SSH to VPS
ssh deploy@YOUR_VPS_IP

# Stop container
cd /path/to/deployment/directory
kamal app stop

# Restore data
sudo rm -rf /var/lib/docker/volumes/void-data/_data/*
sudo tar -xzf ~/void-progress-20251017.tar.gz \
  -C /var/lib/docker/volumes/void-data/_data --strip-components=6

# Start container
kamal app start
```

### Clean Up Old Images

```bash
# Remove unused Docker images
ssh deploy@YOUR_VPS_IP 'docker image prune -a --force'

# Clean build cache
ssh deploy@YOUR_VPS_IP 'docker builder prune --force'
```

### Update Dependencies

```bash
# Update Go dependencies
cd /home/paolo/books/space_pirate/ssh-reader
go get -u ./...
go mod tidy

# Rebuild and deploy
cd ..
docker build -f Dockerfile -t ghcr.io/GITHUB_USERNAME/void-chronicles:latest .
docker push ghcr.io/GITHUB_USERNAME/void-chronicles:latest
kamal deploy
```

### Monitor Active SSH Sessions

```bash
# View logs in real-time
kamal app logs --follow | grep -i "session"

# Count active connections
kamal app exec 'netstat -an | grep :2222 | grep ESTABLISHED | wc -l'
```

## Common Commands Reference

```bash
# Deployment
kamal deploy                    # Deploy latest version
kamal app rollback              # Rollback to previous version
kamal app restart               # Restart container
kamal app stop                  # Stop container
kamal app start                 # Start container
kamal remove                    # Remove everything (DANGEROUS)

# Information
kamal config                    # Show configuration
kamal app version               # Show current version
kamal app containers            # List containers
kamal app logs                  # View logs
kamal app logs --follow         # Live logs
kamal app exec 'command'        # Execute command in container
kamal app exec --interactive sh # Interactive shell

# Secrets
kamal secrets                   # List secrets
kamal secrets set KEY=value     # Set secret
kamal secrets remove KEY        # Remove secret

# Doppler
doppler secrets --project void-reader --config prd           # List secrets
doppler secrets set KEY=value --project void-reader --config prd  # Set secret
doppler run -- command          # Run command with secrets

# Testing
ssh -p 2222 YOUR_VPS_IP         # Connect via SSH
curl https://vc.yourdomain.com  # Test HTTP
kamal app exec 'curl http://localhost:8080'  # Test health
```

## Security Best Practices

1. **SSH Password**
   - Change default password from `Amigos4Life!`
   - Use strong password (16+ characters)
   - Consider implementing SSH key authentication instead

2. **GitHub Token**
   - Use minimal scopes (only `write:packages`)
   - Set expiration (1 year max)
   - Rotate regularly

3. **Doppler Token**
   - Never commit to Git
   - Store in secure location
   - Rotate if compromised

4. **VPS Access**
   - Use SSH keys for VPS login
   - Disable password authentication
   - Keep sudo access restricted

5. **Firewall**
   - Only open necessary ports
   - Use UFW or iptables
   - Monitor for suspicious activity

## Getting Help

If you encounter issues:

1. **Check logs first:**
   ```bash
   kamal app logs --tail 200
   ```

2. **Search this guide** for troubleshooting steps

3. **Check Kamal documentation:** https://kamal-deploy.org

4. **Check Doppler docs:** https://docs.doppler.com

5. **Review GitHub Actions** (if CI/CD is set up)

## Next Steps

After successful deployment:

1. ✅ Test all features thoroughly
2. ✅ Set up monitoring (optional: Sentry integration)
3. ✅ Configure automated backups
4. ✅ Change default SSH password
5. ✅ Add remaining books to the series
6. ✅ Share SSH access with readers!

---

**Deployment Status Checklist:**

After deployment, verify these items:

- [ ] Container running: `kamal app containers`
- [ ] HTTP accessible: `curl https://vc.yourdomain.com`
- [ ] SSH accessible: `ssh -p 2222 YOUR_VPS_IP`
- [ ] SSL certificate valid: `openssl s_client` test
- [ ] Logs show no errors: `kamal app logs`
- [ ] User progress persists after reconnect
- [ ] Doppler secrets loading correctly
- [ ] All chapters visible in menu
- [ ] Navigation works (h/l/q keys)
- [ ] Bookmark feature works (b key)

---

**Document Version:** 1.0  
**Last Updated:** 2025-10-17  
**Application:** Void Chronicles SSH Reader  
**Deployment Method:** Kamal 2.x + Docker + GitHub Container Registry + Doppler
