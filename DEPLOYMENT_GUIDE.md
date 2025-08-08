# Deployment Guide for Void Chronicles SSH Reader

## ⚠️ Important: Platform Limitations

Most Platform-as-a-Service (PaaS) providers like Railway, Heroku, and Render **only support HTTP/HTTPS traffic** on their standard plans. They do not expose arbitrary TCP ports needed for SSH connections.

### Current Status on Railway

✅ **Working:**
- HTTP server on port 80 (serves the 90s homepage)
- Health endpoint for monitoring
- Application builds and deploys successfully

❌ **Not Working:**
- SSH connections (Railway doesn't expose port 23234)
- Direct terminal access to the book reader

## Deployment Options

### Option 1: VPS Deployment (Recommended)
Deploy to a Virtual Private Server where you control the networking:

**Providers that support SSH:**
- DigitalOcean Droplets
- AWS EC2
- Linode
- Vultr
- Hetzner Cloud

**Quick Deploy to DigitalOcean:**
```bash
# 1. Create a droplet (Ubuntu 22.04)
# 2. SSH into your droplet
ssh root@your-droplet-ip

# 3. Clone the repository
git clone https://github.com/yourusername/space_pirate.git
cd space_pirate

# 4. Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 5. Build and run
./build.sh
sudo ./deploy.sh

# 6. Open firewall ports
ufw allow 80/tcp   # HTTP
ufw allow 23234/tcp # SSH Reader
ufw enable
```

### Option 2: Docker with Port Forwarding
Use Docker on any server with proper port mapping:

```bash
docker-compose up -d
# Exposes both HTTP (80) and SSH (23234)
```

### Option 3: Tunneling Service (Development)
For testing, use ngrok or similar to expose the SSH port:

```bash
# Terminal 1: Run the app
./run.sh

# Terminal 2: Expose SSH port with ngrok
ngrok tcp 23234

# Users connect to the ngrok URL
ssh ngrok-url.tcp.ngrok.io -p PORT_FROM_NGROK
```

### Option 4: Local Network Deployment
Run on a Raspberry Pi or home server:

```bash
# On your local server
./build.sh
./run.sh

# Forward ports on your router:
# - External 80 -> Internal 8080 (HTTP)
# - External 23234 -> Internal 23234 (SSH)

# Use dynamic DNS for a stable hostname
```

## Railway Deployment (HTTP Only)

Railway is great for the HTTP frontend but cannot handle SSH connections. The current deployment:

1. **Serves the disguised 90s homepage** at your Railway URL
2. **Provides health monitoring** at `/health`
3. **Cannot accept SSH connections** due to platform limitations

To use Railway as a frontend while running SSH elsewhere:
1. Deploy HTTP to Railway
2. Run SSH server on a VPS
3. Update the homepage to show the VPS SSH connection details

## Security Considerations

1. **Change the default password** in production:
   ```bash
   export SSH_PASSWORD="your-secure-password"
   ```

2. **Use SSH keys** instead of passwords when possible

3. **Restrict access** by IP if running publicly:
   ```bash
   # iptables example
   iptables -A INPUT -p tcp --dport 23234 -s YOUR_IP/32 -j ACCEPT
   iptables -A INPUT -p tcp --dport 23234 -j DROP
   ```

4. **Monitor access logs** regularly:
   ```bash
   journalctl -u void-reader -f
   ```

## Recommended Production Setup

For a production deployment:

1. **Frontend**: Railway or Vercel (serves the homepage)
2. **SSH Server**: DigitalOcean Droplet ($6/month)
3. **Monitoring**: UptimeRobot for both services
4. **DNS**: Cloudflare for DDoS protection

This separates concerns and uses each platform's strengths.

## Testing Your Deployment

After deployment, test both services:

```bash
# Test HTTP
curl http://your-domain.com/health

# Test SSH
ssh your-domain.com -p 23234
# Password: Amigos4Life! (or your custom password)
```

## Troubleshooting

### SSH Connection Hangs
- Platform doesn't support SSH (Railway, Heroku)
- Firewall blocking port 23234
- SSH server not running

### "Connection Refused"
- Check if both servers are running
- Verify ports are open
- Check logs: `journalctl -u void-reader`

### Build Failures
- Ensure Go 1.21+ is installed
- Run `go mod download` in ssh-reader directory
- Check for missing dependencies

## Support

For deployment help, please check:
- The logs in your deployment platform
- The systemd journal if using Linux
- Open an issue on GitHub with deployment logs