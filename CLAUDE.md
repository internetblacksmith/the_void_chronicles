# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## CRITICAL: Commit and Documentation Policy

**📝 Documentation-First Commit Policy:**

You MAY commit changes autonomously, but ONLY after:
1. **Documentation is FULLY updated** to reflect all code changes
2. **README files** match the current implementation exactly
3. **Deployment guides** reflect current configuration
4. **File paths in docs** are verified and correct
5. **CLAUDE.md** is updated if architecture or commands changed

**The workflow MUST be:**
1. Make code changes
2. Update ALL relevant documentation
3. Verify docs match code state
4. Then commit with a clear message

Documentation drift is unacceptable - docs must ALWAYS match the code state before any commit.

## Project Overview

This repository contains **two main components**:

1. **"The Void Chronicles" Book Series** - Starting with Book 1: "Void Reavers: A Tale of Space Pirates and Cosmic Plunder" - a complete 20-chapter science fiction novel about Captain Zara "Bloodhawk" Vega's transformation from pirate to diplomat in a universe where humanity must prove itself to alien Architects.

2. **Void Reavers SSH Reader** - A Go application that creates an SSH server allowing users to read books through a beautiful terminal user interface (TUI) built with Bubbletea and Wish.

## Architecture & File Structure

### Book Content Structure
The series is organized with each book in its own directory:
- `book1_void_reavers_source/` - Contains the complete Markdown source
- `void_chronicles_series_bible.md` - Master document with series overview and future book plans
- Book conversion scripts in root directory for processing any book

Each book's source is written in Markdown format:

1. **Source**: Markdown files in `book1_void_reavers_source/chapters/`
2. **Output**: PDF (via Pandoc) and EPUB for publishing
3. **Reading**: SSH reader loads Markdown directly

The Markdown source is the canonical version. All other formats are generated from it.

### SSH Reader Application Structure
The application is located in the `ssh-reader/` directory and consists of:
- `main.go` (~950 lines) - Dual HTTP/SSH servers, split-view TUI, 10-book series menu
- `book.go` (~250 lines) - Book loading from Markdown source
- `progress.go` (~150 lines) - User progress tracking and bookmarks
- `go.mod/go.sum` - Dependencies (Bubbletea, Lipgloss, Wish, Gorilla WebSocket)
- Railway deployment configs (`railway.toml`, `nixpacks.toml`, `Procfile`)
- Docker configuration (`docker-compose.yml`)
- Build and deployment scripts in root directory

## Essential Commands

### SSH Reader Application Commands

#### Quick Setup and Run
```bash
# Build the SSH reader application
./build.sh

# Start both HTTP (8080) and SSH (23234) servers
./run.sh

# View the 90s homepage
open http://localhost:8080

# Connect to read the book
ssh localhost -p 23234
# Password: Amigos4Life!
```

#### Production Deployment
```bash
# Deploy as systemd service (requires sudo)
sudo ./deploy.sh

# Manage the service
sudo systemctl start/stop/restart void-reader
sudo systemctl status void-reader
sudo journalctl -u void-reader -f
```

#### Railway Deployment
```bash
# Deploy to Railway
railway up

# Set environment variables in Railway dashboard:
# HTTP_PORT=8080
# SSH_PORT=2222
# SSH_PASSWORD=Amigos4Life!

# Configure TCP Proxy in Railway dashboard to port 2222
# Connect via: ssh trolley.proxy.rlwy.net -p 10120
```

#### Container Deployment
```bash
# Docker Compose (recommended)
cd ssh-reader && docker-compose up -d

# Or direct Docker
docker build -t void-reader ssh-reader/
docker run -d -p 8080:8080 -p 23234:23234 void-reader
```

#### Development
```bash
# Run in development mode
go run .

# Run tests
go test ./...

# Debug with verbose logging
DEBUG=1 go run .
```

### Book Content Management Commands


#### Generate PDF for Amazon KDP
```bash
./markdown_to_kdp_pdf.rb book1_void_reavers_source void_reavers.pdf
# Creates print-ready PDF with proper formatting
```

