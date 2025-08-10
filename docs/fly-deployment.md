# Deploying to Fly.io

This guide walks through deploying the Void Chronicles SSH Reader to Fly.io.

## Prerequisites

1. Install the Fly CLI:
```bash
curl -L https://fly.io/install.sh | sh
```

2. Create a Fly.io account:
```bash
fly auth signup
```

3. Login to Fly.io:
```bash
fly auth login
```

## Deployment via GitHub Actions (Recommended)

The project is configured for automatic deployment via GitHub Actions.

1. Fork or clone the repository
2. Get your Fly API token:
```bash
fly auth token
```
3. Add it as a GitHub secret named `FLY_API_TOKEN` in your repository settings
4. Push to main branch to trigger deployment

## Manual Deployment

If you need to deploy manually:

```bash
fly deploy -a the-void-chronicles
```

## Configuration

The app is configured in `fly.toml`:
- HTTP service on ports 80/443 (with HTTPS redirect)
- SSH service on port 22
- 256MB RAM, 1 shared CPU
- London (lhr) region by default

## Environment Variables

Set via Fly.io secrets:
```bash
fly secrets set SSH_PASSWORD="YourSecurePassword" -a the-void-chronicles
```

Current environment variables:
- `HTTP_PORT`: 8080 (internal)
- `SSH_PORT`: 2222 (internal)
- `SSH_HOST`: 0.0.0.0
- `SSH_PASSWORD`: Your chosen password (default: Amigos4Life!)

## Accessing Your App

Once deployed, your app will be available at:

- **Web Interface**: `https://the-void-chronicles.fly.dev`
- **SSH Access**: `ssh the-void-chronicles.fly.dev`

Connect via SSH:
```bash
ssh the-void-chronicles.fly.dev
# Enter password when prompted
```

## Custom Domain

To add a custom domain:

1. Add the certificate:
```bash
fly certs add yourdomain.com -a the-void-chronicles
```

2. Add DNS records as instructed (usually a CNAME to your-app.fly.dev)

3. Wait for DNS propagation and certificate issuance

Example:
```
Type: CNAME
Name: vc
Target: the-void-chronicles.fly.dev
```

## Managing the App

### View logs:
```bash
fly logs -a the-void-chronicles
```

### Check app status:
```bash
fly status -a the-void-chronicles
```

### List certificates:
```bash
fly certs list -a the-void-chronicles
```

### SSH into the container (for debugging):
```bash
fly ssh console -a the-void-chronicles
```

## Costs

Fly.io free tier includes:
- 3 shared-cpu-1x VMs (256MB RAM each)
- 3GB persistent storage
- 160GB outbound transfer

This is perfect for the Void Reader app!

## Troubleshooting

### SSH connection refused
- Ensure the SSH service is configured in fly.toml
- Check that port 22 is exposed externally
- Verify SSH_PASSWORD is set

### Book content not showing
- Ensure book1_void_reavers_source is included in Docker image
- Check .dockerignore doesn't exclude .md files in book directory
- Verify COPY statement in Dockerfile

### SSL/Certificate issues with custom domain
- If using Cloudflare, set SSL/TLS mode to "Full" (not "Full strict")
- Or disable Cloudflare proxy (grey cloud) for direct connection
- Wait for certificate to be issued (check with `fly certs list`)

## GitHub Actions Workflow

The project includes `.github/workflows/fly-deploy.yml` for automatic deployment:
- Triggers on push to main branch
- Uses Fly.io's official GitHub Action
- Requires FLY_API_TOKEN secret in repository settings