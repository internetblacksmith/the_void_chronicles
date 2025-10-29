# 📚 The Void Chronicles: Space Pirates SSH Reader

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

An innovative SSH-based book reader that lets you experience "The Void Chronicles" science fiction series through your terminal. Features a beautiful TUI interface, progress tracking, and a clever 90s-style web disguise.

**🚀 Try it now:** 
- SSH: `ssh vc.internetblacksmith.dev` (Password: `Amigos4Life!`)
- HTTPS: https://vc.internetblacksmith.dev
- HTTP: http://vc.internetblacksmith.dev

## 🌟 Features

### 📖 Interactive SSH Book Reader
- **Beautiful TUI**: Split-view interface showing all 10 planned books in the series
- **Progress Tracking**: Automatic bookmarking and reading statistics per user
- **Keyboard Navigation**: Intuitive controls (h/l for chapters, arrows for scrolling)
- **Multi-User Support**: Individual progress tracking for each SSH connection

### 🕹️ 90s Web Disguise
- Serves a convincing "Bob's Personal Homepage" from 1998
- Perfect for discrete deployment on public servers
- Complete with marquee tags, "Under Construction" notices, and broken links

### 🚀 Modern Deployment
- **Kamal Deployment**: Zero-downtime deployments with Doppler secret management
- **Docker Support**: Full containerization with docker-compose
- **Systemd Integration**: Production-ready Linux service configuration

## Project Structure

```
.
├── book1_void_reavers_source/    # 📝 Markdown source files (edit these!)
│   ├── book.md                   # Main book file
│   ├── metadata.yaml             # Book metadata
│   └── chapters/                 # Individual chapter files
│
├── ssh-reader/                   # 💻 Terminal book reader
│   ├── main.go                   # SSH server and TUI
│   ├── book.go                   # Book loading logic
│   ├── progress.go               # Reading progress tracking
│   ├── void-reader               # Compiled binary
│   ├── Dockerfile                # Docker configuration
│   └── systemd/                  # Service files
│
└── Publishing Tools              # 🖨️ Format converters
    ├── markdown_to_kdp_pdf.rb    # → PDF for Amazon KDP
    ├── markdown_to_epub.rb       # → EPUB for e-readers
    └── latex_to_markdown_source.rb # LaTeX → Markdown converter

```

## ⚙️ Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and customize:

```bash
# HTTP Server
HTTP_PORT=8080         # HTTP server port

# HTTPS Server
HTTPS_PORT=8443        # HTTPS server port
TLS_CERT_PATH=/data/ssl/cert.pem  # Path to TLS certificate
TLS_KEY_PATH=/data/ssl/key.pem    # Path to TLS private key

# SSH Server  
SSH_PORT=2222         # SSH server port (internal container port)
SSH_HOST=0.0.0.0      # Bind address
SSH_PASSWORD=Amigos4Life!  # Authentication password (if password auth enabled)
SSH_REQUIRE_PASSWORD=true  # Set to "false" to disable password auth (allows any public key)

# Monitoring (Optional)
# SENTRY_DSN=https://your-sentry-dsn@sentry.io/project-id  # Sentry error tracking
# ENVIRONMENT=production                                    # Environment name
# RELEASE=void-reader@1.0.0                                 # Release version
# POSTHOG_API_KEY=phc_your_posthog_api_key                  # PostHog analytics
# POSTHOG_HOST=https://app.posthog.com                      # PostHog host URL
```

For Kamal deployment, secrets are managed via Doppler (see KAMAL_CONFIG_INSTRUCTIONS.md).

**Optional Monitoring**: The application supports Sentry (error tracking) and PostHog (analytics) integration. If environment variables are not set, the application runs normally without monitoring.

## 🎯 Quick Start

### Local Development

