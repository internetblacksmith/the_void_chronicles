# Railway TCP Proxy Setup for SSH Reader

## Good News! Railway Now Supports TCP Proxy (2024)

Railway added TCP proxy support which allows exposing non-HTTP services like our SSH server!

## Setup Instructions

### 1. Deploy to Railway First
```bash
railway up
```

### 2. Configure TCP Proxy in Railway Dashboard

1. Go to your Railway project dashboard
2. Click on your `void_reavers` service
3. Navigate to **Settings** → **Networking**
4. In the **Public Networking** section, look for **TCP Proxy**
5. Click **Enable TCP Proxy**
6. Enter port: `23234` (our SSH port)
7. Railway will generate:
   - A proxy domain (like `tcp.railway.internal`)
   - A proxy port (like `12345`)

### 3. Update Your Application

The SSH server should already be configured to listen on port 23234. No changes needed in the code!

### 4. Connect via SSH

Once TCP proxy is enabled, connect using Railway's generated domain and port:

```bash
ssh <railway-tcp-domain> -p <railway-tcp-port>
# Password: Amigos4Life!
```

Example:
```bash
ssh voidreavers.tcp.railway.app -p 12345
# Password: Amigos4Life!
```

## Important Notes

- Railway can expose **both HTTP and TCP** on the same service
- HTTP continues on port 80 (the 90s homepage)
- TCP proxy handles SSH on the assigned proxy port
- The proxy uses random load balancing if you have multiple replicas

## What This Means

✅ **HTTP Homepage**: Still accessible at `voidreavers-production.up.railway.app`
✅ **SSH Reader**: Now accessible via TCP proxy domain and port
✅ **Both services**: Running from the same Railway deployment!

## Troubleshooting

If SSH still doesn't work after enabling TCP proxy:

1. **Check the logs** in Railway dashboard
2. **Verify the port** - ensure it's set to 23234
3. **Try redeploying** after enabling TCP proxy
4. **Check firewall** - some corporate networks block non-standard ports

## Alternative: If TCP Proxy Doesn't Work

Some users report SSH servers specifically have issues with Railway's TCP proxy. If that's the case, you can:

1. Use Railway's native SSH: `railway run bash`
2. Deploy SSH server to a VPS
3. Use a WebSocket-based solution (more complex)

## Testing After Setup

```bash
# Test HTTP (should show 90s homepage)
curl https://voidreavers-production.up.railway.app

# Test SSH (with Railway's TCP proxy domain)
ssh <your-tcp-domain> -p <your-tcp-port>
```

The TCP proxy feature is relatively new (2024), so if you encounter issues, Railway's support team can help!