# Deploying to Fly.io

This guide walks through deploying the Void Reader SSH application to Fly.io.

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

## Initial Deployment

1. From the project root directory, launch the app:
```bash
fly launch
```

When prompted:
- **App name**: Choose a unique name (e.g., `void-reader-yourname`)
- **Region**: Select the closest region to you
- **Database**: Select "No" (we don't need one)
- **Redis**: Select "No"
- **Deploy now**: Select "Yes"

2. Set the SSH password secret:
```bash
fly secrets set SSH_PASSWORD="YourSecurePassword"
```

## Accessing Your App

Once deployed, your app will be available at:

- **Web Interface**: `https://your-app-name.fly.dev`
- **SSH Access**: `ssh your-app-name.fly.dev` (port 22!)

Connect via SSH:
```bash
ssh your-app-name.fly.dev
# Enter password when prompted
```

## Updating the App

After making changes to the code:

```bash
fly deploy
```

## Managing the App

### View logs:
```bash
fly logs
```

### Check app status:
```bash
fly status
```

### Scale the app:
```bash
# Scale to 2 instances
fly scale count 2

# Scale back to 1
fly scale count 1
```

### SSH into the container (for debugging):
```bash
fly ssh console
```

## Configuration

The `fly.toml` file contains all the configuration:

- **HTTP Service**: Runs on port 8080 internally, exposed on 80/443
- **SSH Service**: Runs on port 2222 internally, exposed on port 22
- **Persistent Storage**: 1GB volume mounted at `/data`
- **Auto-scaling**: Scales to zero when not in use

## Environment Variables

Set additional environment variables:
```bash
fly secrets set KEY=value
```

View current secrets (names only):
```bash
fly secrets list
```

## Custom Domain

Add a custom domain:
```bash
fly certs add yourdomain.com
```

Then update your DNS:
- **A Record**: Point to the IP from `fly ips list`
- **AAAA Record**: Point to the IPv6 from `fly ips list`

## Monitoring

View metrics in the dashboard:
```bash
fly dashboard
```

## Costs

Fly.io free tier includes:
- 3 shared-cpu-1x VMs (256MB RAM each)
- 3GB persistent storage
- 160GB outbound transfer

This is perfect for the Void Reader app!

## Troubleshooting

### App won't start
Check the logs:
```bash
fly logs
```

### SSH connection refused
Ensure the SSH service is running:
```bash
fly status
```

### Can't connect to SSH
Check that port 22 is exposed in `fly.toml`:
```toml
[[services]]
  internal_port = 2222
  protocol = "tcp"
  
  [[services.ports]]
    port = 22
```

## Backup and Data

The app stores user progress in `/data` which is persisted across deployments.

To backup:
```bash
fly ssh console -C "tar czf - /data" > backup.tar.gz
```

To restore:
```bash
cat backup.tar.gz | fly ssh console -C "tar xzf - -C /"
```