# Quick Start Guide

Get up and running with the Void Reavers SSH Reader in just a few minutes!

## 🎯 Prerequisites

Before you begin, ensure you have:

- **Go 1.21+** installed ([Download Go](https://golang.org/dl/))
- **SSH client** (built into most systems)
- **Terminal** with UTF-8 support
- **Port 23234** available (or ability to configure different port)

## ⚡ 5-Minute Setup

### 1. Get the Code

```bash
# If you have the source code directory
cd /path/to/void-reavers-reader

# Or if cloning from a repository
git clone <repository-url>
cd void-reavers-reader
```

### 2. Build the Application

```bash
# This script will:
# - Download Go dependencies
# - Generate SSH host keys
# - Build the binary
# - Set up data directories
./build.sh
```

Expected output:
```
🚀 Building Void Reavers SSH Reader...
==================================
🔑 Generating SSH host key...
✅ SSH host key generated
📁 Created data directory for user progress
📦 Downloading Go dependencies...
🔨 Building application...
✅ Build complete!
```

### 3. Start the Server

```bash
# This will start the SSH server on port 23234
./run.sh
```

You should see:
```
🚀 Starting Void Reavers SSH Reader...
=====================================
📚 Book: Void Reavers
🌐 Server: localhost:23234
🔑 SSH Key: .ssh/id_ed25519
💾 Data Dir: .void_reader_data/

🎯 To connect: ssh localhost -p 23234

Starting server...
```

### 4. Connect and Read!

Open a **new terminal** and connect:

```bash
ssh localhost -p 23234
```

You'll see the main menu:

```
🚀 VOID REAVERS 🚀
A Tale of Space Pirates and Cosmic Plunder

▶ 📖 Continue Reading
  📚 Chapter List  
  📊 Progress
  ℹ️  About
  🚪 Exit

↑/↓: navigate • enter: select • q: quit
```

## 🎮 Basic Controls

### Main Menu
- `↑/↓` or `k/j`: Navigate options
- `Enter` or `Space`: Select option
- `q`: Quit application

### Reading View
- `↑/↓` or `k/j`: Scroll line by line
- `Space` or `Page Down`: Scroll page down
- `Page Up`: Scroll page up
- `h/←` or `p`: Previous chapter
- `l/→` or `n`: Next chapter
- `b`: Add bookmark at current position
- `g`: Go to beginning of chapter
- `G`: Go to end of chapter
- `Esc`: Return to main menu

### Chapter List
- `↑/↓` or `k/j`: Navigate chapters
- `Enter`: Jump to selected chapter
- `Esc`: Back to main menu

## 📊 Your First Reading Session

1. **Start Reading**: Select "📖 Continue Reading" from the main menu
2. **Navigate**: Use arrow keys or `h/l` to move between chapters
3. **Bookmark**: Press `b` to bookmark interesting passages
4. **Check Progress**: Press `Esc` to go back, then select "📊 Progress"
5. **Continue Later**: Your position is automatically saved!

## 🔧 Quick Customization

### Change the Port

Edit `main.go` and change:
```go
const (
    host = "localhost"
    port = "23234"      // Change this
)
```

Then rebuild:
```bash
./build.sh
```

### Connect from Remote Machine

If you want to allow external connections:

1. Change host in `main.go`:
```go
const (
    host = "0.0.0.0"    // Listen on all interfaces
    port = "23234"
)
```

2. Rebuild and restart:
```bash
./build.sh
./run.sh
```

3. Connect from remote machine:
```bash
ssh your-server-ip -p 23234
```

## 🐛 Common Issues

### "Permission denied" when connecting
The SSH server might not be running. Check if `./run.sh` is still active.

### "Connection refused"
Port 23234 might be in use. Try changing the port in `main.go`.

### "Book content not found"
Ensure the `book1_void_reavers/` directory exists with chapter files.

### SSH key warnings
This is normal for first connection. The app generates its own SSH host key.

## 🚀 Next Steps

Now that you're up and running:

- **Read the full book**: Navigate through all 20 chapters of Void Reavers
- **Explore features**: Try bookmarking, check your progress statistics
- **Customize**: See [Configuration Guide](configuration.md) for advanced options
- **Deploy**: Check [Deployment Guide](deployment.md) for production setup
- **Develop**: Read [Development Guide](development.md) to add features

## 💡 Pro Tips

- **Multiple connections**: You can have multiple SSH sessions to the same server
- **Terminal size**: Resize your terminal for better reading experience
- **Progress tracking**: Each user gets individual progress tracking
- **Bookmarks**: Use bookmarks to mark favorite quotes or important passages
- **Keyboard shortcuts**: Learn the shortcuts for faster navigation

## 🆘 Need Help?

- **Issues**: Check [Troubleshooting Guide](troubleshooting.md)
- **Configuration**: See [Configuration Guide](configuration.md)
- **Features**: Read [User Guide](user-guide.md) for detailed feature explanations

---

**Congratulations!** 🎉 You're now ready to explore the universe of Void Reavers through your terminal. Enjoy the reading experience!

*"Every great journey begins with a single step into the void."* ✨