# 📚 The Void Chronicles: Space Pirates SSH Reader

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://go.dev/)
[![Fly.io Deploy](https://img.shields.io/badge/Deploy-Fly.io-783EF9?logo=fly.io)](https://fly.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

An innovative SSH-based book reader that lets you experience "The Void Chronicles" science fiction series through your terminal. Features a beautiful TUI interface, progress tracking, and a clever 90s-style web disguise.

**🚀 Try it now:** 
- SSH: `ssh vc.internetblacksmith.dev` (Password: `Amigos4Life!`)
- Web: https://vc.internetblacksmith.dev

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
- **Fly.io Deployed**: Automatic deployment via GitHub Actions
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

# SSH Server  
SSH_PORT=2222         # SSH server port
SSH_HOST=0.0.0.0      # Bind address
SSH_PASSWORD=Amigos4Life!  # Authentication password
```

For Fly.io deployment, use `fly secrets set` instead of a `.env` file.

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

# In another terminal, connect via SSH
ssh localhost -p 2222
# Password: Amigos4Life! (or your custom password from .env)

# View the 90s homepage disguise
open http://localhost:8080
```

### 🚢 Deploy to Fly.io

```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Sign up and authenticate
fly auth signup
fly auth login

# Launch the app (first time only)
fly launch

# Deploy updates
fly deploy

# Connect via SSH (standard port 22!)
ssh void-chronicles.fly.dev
# Password: Amigos4Life!
```

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

- [Deployment Guide](DEPLOYMENT_GUIDE.md) - Deploy to 12+ platforms
- [Fly.io Deployment](docs/fly-deployment.md) - Production deployment guide
- [Style Guide](MARKDOWN_STYLE_GUIDE.md) - Markdown formatting conventions
- [Series Bible](void_chronicles_series_bible.md) - Complete series planning
- [Contributing](CONTRIBUTING.md) - How to contribute
- [API Reference](docs/api-reference.md) - Technical documentation

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

See [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) for detailed guides on:
- **Fly.io** (Currently deployed) - Native SSH and HTTP support
- **DigitalOcean** - $6/month VPS with full control
- **Fly.io** - Global edge deployment
- **Self-hosting** - Raspberry Pi or home server

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

## 🔒 Security

- SSH connections use password authentication (configurable)
- No telemetry or data collection
- User progress stored locally in `.void_reader_data/`
- Report security issues to: security@voidchronicles.space

## 📄 License

This project is licensed under the **GNU Affero General Public License v3.0** - see [LICENSE](LICENSE) for details.

### Additional Terms
- The Void Chronicles narrative content is also available under **CC-BY-SA 4.0**
- Contributions are welcome under the same license terms

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