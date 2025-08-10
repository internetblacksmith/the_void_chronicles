# Alternative Deployment Options

## Fly.io Deployment (Recommended - Currently Used)

Fly.io handles multiple ports naturally and has a generous free tier. **This is the current deployment method.**

### Setup

1. Install Fly CLI:
```bash
curl -L https://fly.io/install.sh | sh
```

2. Sign up and login:
```bash
fly auth signup
fly auth login
```

3. Deploy via GitHub Actions (recommended):
- Add `FLY_API_TOKEN` to GitHub secrets
- Push to main branch

4. Access your app:
- HTTP: `https://the-void-chronicles.fly.dev` or custom domain
- SSH: `ssh the-void-chronicles.fly.dev` (port 22)

### Advantages
- No port conflicts
- Native SSH support on port 22
- Automatic HTTPS
- Global deployment options
- Persistent storage included

## Oracle Cloud Free Tier

Best for full control with always-free resources.

### Setup
1. Sign up at https://cloud.oracle.com/free
2. Create an Ubuntu instance (ARM or x86)
3. Deploy with Docker or systemd
4. Configure firewall for ports 80 and 22

### Advantages
- Completely free forever (not trial)
- Full root access
- No restrictions on ports
- 24GB RAM on ARM instances

## Render.com

Good for HTTP, SSH requires paid tier.

### For HTTP-only version:
```yaml
# render.yaml
services:
  - type: web
    name: void-reader
    env: docker
    plan: free
    envVars:
      - key: PORT
        value: 8080
```

## Comparison

| Provider | Free Tier | SSH Support | Setup Ease | Limits |
|----------|-----------|-------------|------------|--------|
| Fly.io | 3 VMs | ✅ Native | Easy | 256MB RAM |
| Oracle | Always Free | ✅ Full | Medium | None |
| Railway | $5 credit | ⚠️ Complex | Easy | Credit expires |
| Render | 750 hrs/mo | ❌ Paid only | Easy | HTTP only |
| Heroku | 550 hrs/mo | ❌ No | Easy | Sleeps |

## Migration Commands

### Export from Railway:
```bash
railway variables export > .env.production
```

### Deploy to Fly.io:
```bash
fly launch
fly secrets import < .env.production
fly deploy
```