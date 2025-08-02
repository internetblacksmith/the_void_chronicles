# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This repository contains **two main components**:

1. **"The Void Chronicles" Book Series** - Starting with Book 1: "Void Reavers: A Tale of Space Pirates and Cosmic Plunder" - a complete 20-chapter science fiction novel about Captain Zara "Bloodhawk" Vega's transformation from pirate to diplomat in a universe where humanity must prove itself to alien Architects.

2. **Void Reavers SSH Reader** - A Go application that creates an SSH server allowing users to read books through a beautiful terminal user interface (TUI) built with Bubbletea and Wish.

## Architecture & File Structure

### Book Content Structure
The series is organized with each book in its own directory:
- `book1_void_reavers/` - Contains the first complete book
- `void_chronicles_series_bible.md` - Master document with series overview and future book plans
- Book conversion scripts in root directory for processing any book

Each book exists in multiple formats through a conversion pipeline:

1. **Source**: LaTeX files (`book.tex` + `chapter01.tex` through `chapter20.tex`)
2. **Intermediate**: Markdown files (in `markdown/` directory within each book folder)
3. **Output**: HTML files (in `html/` directory within each book folder) for web viewing and PDF generation

The LaTeX source is the canonical version. All other formats are generated from it.

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

#### Generate Markdown from LaTeX
```bash
ruby convert_to_md.rb
# Script will detect book directories and prompt for which book to convert
# Choose 'all' to convert all books, or specify a book name/number
```

#### Generate HTML for PDF Creation
```bash
ruby markdown_to_html.rb
# Will work with the current directory structure
```

#### Generate PDF (requires pandoc)
```bash
ruby convert_to_pdf.rb
# For books with pandoc available
```

#### Generate EPUB for Kindle
```bash
ruby convert_to_epub.rb
# Or use the shell script:
./convert_to_epub.sh

# Creates an EPUB file ready for Kindle
# Located at: book1_void_reavers/void_reavers.epub
```

#### View a Book
Navigate to the book directory and open `html/void_reavers_complete.html` in a web browser.
Use the browser's print function to create PDF.

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
- **Book Content**: `book1_void_reavers/` (markdown preferred, LaTeX fallback)
- **Logs**: systemd journal or console output

## Book Content Structure

### Book Files
- **Main book file**: `book.tex` - Contains document structure and chapter includes
- **Chapters**: `chapter01.tex` through `chapter20.tex` - Individual chapter content
- **Conversion scripts**: Ruby-based converters that handle LaTeX → Markdown → HTML transformation

### LaTeX Formatting Conventions
- Chapter titles use `\chapter{Title}`
- Italics use `\textit{text}`
- Bold uses `\textbf{text}`
- Ship names are italicized: `\textit{Ship Name}`
- Em-dashes use `---`
- Quotes use `` ` ` text ' ' `` for double quotes and `` ` text ' `` for single quotes

### Conversion Pipeline Details
The Ruby conversion scripts handle:
- LaTeX markup → Markdown syntax transformation
- Character encoding and special character handling
- Automatic chapter detection from `book.tex` or filesystem
- Generation of complete book files and individual chapters
- HTML output with print-optimized CSS

## Working with the Project

### SSH Reader Development
1. Make changes to Go source files (`main.go`, `book.go`, `progress.go`)
2. Test with `go run .` for development
3. Build with `./build.sh` for testing
4. Deploy with `./deploy.sh` for production
5. Check logs with `sudo journalctl -u void-reader -f`

### Book Content Editing
1. Edit the LaTeX source files directly
2. Run conversion scripts to update other formats
3. Review HTML output for formatting issues
4. The markdown versions are generated - don't edit them directly

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