# Void Reavers SSH Reader 🚀

A terminal-based book reader that allows users to read "Void Reavers" over SSH using Go, Bubbletea, and Wish.

## Features

- 📖 **SSH-based Reading**: Connect via SSH to read the book in a beautiful terminal interface
- 🎨 **Rich TUI**: Built with Bubbletea for smooth, responsive terminal UI
- 📚 **Chapter Navigation**: Easy navigation between chapters with keyboard shortcuts
- 🔄 **Auto-loading**: Automatically loads from Markdown or LaTeX sources
- 📱 **Responsive**: Adapts to different terminal sizes
- 🎯 **Progress Tracking**: See current chapter and progress through the book

## Quick Start

### Prerequisites

- Go 1.21 or later
- SSH key pair (will be generated if not present)

### Installation

```bash
# Install dependencies
go mod tidy

# Generate SSH host key (if needed)
ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ""

# Build and run
go build -o void-reader
./void-reader
```

### Connecting

Once the server is running, connect via SSH:

```bash
ssh localhost -p 23234
```

## Controls

### Main Menu
- `↑/↓` or `k/j`: Navigate menu items
- `Enter` or `Space`: Select item
- `q`: Quit

### Chapter List
- `↑/↓` or `k/j`: Navigate chapters
- `Enter` or `Space`: Read selected chapter
- `Esc`: Back to main menu
- `q`: Quit

### Reading View
- `↑/↓` or `k/j`: Scroll line by line
- `Page Up/Down` or `Ctrl+B/F`: Scroll by page
- `Space`: Page down
- `Home/End` or `g/G`: Go to beginning/end
- `←/→` or `h/l`: Previous/next chapter
- `p/n`: Previous/next chapter (alternative)
- `Esc`: Back to main menu
- `q`: Quit

## Configuration

### Server Settings

Edit the constants in `main.go`:

```go
const (
    host = "localhost"  // Change to "0.0.0.0" for external access
    port = "23234"      // Change port as needed
)
```

### SSH Security

For production use:
- Use proper SSH host keys
- Consider authentication methods
- Restrict access via firewall rules
- Use non-standard ports

## Book Loading

The reader automatically attempts to load the book in this order:

1. **Markdown**: Looks for `book1_void_reavers/markdown/chapter*.md`
2. **LaTeX**: Falls back to `book1_void_reavers/chapter*.tex`

### Supported Formats

#### Markdown
- Chapter files: `chapter01.md`, `chapter02.md`, etc.
- Format: `# Chapter Title` followed by content

#### LaTeX
- Chapter files: `chapter01.tex`, `chapter02.tex`, etc.
- Format: `\chapter{Chapter Title}` followed by content
- Automatically converts LaTeX formatting to plain text

## Development

### Project Structure

```
void-reader/
├── main.go              # Main application and SSH server
├── book.go              # Book loading and parsing logic
├── go.mod               # Go module dependencies
├── .ssh/                # SSH host keys
│   └── id_ed25519
└── book1_void_reavers/  # Book content
    ├── markdown/        # Markdown chapters (preferred)
    └── *.tex           # LaTeX chapters (fallback)
```

### Adding Features

The code is structured for easy extension:

- **New book formats**: Add parsers in `book.go`
- **UI improvements**: Modify views in `main.go`
- **Additional navigation**: Add to update functions
- **Progress tracking**: Extend the model struct

### Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling and layout
- `github.com/charmbracelet/wish` - SSH server middleware
- `github.com/muesli/reflow` - Text wrapping utilities

## Deployment

### Local Development
```bash
go run . 
```

### Production Build
```bash
# Build optimized binary
go build -ldflags="-s -w" -o void-reader

# Run as service (example with systemd)
sudo systemctl enable --now void-reader.service
```

### Docker Deployment
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o void-reader

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/void-reader .
COPY --from=builder /app/book1_void_reavers ./book1_void_reavers
EXPOSE 23234
CMD ["./void-reader"]
```

## Troubleshooting

### Connection Issues
- Check if port 23234 is available: `netstat -tlnp | grep 23234`
- Verify SSH client can connect: `ssh -v localhost -p 23234`
- Check firewall settings if connecting remotely

### Book Loading Issues
- Ensure chapter files exist in expected locations
- Check file permissions for read access
- Verify file encoding (should be UTF-8)

### Performance
- Large books may need pagination adjustments
- Terminal size affects text wrapping
- SSH connection quality impacts responsiveness

## Future Enhancements

- [ ] User authentication and sessions
- [ ] Bookmark and progress persistence
- [ ] Multiple book support
- [ ] Search within book content
- [ ] Themes and color customization
- [ ] Reading statistics and analytics
- [ ] Export functionality
- [ ] Multi-language support

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is part of the Void Chronicles universe. The reader application is open source, while the book content follows its own licensing terms.

---

*"In the void between stars, even pirates can find their way home."* 🌌