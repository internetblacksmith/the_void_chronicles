# Railway TCP Proxy Configuration Fix

## The Problem
You're getting "Connection closed" or "Connection reset" errors when trying to SSH.

## The Solution

### 1. Redeploy with New Configuration
```bash
railway up
```

### 2. Update TCP Proxy in Railway Dashboard

1. Go to your Railway service dashboard
2. Navigate to **Settings** → **Networking** → **TCP Proxy**
3. Make sure the **Target Port** is set to: `2222`
   - This is the port our SSH server listens on internally
4. Railway shows you:
   - **Proxy Domain**: yamabiko.proxy.rlwy.net
   - **Proxy Port**: 30460 (external port users connect to)

### 3. How It Works

```
User → yamabiko.proxy.rlwy.net:30460 → Railway TCP Proxy → Your App Port 2222
```

- Users connect to: `yamabiko.proxy.rlwy.net -p 30460`
- Railway forwards to: Your container's port `2222`
- SSH server listens on: Internal port `2222`

### 4. Connect

After redeploying and updating the TCP proxy target port to 2222:

```bash
ssh yamabiko.proxy.rlwy.net -p 30460
# Password: Amigos4Life!
```

## Why Port 2222?

- Port 2222 is a common alternative SSH port
- Avoids conflicts with system SSH (port 22)
- Clear separation from HTTP (port 80/8080)
- Standard in containerized environments

## Troubleshooting

If it still doesn't work:

1. **Check Railway logs** - Look for "SSH server listening on internal port 2222"
2. **Verify TCP Proxy Target Port** - Must be 2222 in Railway dashboard
3. **Wait for deployment** - Give it a minute after deploying
4. **Try verbose SSH**: `ssh -v yamabiko.proxy.rlwy.net -p 30460`

## Important Notes

- HTTP continues to work on the main domain
- SSH only works through the TCP proxy domain/port
- Both services run from the same deployment
- The 90s homepage remains your cover story! 🚧