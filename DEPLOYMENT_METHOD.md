# Current Deployment Method

## ✅ Current: Kamal + Docker (Recommended)

The application is deployed using **Kamal** with Docker containers.

### Quick Deploy
```bash
make deploy
```

### What This Does
- Builds Docker image
- Pushes to GitHub Container Registry (ghcr.io)
- Deploys to VPS via Kamal
- Connects to Traefik for HTTPS
- Exposes SSH on port 22

### Requirements
- VPS configured via `vps-config/ansible` (infrastructure management)
- Doppler for secrets management
- Docker on VPS (installed by Ansible)
- Traefik reverse proxy (installed by Ansible)

### Configuration
- **Kamal config**: `config/deploy.yml`
- **Dockerfile**: Root `Dockerfile`
- **Secrets**: Managed via Doppler

### Infrastructure
Infrastructure is managed **separately** in the `vps-config` repository:
- VPS setup (Docker, Traefik, firewall)
- Deploy user creation
- SSH configuration
- System security

**Do not manage infrastructure from this repo!**

## 🗑️ Removed: SystemD Deployment (Old Method)

The following files were **removed** as they are no longer used:

- ❌ `deploy.sh` - Old systemd deployment script
- ❌ `renew-ssl-certs.sh` - Manual SSL renewal (Traefik handles this now)
- ❌ `ssh-reader/systemd/` - SystemD service files

These files are preserved in git history if needed:
```bash
git log --all --full-history -- deploy.sh
git show <commit>:deploy.sh
```

## Why the Change?

### Old Method (SystemD)
```
Problems:
- Manual VPS setup required
- Systemd service management
- Manual SSL certificate renewal
- No zero-downtime deployments
- Infrastructure mixed with app code
```

### Current Method (Kamal)
```
Benefits:
✅ Automated deployments
✅ Zero-downtime rolling updates
✅ Infrastructure as code (separate repo)
✅ Automatic SSL via Traefik
✅ Rollback capability
✅ Consistent across all apps
```

## Deployment Architecture

```
┌─────────────────────────────────────────┐
│  VPS (161.35.165.206)                  │
│                                         │
│  Port 80/443 → Traefik                 │
│                ├─ HTTPS: vc.domain.dev │
│                └─ Auto SSL (Let's Enc) │
│                                         │
│  Port 22 → void-chronicles SSH         │
│            (container port 2222)       │
│                                         │
│  Docker Network: private               │
│  └─ void-chronicles:8080 (HTTP)       │
│                                         │
│  Volumes:                              │
│  ├─ void-data (user progress)         │
│  └─ void-ssl (SSL certs, if needed)   │
└─────────────────────────────────────────┘
```

## See Also

- **Deployment guide**: `DEPLOYMENT_UPDATES.md` (migration notes)
- **Infrastructure**: `../ansible/README.md` (VPS setup)
- **Kamal config**: `config/deploy.yml`
- **Architecture**: `../DEPLOYMENT_ARCHITECTURE.md`
