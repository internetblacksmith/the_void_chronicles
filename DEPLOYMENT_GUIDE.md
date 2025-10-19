# Deployment Guide for Void Chronicles SSH Reader

## ⚠️ Important: Platform Limitations

Most Platform-as-a-Service (PaaS) providers like Railway, Heroku, and Render **only support HTTP/HTTPS traffic** on their standard plans. They do not expose arbitrary TCP ports needed for SSH connections.

**For SSH functionality, you must use a VPS or container platform with direct port mapping.**

## Recommended Deployment Options

### Option 1: Kamal Deployment (Recommended for Production)

**Complete guide:** [KAMAL_CONFIG_INSTRUCTIONS.md](KAMAL_CONFIG_INSTRUCTIONS.md)

Deploy to a VPS using Kamal orchestration with Doppler secret management:

**Features:**
- ✅ Zero-downtime deployments
- ✅ Native HTTPS with Let's Encrypt
- ✅ Direct port mapping (HTTP:80, HTTPS:443, SSH:22)
- ✅ Persistent volumes for data and certificates
- ✅ Doppler secret management

**Providers:**
- DigitalOcean Droplets ($6/month)
- Linode ($5/month)
- Vultr ($6/month)
- Hetzner Cloud (€4/month)

**Quick Deploy:**
```bash
# 1. Configure config/deploy.yml with your VPS details
# 2. Setup Doppler secrets
kamal secrets set DOPPLER_TOKEN="dp.st.prd.YOUR_TOKEN"

# 3. Deploy
kamal deploy

# 4. Setup SSL certificates (see docs/ssl-certificate-renewal.md)
# On VPS:
certbot certonly --standalone -d your-domain.com
sudo cp /etc/letsencrypt/live/your-domain.com/*.pem \
  /var/lib/docker/volumes/void-ssl/_data/
sudo chmod 644 /var/lib/docker/volumes/void-ssl/_data/*.pem

# 5. Connect
ssh your-domain.com  # Port 22, password: Amigos4Life!
https://your-domain.com  # 90s homepage
```

### Option 2: Manual VPS Deployment

Deploy directly to a VPS without Kamal:

```bash
# 1. SSH into your VPS
ssh root@your-vps-ip

# 2. Install Go
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 3. Clone repository
git clone https://github.com/yourusername/space_pirate.git
cd space_pirate

# 4. Build and deploy
./build.sh
sudo ./deploy.sh

# 5. Configure firewall
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 2222/tcp
sudo ufw enable

# 6. Setup SSL (optional)
mkdir -p .ssl
certbot certonly --standalone -d your-domain.com
cp /etc/letsencrypt/live/your-domain.com/fullchain.pem .ssl/cert.pem
cp /etc/letsencrypt/live/your-domain.com/privkey.pem .ssl/key.pem

# 7. Restart service
sudo systemctl restart void-reader
```

### Option 3: Docker Compose

For local or VPS deployment with Docker:

```bash
# Clone repository
git clone https://github.com/yourusername/space_pirate.git
cd space_pirate/ssh-reader

# Create .env file
cp ../.env.example .env
# Edit .env with your settings

# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Connect
ssh localhost -p 2222  # SSH reader
curl http://localhost:8080  # 90s homepage
```

### Option 4: Local Network / Home Server

Run on Raspberry Pi or home server:

```bash
# Build and run
./build.sh
./run.sh

# Forward ports on router:
# - External 80 → Internal 8080 (HTTP)
# - External 2222 → Internal 2222 (SSH)

# Use dynamic DNS for stable hostname (DuckDNS, No-IP, etc.)
```

### Option 5: Development/Testing with ngrok

Expose local instance temporarily:

```bash
# Terminal 1: Start app
./run.sh

# Terminal 2: Expose SSH port
ngrok tcp 2222

# Terminal 3: Expose HTTP port
ngrok http 8080

# Connect using ngrok URLs
ssh tcp.ngrok.io -p XXXXX  # From ngrok output
```

## What Works Where

| Feature | Kamal (VPS) | Manual VPS | Docker Compose | PaaS (Railway/Heroku) |
|---------|-------------|------------|----------------|----------------------|
| HTTP Server | ✅ | ✅ | ✅ | ✅ |
| HTTPS Server | ✅ | ✅ | ✅ (manual cert) | ✅ (auto) |
| SSH Reader | ✅ | ✅ | ✅ | ❌ No TCP ports |
| Zero-downtime | ✅ | ❌ | ❌ | ✅ |
| Auto SSL | ✅ (manual) | ❌ | ❌ | ✅ |
| Cost | $5-12/mo | $5-12/mo | VPS/free | Free-$5/mo |

