# SSL Certificate Renewal Guide

## Overview

Void Chronicles uses Let's Encrypt SSL certificates for HTTPS. Certificates expire after 90 days and must be renewed.

## Certificate Location

- **Let's Encrypt**: `/etc/letsencrypt/live/vc.internetblacksmith.dev/`
- **Docker Volume**: `/var/lib/docker/volumes/void-ssl/_data/`
- **Container Mount**: `/data/ssl/` (inside container)

## Manual Renewal

1. Stop the container (frees port 443):
   ```bash
   docker stop $(docker ps -q --filter "name=void-chronicles-web")
   ```

2. Renew certificate:
   ```bash
   certbot renew --standalone
   ```

3. Copy to Docker volume:
   ```bash
   sudo cp /etc/letsencrypt/live/vc.internetblacksmith.dev/fullchain.pem \
     /var/lib/docker/volumes/void-ssl/_data/cert.pem
   sudo cp /etc/letsencrypt/live/vc.internetblacksmith.dev/privkey.pem \
     /var/lib/docker/volumes/void-ssl/_data/key.pem
   sudo chmod 644 /var/lib/docker/volumes/void-ssl/_data/*.pem
   ```

4. Start the container:
   ```bash
   docker start $(docker ps -aq --filter "name=void-chronicles-web" | head -1)
   ```

## Automated Renewal (Recommended)

Use the provided `renew-ssl-certs.sh` script:

```bash
# Test renewal process (dry-run)
certbot renew --dry-run

# Run renewal script
./renew-ssl-certs.sh
```

### Cron Setup

Add to VPS crontab to run monthly:

```bash
# On VPS: edit crontab
sudo crontab -e

# Add this line (runs on 1st of each month at 3 AM)
0 3 1 * * /root/renew-ssl-certs.sh >> /var/log/ssl-renewal.log 2>&1
```

Or deploy the script to VPS and schedule via Kamal hooks.

## Verification

Check certificate expiration:
```bash
# On VPS
openssl x509 -enddate -noout -in /var/lib/docker/volumes/void-ssl/_data/cert.pem

# From client
echo | openssl s_client -connect vc.internetblacksmith.dev:443 2>/dev/null | openssl x509 -noout -enddate
```

## Troubleshooting

**Port 443 already in use during renewal:**
- Certbot standalone needs port 443. Stop the container first.

**Permission denied errors:**
- Ensure key.pem is readable: `sudo chmod 644 /var/lib/docker/volumes/void-ssl/_data/key.pem`

**Certificate not updating in container:**
- Restart container: `docker restart $(docker ps -q --filter "name=void-chronicles-web")`

## Current Certificate

- **Domain**: vc.internetblacksmith.dev
- **Issued**: 2025-10-19
- **Expires**: 2026-01-17
- **Issuer**: Let's Encrypt
