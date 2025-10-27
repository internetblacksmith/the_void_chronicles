# Deployment Configuration Updates

## Changes Made

### 1. SSH Configuration
**Before:**
```yaml
ssh:
  user: root
  port: 1447
  keys:
    - ~/.ssh/digitalocean_vps
```

**After:**
```yaml
ssh:
  config: true
```

**Why:** Uses `~/.ssh/config` for SSH settings, making it consistent with gcal-sinatra and easier to maintain.

### 2. Server Host
**Before:**
```yaml
servers:
  web:
    hosts:
      - 161.35.165.206
```

**After:**
```yaml
servers:
  web:
    hosts:
      - digitalocean
```

**Why:** Uses SSH config alias instead of hardcoded IP address.

### 3. Port Configuration (IMPORTANT)
**Before:**
```yaml
options:
  publish:
    - "80:8080"   # HTTP - CONFLICT with Traefik!
    - "443:8443"  # HTTPS - CONFLICT with Traefik!
    - "22:2222"   # SSH
proxy: false
```

**After:**
```yaml
options:
  network: "private"
  publish:
    - "22:2222"   # SSH only
proxy: false

labels:
  traefik.enable: true
  traefik.http.routers.void-web.rule: Host(`vc.internetblacksmith.dev`)
  traefik.http.routers.void-web.entrypoints: websecure
  traefik.http.routers.void-web.tls.certresolver: letsencrypt
  traefik.http.services.void-web.loadbalancer.server.port: 8080
```

**Why:** 
- **Fixes port conflict** - Both apps were trying to use ports 80/443
- **Shares Traefik** - Uses gcal-sinatra's Traefik instance for HTTP/HTTPS
- **Simplifies SSL** - Traefik handles Let's Encrypt certificates automatically
- **Maintains SSH** - SSH still accessible on port 22 (direct mapping)

## Architecture

### Before (Port Conflict):
```
VPS Ports:
  80, 443  ← gcal-sinatra (via Traefik)
  80, 443  ← void-chronicles (direct) ❌ CONFLICT!
  22       ← void-chronicles SSH
```

### After (Shared Traefik):
```
VPS Ports:
  80, 443  ← Traefik (handles both apps)
             ├─ standup.internetblacksmith.dev → gcal-sinatra
             └─ vc.internetblacksmith.dev → void-chronicles HTTP
  22       ← void-chronicles SSH (direct)

Docker Network (private):
  gcal-sinatra:4567 ← Traefik
  void-chronicles:8080 ← Traefik (HTTP interface)
  void-chronicles:2222 ← Exposed as port 22
```

## What This Means

### HTTP/HTTPS Access
- **Before**: `https://vc.internetblacksmith.dev:443` (app's own HTTPS server)
- **After**: `https://vc.internetblacksmith.dev` (via Traefik, automatic SSL)

### SSH Access
- **Before**: `ssh 161.35.165.206` (port 22)
- **After**: `ssh vc.internetblacksmith.dev` (port 22) - Same!

### App Code Changes
**NONE REQUIRED!** The app still runs HTTP on 8080 and SSH on 2222 inside the container. The difference is:
- HTTP 8080 → Now accessible via Traefik (not direct port mapping)
- HTTPS 8443 → No longer used (Traefik handles SSL)
- SSH 2222 → Still mapped to host port 22

## Migration Steps

### 1. Update ~/.ssh/config (if needed)
Ensure you have:
```
Host digitalocean
  Hostname 161.35.165.206
  User deploy
  Port 1447
  IdentityFile ~/.ssh/digitalocean_vps
```

### 2. Ensure Traefik is Running
```bash
ssh digitalocean
docker ps | grep traefik
```

If not running, the Ansible playbook will set it up:
```bash
cd /path/to/vps-config/ansible
make setup
```

### 3. Deploy Updated Configuration
```bash
cd /path/to/vps-config/the_void_chronicles

# First deployment with new config might require cleanup
make deploy-cleanup
make deploy
```

### 4. Verify Deployment
```bash
# Check HTTP (should redirect to HTTPS)
curl -I http://vc.internetblacksmith.dev

# Check HTTPS
curl -I https://vc.internetblacksmith.dev

# Check SSH
ssh vc.internetblacksmith.dev
# Password: Amigos4Life!
```

## Troubleshooting

### Port already in use
If you see "port 22 already in use":
```bash
ssh digitalocean
docker ps -a | grep void
docker stop void-chronicles-web
docker rm void-chronicles-web
# Then redeploy
```

### Traefik not routing
```bash
ssh digitalocean

# Check Traefik logs
docker logs traefik

# Check container is on private network
docker network inspect private

# Should see void-chronicles-web in the network
```

### SSH not accessible
```bash
ssh digitalocean

# Check port mapping
docker ps | grep void
# Should show: 0.0.0.0:22->2222/tcp

# Check firewall
sudo ufw status
# Should allow port 22
```

## Benefits of This Approach

1. ✅ **No port conflicts** - Apps share Traefik properly
2. ✅ **Automatic SSL** - Traefik handles Let's Encrypt for both apps
3. ✅ **Unified config** - Both apps use same deployment pattern
4. ✅ **Easier maintenance** - Single Traefik instance to manage
5. ✅ **Better security** - Apps on private network, only necessary ports exposed
6. ✅ **No code changes** - App continues to work as before

## Notes

- The app's built-in HTTPS server (port 8443) is no longer used in production
- Traefik provides HTTPS termination instead
- The app's HTTP server (port 8080) now only serves HTTP to Traefik on the private network
- SSH functionality remains unchanged (direct port mapping to port 22)