#### Generate EPUB for E-readers
```bash
./markdown_to_epub.rb book1_void_reavers_source void_reavers.epub
# Creates EPUB file for Kindle and other e-readers
```


## SSH Reader Application Details

### Architecture
- **Dual Servers**: HTTP on port 8080 (90s homepage), SSH on port 2222/23234
- **SSH Server**: Uses Charm's Wish library with password authentication
- **TUI Interface**: Split-view design showing all 10 books with summaries
- **Book Loading**: Loads from `book1_void_reavers_source/chapters/` Markdown files
- **Progress System**: JSON-based user progress persistence
- **Railway Support**: TCP proxy compatible with environment-based port configuration

### Key Features
- Split-view menu showing all 10 Void Chronicles books
- Spoiler-free book summaries to intrigue readers
- 90s-style homepage disguise for public deployments
- Beautiful terminal interface with emojis and colors
- Chapter navigation with keyboard shortcuts (h/l, ←/→)
- Progress tracking with auto-save on chapter change
- Bookmark system (press 'b' while reading)
- Railway deployment with TCP proxy support

### User Interface States
1. **Main Menu** - Split-view with book library on left, details on right
2. **Chapter List** - Browse all chapters with completion indicators
3. **Reading View** - Main reading interface with scrolling
4. **Progress View** - Statistics, completion percentage, bookmarks
5. **About View** - Book and series information
6. **HTTP Homepage** - 90s-style "Bob's Personal Homepage" disguise

### File Locations
- **User Progress**: `.void_reader_data/username.json`
- **SSH Host Keys**: `.ssh/id_ed25519`
- **Book Content**: `book1_void_reavers_source/` (Markdown source)
- **Logs**: systemd journal or console output

## Book Content Structure

### Book Files
- **Main book file**: `book.md` - Contains book metadata and structure
- **Chapters**: `chapters/chapter-01.md` through `chapters/chapter-20.md` - Individual chapter content
- **Metadata**: `metadata.yaml` - Book information for publishing
- **Conversion scripts**: Ruby scripts for PDF and EPUB generation

### Markdown Formatting Conventions
- Chapter titles use `# Chapter Title`
- Italics use `*text*`
- Bold uses `**text**`
- Ship names are italicized: `*Ship Name*`
- Scene breaks use `* * *`
- Quotes use standard double quotes

### Conversion Pipeline Details
The Ruby conversion scripts handle:
- Markdown → PDF via Pandoc with custom LaTeX templates
- Markdown → EPUB with proper metadata
- Multiple trim sizes for print publishing
- Professional typography and formatting

## Working with the Project

### SSH Reader Development
1. Make changes to Go source files in `ssh-reader/` directory
2. Test with `./run.sh` for local development (HTTP on 8080, SSH on 23234)
3. Build with `./build.sh` for binary generation
4. Deploy to Railway with `railway up`
5. Configure Railway TCP proxy for SSH access
6. Check logs with `railway logs` or local console output

### Book Content Editing
1. Edit the Markdown source files in `book1_void_reavers_source/chapters/`
2. Follow the style guide in `MARKDOWN_STYLE_GUIDE.md`
3. Generate PDF/EPUB to review formatting
4. Use any text editor - no special tools required

### Documentation
Complete documentation is available in the `docs/` directory:
- User guides, installation, configuration
- Deployment options for 12+ platforms
- Development guide and API reference
- Troubleshooting and contributing guides

## Story Context

The book follows a 50-year timeline chronicling humanity's evolution from chaotic expansion to galactic citizenship, with pirates serving as both catalyst and conscience for this transformation. The alien "Architects" (formerly called "Watchers") serve as probation officers testing humanity's worthiness for cosmic citizenship.

## Key Technologies

- **Go 1.21+** - Main application language
- **Bubbletea** - Terminal user interface framework
- **Lipgloss** - Styling and layout for TUI
- **Wish** - SSH server middleware
- **Keygen** - Go-native SSH key generation
- **Markdown** - Book source format (migrated from LaTeX)
- **Ruby** - Conversion scripts for PDF/EPUB
- **Railway** - Cloud deployment platform with TCP proxy
- **Docker** - Containerization support