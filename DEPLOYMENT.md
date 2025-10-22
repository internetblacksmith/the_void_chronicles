# SSH Reader Deployment Guide

## 🚀 Deployment Options

### 1. Kamal Deployment (Recommended)

**See [KAMAL_CONFIG_INSTRUCTIONS.md](KAMAL_CONFIG_INSTRUCTIONS.md) for complete setup guide.**

#### Features
- **Zero-downtime deployments** with Kamal orchestration
- **Doppler secret management** for secure environment variables
- **Native HTTPS support** with Let's Encrypt certificates
- **Direct port mapping**: HTTP (80), HTTPS (443), SSH (22)
- **Persistent volumes** for SSH keys, SSL certificates, and user progress

#### Requirements
- VPS with Docker and Docker Compose
- Doppler account for secret management
- Domain name with DNS configured
- Let's Encrypt SSL certificates (or self-signed for testing)

#### Quick Setup
```bash
# 1. Configure config/deploy.yml with your VPS IP and domain
# 2. Setup Doppler secrets (SSH_PASSWORD, etc.)
# 3. Deploy
kamal deploy

# Connect via SSH (standard port 22 mapped to container port 2222)
ssh your-domain.com
# Password: (from Doppler SSH_PASSWORD)

# HTTPS available at https://your-domain.com
# HTTP available at http://your-domain.com
```

#### SSL Certificate Setup
See [docs/ssl-certificate-renewal.md](docs/ssl-certificate-renewal.md) for certificate generation and renewal:
```bash
# Generate Let's Encrypt certificate
certbot certonly --standalone -d your-domain.com

# Copy to Docker volume
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem \
  /var/lib/docker/volumes/void-ssl/_data/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem \
  /var/lib/docker/volumes/void-ssl/_data/key.pem
sudo chmod 644 /var/lib/docker/volumes/void-ssl/_data/*.pem

# Auto-renewal with cron
./renew-ssl-certs.sh
```

**Cost**: VPS pricing (~$6-12/month for DigitalOcean/Linode/Vultr)

### 2. Budget VPS Options ($5-10/month)

#### **DigitalOcean Droplet**
- **Cost**: $6/month (1GB RAM, 25GB SSD)
- **Setup**:
```bash
# After creating droplet, SSH in
ssh root@your_droplet_ip

# Clone and setup
git clone https://github.com/yourusername/void-reader
cd void-reader/ssh-reader
./build.sh

# Install as systemd service
sudo ./deploy.sh
```

#### **Linode**
- **Cost**: $5/month (1GB RAM, 25GB SSD)
- **Similar setup to DigitalOcean**

#### **Vultr**
- **Cost**: $6/month (1GB RAM, 25GB SSD)
- **Good global locations**

### 3. Platform-as-a-Service

#### **Heroku** (Limited free tier)
- **Cost**: Free tier available with limitations
- **Requires**: Dockerfile and heroku.yml
```yaml
# heroku.yml
build:
  docker:
    web: Dockerfile
run:
  web: ./void-reader
```

#### **Railway** (Recommended for simple deployments)
- **Free tier**: 500 hours/month (~20 days)
- **Pros**: Easy deployment, automatic HTTPS, good performance
- **Setup**:
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login and initialize
railway login
railway init

# Deploy
railway up

# Get your app URL
railway domain
```

#### **Render**
- **Cost**: Free tier with spin-down
- **Pros**: Easy Docker deployment
```yaml
# render.yaml
services:
  - type: web
    name: void-reader
    env: docker
    dockerfilePath: ./Dockerfile
    envVars:
      - key: PORT
        value: 23234
```

### 4. Self-Hosted Options

#### **Raspberry Pi**
Perfect for home network access:
```bash
# Install Go
wget https://go.dev/dl/go1.21.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-armv6l.tar.gz

# Build and run
cd ssh-reader
./build.sh
./deploy.sh
```

#### **Home Server with Cloudflare Tunnel**
Expose your home server safely:
```bash
# Install cloudflared
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64
chmod +x cloudflared-linux-amd64
sudo mv cloudflared-linux-amd64 /usr/local/bin/cloudflared

# Create tunnel
cloudflared tunnel create void-reader
cloudflared tunnel route dns void-reader reader.yourdomain.com

# Configure and run
cloudflared tunnel run void-reader
```

## 🐳 Docker Deployment

### Build and Run with Docker
```bash
# Build image
docker build -t void-reader .

# Run container
docker run -d \
  --name void-reader \
  -p 2222:2222 \
  -v $(pwd)/.ssh:/app/.ssh \
  -v $(pwd)/.void_reader_data:/app/.void_reader_data \
  void-reader
```

### Docker Compose
```yaml
version: '3.8'
services:
  void-reader:
    build: .
    ports:
      - "2222:2222"
    volumes:
      - ./.ssh:/app/.ssh
      - ./.void_reader_data:/app/.void_reader_data
    restart: unless-stopped
