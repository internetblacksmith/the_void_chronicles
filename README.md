# Void Reavers - Space Pirate Novel Project

A complete 20-chapter science fiction novel with multiple reading formats and an SSH-based terminal reader.

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

## Quick Start

### 1. Read the Book via SSH
```bash
./run.sh                    # Start SSH server
ssh localhost -p 23234      # Connect (password: Amigos4Life!)
```

### 2. Edit the Book
```bash
# Edit in your favorite editor
vim book1_void_reavers_source/chapters/chapter-01.md
```

### 3. Generate Publishing Formats
```bash
# PDF for print/Amazon KDP
./markdown_to_kdp_pdf.rb book1_void_reavers_source void_reavers.pdf

# EPUB for e-readers
./markdown_to_epub.rb book1_void_reavers_source void_reavers.epub
```

## Key Features

- **Complete Novel**: 20 chapters of space pirate adventure
- **SSH Terminal Reader**: Beautiful TUI with progress tracking
- **Multiple Formats**: PDF, EPUB, and terminal-readable
- **British English**: Properly formatted for UK spelling
- **Docker Support**: Easy deployment anywhere

## Documentation

- `MARKDOWN_STYLE_GUIDE.md` - Formatting guidelines for editing
- `README_kdp_publishing.md` - Amazon self-publishing guide
- `DEPLOYMENT.md` - Complete deployment guide for SSH reader
- `void_chronicles_series_bible.md` - Series planning and world-building
- `CLAUDE.md` - AI assistant guidance for this project

## Deployment

### Quick Start (Local)
```bash
cd ssh-reader
./run.sh
# Connect: ssh localhost -p 23234
```

### Cloud Deployment Options
See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed guides on:
- **Free hosting** (Fly.io, Railway, Google Cloud Run)
- **VPS hosting** ($5-10/month on DigitalOcean, Linode, Vultr)
- **Self-hosting** (Raspberry Pi, home server with Cloudflare)
- **Container deployment** (Docker, Kubernetes)

Popular options:
- **Fly.io** - Best free tier, global deployment
- **DigitalOcean** - $6/month, reliable and simple
- **Railway** - Easy deployment with free tier

## Requirements

- **SSH Reader**: Go 1.21+
- **PDF Generation**: Pandoc and LaTeX (or use HTML fallback)
- **EPUB Generation**: Pandoc

## License

© 2024 Captain J. Starwind. All rights reserved.