```bash
# Clone the repository
git clone https://github.com/yourusername/the-void-chronicles.git
cd the-void-chronicles

# Set up environment variables
cp .env.example .env
# Edit .env to customize ports and password

# Build and run
./build.sh
./run.sh

# In another terminal, connect via SSH (internal port 2222)
ssh localhost -p 2222
# Password: Amigos4Life! (or your custom password from .env)

# View the 90s homepage disguise
open http://localhost:8080

# HTTPS available with SSL certificates in .ssl/ directory
# open https://localhost:8443
```

### 🚢 Deploy with Kamal

```bash
# See KAMAL_CONFIG_INSTRUCTIONS.md for complete setup guide

# Quick deployment:
# 1. Configure config/deploy.yml with your VPS details
# 2. Setup Doppler and authenticate: doppler login
# 3. Setup project: doppler setup (select void-reader, config: prd)
# 4. Set Kamal Doppler token: kamal secrets set DOPPLER_TOKEN="dp.st.prd.YOUR_TOKEN"
# 5. Deploy: kamal deploy

# Connect via SSH (standard port 22 mapped to container port 2222)
ssh your-domain.com
# Password: (from Doppler SSH_PASSWORD)

# HTTPS available at https://your-domain.com (Let's Encrypt certificates)
# HTTP available at http://your-domain.com
```

