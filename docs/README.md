# Void Reavers SSH Reader Documentation

Welcome to the complete documentation for the Void Reavers SSH Reader - a terminal-based book reader that allows users to read books over SSH connections using a beautiful TUI interface.

## 📚 Table of Contents

- [Quick Start Guide](quick-start.md)
- [Installation Guide](installation.md)
- [User Guide](user-guide.md)
- [Configuration](configuration.md)
- [Deployment Guide](deployment.md)
- [Deployment Options](deployment-options.md)
- [Development Guide](development.md)
- [API Reference](api-reference.md)
- [Troubleshooting](troubleshooting.md)
- [Contributing](contributing.md)

## 📖 What is Void Reavers SSH Reader?

The Void Reavers SSH Reader is a Go application that creates an SSH server allowing users to read books through a terminal user interface (TUI). Built with modern Go libraries including Bubbletea for the UI and Wish for SSH handling, it provides a unique and engaging way to read literature.

### Key Features

- **🌐 SSH-Based Access**: Connect from anywhere via SSH
- **🎨 Beautiful TUI**: Rich terminal interface with colors and styling
- **📊 Progress Tracking**: Automatic save/resume functionality
- **🔖 Bookmarking System**: Mark favorite passages and locations
- **📱 Responsive Design**: Adapts to different terminal sizes
- **🔄 Multi-Format Support**: Reads from Markdown or LaTeX sources
- **👥 Multi-User Support**: Individual progress tracking per user
- **🐳 Container Ready**: Docker and Docker Compose support
- **⚙️ System Service**: Systemd integration for production

### Use Cases

- **Personal Reading**: Read books in a distraction-free terminal environment
- **Remote Access**: Access your library from any device with SSH
- **Educational**: Great for technical documentation or coding books
- **Unique Experience**: Novel way to experience literature
- **Offline Reading**: No web browser required, works over basic SSH

## 🚀 Quick Start

```bash
# Clone and build
git clone <repository>
cd void-reavers-reader
./build.sh

# Start the server
./run.sh

# Connect from another terminal
ssh localhost -p 23234
```

## 🏗️ Architecture Overview

```
┌─────────────────┐    SSH     ┌─────────────────┐
│   SSH Client    │ ◄────────► │   Wish Server   │
│  (Terminal)     │  Port 23234│  (Go App)       │
└─────────────────┘            └─────────────────┘
                                        │
                                        ▼
                               ┌─────────────────┐
                               │  Bubbletea TUI  │
                               │   (Interface)   │
                               └─────────────────┘
                                        │
                        ┌───────────────┼───────────────┐
                        ▼               ▼               ▼
                ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
                │Book Loader  │ │Progress Mgr │ │  Styling    │
                │(MD/LaTeX)   │ │(JSON Files) │ │(Lipgloss)   │
                └─────────────┘ └─────────────┘ └─────────────┘
```

## 📋 System Requirements

### Minimum Requirements
- **OS**: Linux, macOS, Windows (with WSL)
- **Go**: Version 1.21 or later
- **Memory**: 50MB RAM
- **Storage**: 100MB for application and book content
- **Network**: Port 23234 available (configurable)

### Recommended Requirements
- **CPU**: 1 core at 1GHz
- **Memory**: 100MB RAM for smooth operation
- **Storage**: 500MB for multiple books and user data
- **Terminal**: Modern terminal with UTF-8 support

## 🛡️ Security Considerations

- SSH host key authentication
- Non-privileged user execution
- Sandboxed file access
- Read-only book content
- Isolated user data directories

## 📄 License

This project is part of the Void Chronicles universe. The reader application is open source under the MIT License, while book content follows its own licensing terms.

## 🤝 Support

- **Documentation**: See individual guide files in this directory
- **Issues**: Report bugs and feature requests on GitHub
- **Discussions**: Join community discussions for help and ideas

---

*"In the void between stars, even code can tell a story."* 🌌

---

**Next Steps:**
- New to the system? Start with the [Quick Start Guide](quick-start.md)
- Need to install? Check the [Installation Guide](installation.md)
- Ready to deploy? See the [Deployment Guide](deployment.md)
- Want to contribute? Read the [Development Guide](development.md)