## Security Best Practices

### 1. Change Default Password
```bash
# In .env or Doppler
SSH_PASSWORD="YourVerySecurePasswordHere123!"
```

### 2. Use Strong SSL Certificates
```bash
# Let's Encrypt (recommended)
certbot certonly --standalone -d your-domain.com

# Or self-signed for testing
openssl req -x509 -newkey rsa:4096 -nodes \
  -keyout .ssl/key.pem -out .ssl/cert.pem \
  -days 365 -subj "/CN=your-domain.com"
```

### 3. Configure Firewall
```bash
# Allow only necessary ports
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp    # VPS SSH (if different from app)
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 2222/tcp  # SSH Reader (or 22 for Kamal)
sudo ufw enable
```

### 4. Monitor Logs
```bash
# Kamal deployment
kamal app logs --tail

# Systemd service
journalctl -u void-reader -f

# Docker
docker logs -f void-chronicles-web
```

### 5. Rate Limiting
The SSH reader includes built-in rate limiting:
- Max 3 failed login attempts
- 1-minute cooldown
- Session timeout after inactivity

## SSL Certificate Management

See [docs/ssl-certificate-renewal.md](docs/ssl-certificate-renewal.md) for complete guide.

**Quick renewal:**
```bash
# Automated renewal script
./renew-ssl-certs.sh

# Setup cron for monthly renewal
sudo crontab -e
# Add: 0 3 1 * * /path/to/renew-ssl-certs.sh >> /var/log/ssl-renewal.log 2>&1
```

## Testing Your Deployment

```bash
# Test HTTP
curl http://your-domain.com/health

# Test HTTPS
curl https://your-domain.com

# Test SSH (interactive)
ssh your-domain.com -p 22  # Kamal
# or
ssh your-domain.com -p 2222  # Manual/Docker

# Test SSH (automated)
sshpass -p "Amigos4Life!" ssh -o StrictHostKeyChecking=no \
  your-domain.com -p 22 exit
```

## Troubleshooting

### SSH Connection Hangs
- **PaaS platforms**: Most don't support SSH (use VPS instead)
- **Firewall**: Check ports are open
- **Container**: Verify container is running

### HTTPS Certificate Errors
- **Self-signed**: Browser warns (expected), use `-k` with curl
- **Let's Encrypt**: Check domain DNS and port 443 accessibility
- **Permissions**: Ensure certificate files are readable (644)

### "Connection Refused"
```bash
# Check if services are running
docker ps | grep void-chronicles  # Kamal
sudo systemctl status void-reader  # Systemd
docker-compose ps  # Docker Compose

# Check logs
kamal app logs --tail  # Kamal
journalctl -u void-reader -n 50  # Systemd
docker-compose logs  # Docker Compose

# Verify ports
netstat -tlnp | grep -E '(80|443|2222)'
```

### Build Failures
```bash
# Install dependencies
cd ssh-reader
go mod download

# Clean build
go clean -cache
go build -v

# Check Go version (need 1.21+)
go version
```

## Monitoring

### Health Checks
```bash
# HTTP health endpoint (returns "OK")
curl http://your-domain.com/health

# Check certificate expiry
echo | openssl s_client -connect your-domain.com:443 2>/dev/null | \
  openssl x509 -noout -enddate
```

### Uptime Monitoring
- **UptimeRobot** - Free for 50 monitors
- **Pingdom** - Free tier available
- **StatusCake** - Free tier available

Monitor both:
- HTTP/HTTPS endpoint: `https://your-domain.com/health`
- SSH port: TCP check on port 22 (or 2222)

## Recommended Production Setup

**Architecture:**
1. **VPS**: DigitalOcean Droplet ($6/month)
2. **Deployment**: Kamal with Doppler secrets
3. **SSL**: Let's Encrypt with auto-renewal
4. **Monitoring**: UptimeRobot for health checks
5. **DNS/CDN**: Cloudflare for DDoS protection

**Estimated monthly cost:** $6-12

## Support

For deployment help:
- Check logs first (see monitoring section)
- Review [KAMAL_CONFIG_INSTRUCTIONS.md](KAMAL_CONFIG_INSTRUCTIONS.md)
- Open GitHub issue with deployment platform and logs