```

## 🔒 Security Considerations

### SSL/TLS Certificates

The SSH reader supports native HTTPS with graceful fallback:

**For Kamal Deployment** (see [docs/ssl-certificate-renewal.md](docs/ssl-certificate-renewal.md)):
```bash
# Let's Encrypt certificates (recommended)
certbot certonly --standalone -d yourdomain.com
cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem \
  /var/lib/docker/volumes/void-ssl/_data/cert.pem
cp /etc/letsencrypt/live/yourdomain.com/privkey.pem \
  /var/lib/docker/volumes/void-ssl/_data/key.pem
chmod 644 /var/lib/docker/volumes/void-ssl/_data/*.pem
```

**For Local Development**:
```bash
# Self-signed certificates
mkdir -p .ssl
openssl req -x509 -newkey rsa:4096 -nodes \
  -keyout .ssl/key.pem \
  -out .ssl/cert.pem \
  -days 365 -subj "/CN=localhost"

# Set environment variables
export TLS_CERT_PATH=.ssl/cert.pem
export TLS_KEY_PATH=.ssl/key.pem
export HTTPS_PORT=8443
```

**Note**: If certificates are not found, HTTPS server gracefully skips startup and only HTTP is available.

### Password Authentication
The SSH reader supports both password and public key authentication:
- Default: Password authentication enabled
- Default password: `Amigos4Life!`
- Customizable via `SSH_PASSWORD` environment variable (required for production)
- Can be disabled by setting `SSH_REQUIRE_PASSWORD=false` (allows any public key)

```bash
# With password authentication (default)
ssh your-domain.com
# Enter password when prompted: Amigos4Life!

# Disable password authentication (allow any public key)
export SSH_REQUIRE_PASSWORD=false
./void-reader

# Set custom password via environment
export SSH_PASSWORD="YourSecurePassword"
export SSH_REQUIRE_PASSWORD=true
./void-reader
```

### SSH Key Management
```bash
# Generate a strong host key
ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N "" -C "void-reader-host"

# Set proper permissions
chmod 700 .ssh
chmod 600 .ssh/id_ed25519
```

### Firewall Configuration
```bash
# UFW (Ubuntu/Debian) - Kamal deployment
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 22/tcp    # SSH (mapped to container 2222)
sudo ufw enable

# firewalld (CentOS/RHEL)
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --permanent --add-port=22/tcp
sudo firewall-cmd --reload
```

### Rate Limiting with fail2ban
```ini
# /etc/fail2ban/jail.local
[void-reader]
enabled = true
port = 2222
filter = void-reader
logpath = /var/log/void-reader.log
maxretry = 5
bantime = 3600
```

## 📊 Monitoring

### Health Check Endpoint
Add a simple HTTP health check:
```go
// In main.go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})
go http.ListenAndServe(":8080", nil)
```

### Uptime Monitoring Services
- **UptimeRobot** - Free for 50 monitors
- **Pingdom** - Free tier available
- **StatusCake** - Free tier available

## 🌍 CDN and Performance

### Cloudflare (Recommended)
1. Add your domain to Cloudflare
2. Create a subdomain for the SSH reader
3. Enable Cloudflare Spectrum for SSH traffic (paid feature)

### Alternative: Multiple Deployments
Deploy to multiple regions for better latency:
- US: Fly.io (us-west)
- EU: Fly.io (eu-central)
- Asia: Fly.io (asia-pacific)

## 📝 Environment Variables

Create `.env` file for configuration:
```bash
# .env
SSH_PORT=2222
SSH_HOST=0.0.0.0
SSH_PASSWORD=Amigos4Life!  # Change this for production!
BOOK_PATH=./book1_void_reavers_source
DATA_PATH=./.void_reader_data
LOG_LEVEL=info
```

## 🔄 Continuous Deployment

### GitHub Actions
```yaml
# .github/workflows/deploy.yml
name: Deploy to Fly.io
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

## 📱 Client Access

### Desktop/Laptop
```bash
ssh reader.yourdomain.com -p 2222
# Password: Amigos4Life!
```

### Mobile Clients
- **iOS**: Termius, Prompt 3, Blink Shell
- **Android**: JuiceSSH, Termux, ConnectBot

### Web-based Terminal (Optional)
Use `ttyd` to provide web access:
```bash
ttyd -p 8080 ssh localhost -p 2222
```

## 🎯 Quick Start Recommendations

1. **For Testing**: Use Fly.io (free, easy setup)
2. **For Production**: DigitalOcean Droplet ($6/month)
3. **For Home Use**: Raspberry Pi with Cloudflare Tunnel
4. **For Scale**: Railway or Google Cloud Run

## 📞 Support

- Check application logs: `journalctl -u void-reader -f`
- Test SSH connection: `ssh -v localhost -p 2222`
- Verify port is open: `netstat -tlnp | grep 2222`