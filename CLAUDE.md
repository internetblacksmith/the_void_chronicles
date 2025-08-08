# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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
The Go application consists of:
- `main.go` (546 lines) - SSH server, TUI interface, user interaction
- `book.go` (242 lines) - Book loading, parsing, format conversion
- `progress.go` (152 lines) - User progress tracking, bookmarks, statistics
- `go.mod/go.sum` - Dependencies (Bubbletea, Lipgloss, Wish, etc.)
- Build and deployment scripts (`build.sh`, `run.sh`, `deploy.sh`)
- Docker configuration (`Dockerfile`, `docker-compose.yml`)
- Systemd service configuration (`systemd/void-reader.service`)
- Comprehensive documentation in `docs/` directory

## Essential Commands

### SSH Reader Application Commands

#### Quick Setup and Run
```bash
# Build the SSH reader application
./build.sh

# Start the SSH server (localhost:23234)
./run.sh

# Connect to read the book
ssh localhost -p 23234
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

#### Container Deployment
```bash
# Docker Compose (recommended)
docker-compose up -d

# Or direct Docker
docker build -t void-reader .
docker run -d -p 23234:23234 void-reader
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
- **SSH Server**: Uses Charm's Wish library for SSH handling
- **TUI Interface**: Built with Bubbletea for responsive terminal UI
- **Styling**: Lipgloss for colors, borders, and layout
- **Book Loading**: Supports both Markdown and LaTeX sources
- **Progress System**: JSON-based user progress persistence
- **Multi-User**: Individual progress tracking per SSH user

### Key Features
- Beautiful terminal interface with emojis and colors
- Chapter navigation with keyboard shortcuts (h/l, ←/→)
- Progress tracking with auto-save on chapter change
- Bookmark system (press 'b' while reading)
- Multi-user support with separate progress files
- Responsive design adapting to terminal size
- Production-ready with systemd service integration

### User Interface States
1. **Main Menu** - Continue reading, chapter list, progress, about, exit
2. **Chapter List** - Browse all chapters with completion indicators
3. **Reading View** - Main reading interface with scrolling
4. **Progress View** - Statistics, completion percentage, bookmarks
5. **About View** - Book and application information

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
1. Make changes to Go source files (`main.go`, `book.go`, `progress.go`)
2. Test with `go run .` for development
3. Build with `./build.sh` for testing
4. Deploy with `./deploy.sh` for production
5. Check logs with `sudo journalctl -u void-reader -f`

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
- **LaTeX** - Book source format
- **Ruby** - Conversion scripts
- **Docker** - Containerization
- **Systemd** - Linux service management