**Setting up on another PC**: See [KAMAL_CONFIG_INSTRUCTIONS.md](KAMAL_CONFIG_INSTRUCTIONS.md#setting-up-development-on-another-pc) for complete instructions on deploying from a new machine using Doppler.

## 📖 The Void Chronicles Series

### Book 1: Void Reavers (Available Now!)
*A Tale of Space Pirates and Cosmic Plunder*

Captain Zara "Bloodhawk" Vega leads her crew through the lawless void between solar systems. When humanity attracts the attention of ancient alien Architects, pirates become unlikely diplomats in a test that will determine if humans deserve a place among the stars.

**20 chapters** of space opera adventure featuring:
- Epic space battles and heists
- First contact with alien civilizations
- Pirates as humanity's unlikely saviors
- A 50-year saga of transformation

### Upcoming Books (Planned)
- Book 2: Shadow Dancers - The rise of void-born humans
- Book 3: The Quantum Academy - Training the next generation
- Book 4: Empire of Stars - Consolidating human space
- Book 5-10: The complete saga through 2127

## 🎮 Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate menu/scroll text |
| `←/→` or `h/l` | Previous/Next chapter |
| `Enter` | Select |
| `b` | Set bookmark |
| `q` | Back/Quit |
| `?` | Show help |

## 📚 Documentation

- [Kamal Deployment Guide](KAMAL_CONFIG_INSTRUCTIONS.md) - Complete Kamal/Doppler setup
- [SSL Certificate Renewal](docs/ssl-certificate-renewal.md) - HTTPS certificate management
- [Deployment Options](DEPLOYMENT.md) - Alternative deployment methods
- [Style Guide](MARKDOWN_STYLE_GUIDE.md) - Markdown formatting conventions
- [Series Bible](void_chronicles_series_bible.md) - Complete series planning
- [Contributing](CONTRIBUTING.md) - How to contribute

## 🛠️ Installation

### Prerequisites
- Go 1.21 or higher
- Git
- SSH client

### From Source

```bash
# Install dependencies
cd ssh-reader
go mod download

# Build
go build -o void-reader .

# Run
./void-reader
```

### Using Docker

```bash
cd ssh-reader
docker-compose up -d
```

### Deploy to Production

See deployment documentation for detailed guides:
- **Kamal** (Recommended) - Zero-downtime VPS deployment with Doppler secrets (see KAMAL_CONFIG_INSTRUCTIONS.md)
- **Docker** - Containerized deployment with docker-compose
- **Systemd** - Production Linux service (see deploy.sh)
- **Alternative platforms** - See DEPLOYMENT.md for more options

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Workflow

```bash
# Fork and clone
git clone https://github.com/yourusername/void-chronicles.git

# Create feature branch
git checkout -b feature/amazing-feature

# Make changes and test
./run.sh

# Commit with conventional commits
git commit -m "feat: add amazing feature"

# Push and create PR
git push origin feature/amazing-feature
```

### How to Help Development

The project uses AI coding assistance. To help the AI agent work effectively:

**Build & Test Commands:**
- Test: `cd ssh-reader && go test ./...` or `make test`
- Single test: `cd ssh-reader && go test -run TestName`
- Coverage: `make test-coverage`
- Build: `cd ssh-reader && go build` or `make build`
- Lint: `cd ssh-reader && go fmt ./... && go vet ./...` or `make lint`
- Security scan: `make security-scan` (runs gosec)
- Pre-commit: `make pre-commit` (runs lint + tests + security scan)
- Deploy: `make deploy` (runs pre-commit checks first, then deploys)
- Local dev: `./run.sh` (HTTP:8080, HTTPS:8443, SSH:2222, password: Amigos4Life!)

**Deployment Safety**: The `make deploy` command automatically runs `make pre-commit` first, which includes:
- ✅ **Linting** (`go fmt`, `go vet`, `go mod tidy`) - Code formatting and vet checks
- ✅ **Tests** (all Go tests) - Ensures code works as expected
- ✅ **Security scan** (`gosec`) - Checks for security vulnerabilities

If any check fails, deployment is automatically cancelled. This prevents broken or vulnerable code from reaching production.

**Code Style:**
- All Go files start with AGPL-3.0 copyright header
- Imports: Standard library first, blank line, then external packages
- Naming: camelCase private, PascalCase exported
- Errors: Always wrap with context: `fmt.Errorf("failed to load book: %w", err)`
- Paths: Use `filepath.Join()` for cross-platform compatibility
- Comments: Document exported functions only

**Key Architecture:**
- Triple servers: HTTP (8080), HTTPS (8443), SSH (2222) in `main.go`
- TUI states: Main menu, chapter list, reading view, progress, about
- Progress tracking: JSON persistence in `.void_reader_data/username.json`
- Book loading: Markdown parser in `book.go` reads from `chapters/*.md`
- Deployment: Kamal with direct port mapping, Doppler secrets, persistent volumes

**Critical Policy:**
- Documentation-First: Before ANY commit, verify ALL documentation matches code
- NEVER commit changes unless explicitly requested
- Always run lint and typecheck before committing

## 🔒 Security

- SSH authentication: Password (default) or public key (configurable via `SSH_REQUIRE_PASSWORD`)
- Optional monitoring: Sentry and PostHog (no data collection if not configured)
- User progress stored locally in `.void_reader_data/`
- Report security issues to: security@voidchronicles.space

## 📄 License

This project uses a dual-license structure:

### 📖 Book Content
The Void Chronicles book series is licensed under **Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International (CC BY-NC-SA 4.0)**.

**You are free to:**
- Share and redistribute the books
- Remix, transform, and build upon the material

**Under these terms:**
- **Attribution** - Give appropriate credit
- **NonCommercial** - No commercial use without permission
- **ShareAlike** - Derivatives must use the same license

See [LICENSE-CONTENT.md](LICENSE-CONTENT.md) for full details.

### 💻 SSH Reader Application
The SSH reader software is licensed under **GNU Affero General Public License v3.0 (AGPL-3.0)**.

See [LICENSE](LICENSE) for full details.

**Note:** Contributions to either the books or the software are welcome under their respective license terms

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- SSH server powered by [Wish](https://github.com/charmbracelet/wish)
- Styling with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Inspired by classic BBS systems and terminal adventures

## 📮 Contact

- **Author**: Paolo Fabbri
- **Project**: [github.com/yourusername/void-chronicles](https://github.com/yourusername/void-chronicles)
- **Issues**: [Bug reports and feature requests](https://github.com/yourusername/void-chronicles/issues)

---

<p align="center">
  Made with ❤️ and lots of ☕ by terminal enthusiasts
  <br>
  <i>Experience books the way hackers do - through SSH!</i>
</p>