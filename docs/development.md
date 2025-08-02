# Development Guide

Complete guide for developers who want to understand, modify, or contribute to the Void Reavers SSH Reader.

## 🏗️ Architecture Overview

### System Architecture

```
┌─────────────────┐    SSH/TCP     ┌─────────────────┐
│   SSH Client    │ ◄────────────► │   Wish Server   │
│  (Any Terminal) │  Port 23234    │  (Entry Point)  │
└─────────────────┘                └─────────────────┘
                                            │
                                            ▼
                                   ┌─────────────────┐
                                   │  Bubbletea App  │
                                   │   (TUI Logic)   │
                                   └─────────────────┘
                                            │
                    ┌───────────────────────┼───────────────────────┐
                    ▼                       ▼                       ▼
            ┌─────────────┐        ┌─────────────┐        ┌─────────────┐
            │Book Manager │        │Progress Mgr │        │   Styling   │
            │(Content)    │        │(User Data)  │        │ (Lipgloss)  │
            └─────────────┘        └─────────────┘        └─────────────┘
                    │                       │
                    ▼                       ▼
            ┌─────────────┐        ┌─────────────┐
            │   Files     │        │JSON Storage │
            │(MD/LaTeX)   │        │(User Data)  │
            └─────────────┘        └─────────────┘
```

### Component Responsibilities

| Component | Responsibility | Files |
|-----------|----------------|-------|
| **Wish Server** | SSH connection handling, session management | `main.go` |
| **Bubbletea App** | TUI interface, user interaction, navigation | `main.go` |
| **Book Manager** | Content loading, parsing, format conversion | `book.go` |
| **Progress Manager** | User data persistence, bookmarks, statistics | `progress.go` |
| **Styling** | Visual appearance, themes, responsive layout | `main.go` (styles) |

## 🚀 Getting Started

### Development Environment Setup

#### Prerequisites
```bash
# Go 1.21+
go version

# Git for version control
git --version

# SSH client for testing
ssh -V

# Optional: Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
```

#### Project Setup
```bash
# Clone repository
git clone <repository-url>
cd void-reavers-reader

# Install dependencies
go mod tidy

# Verify build
go build -o void-reader-dev

# Run tests (if any)
go test ./...
```

#### Development Workflow
```bash
# Development with auto-restart
go run . &
PID=$!

# Make changes, then restart
kill $PID
go run . &
```

### Project Structure

```
void-reavers-reader/
├── main.go              # Main application, SSH server, TUI
├── book.go              # Book loading and content parsing
├── progress.go          # User progress and data management
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── build.sh             # Build automation script
├── run.sh               # Development run script
├── deploy.sh            # Production deployment script
├── Dockerfile           # Container build instructions
├── docker-compose.yml   # Multi-container orchestration
├── systemd/             # System service configuration
│   └── void-reader.service
├── docs/                # Documentation
│   ├── README.md
│   ├── quick-start.md
│   ├── user-guide.md
│   ├── configuration.md
│   ├── deployment.md
│   ├── troubleshooting.md
│   └── development.md   # This file
├── book1_void_reavers/  # Book content
│   ├── markdown/        # Preferred format
│   │   ├── chapter01.md
│   │   └── ...
│   ├── chapter01.tex    # LaTeX source
│   └── ...
├── .ssh/                # SSH host keys (generated)
├── .void_reader_data/   # User progress data (generated)
└── README_ssh_reader.md # Quick reference
```

## 🧩 Core Components

### 1. SSH Server (Wish Integration)

**File**: `main.go` (lines 1-65)

**Key Functions**:
- `main()`: Server initialization and lifecycle
- `teaHandler()`: Session creation for each SSH connection

**Architecture**:
```go
// Server creation with middleware stack
s, err := wish.NewServer(
    wish.WithAddress(net.JoinHostPort(host, port)),
    wish.WithHostKeyPath(".ssh/id_ed25519"),
    wish.WithMiddleware(
        bubbletea.Middleware(teaHandler),  // TUI integration
        logging.Middleware(),              // Request logging
    ),
)
```

**Extension Points**:
- Add authentication middleware
- Custom logging middleware
- Rate limiting middleware
- Metrics collection middleware

### 2. TUI Application (Bubbletea)

**File**: `main.go` (lines 66-548)

**Key Structures**:
```go
type model struct {
    state           state              // Current UI state
    book            *Book             // Book content
    progress        *UserProgress     // User's reading progress
    progressManager *ProgressManager  // Progress persistence
    // ... UI state fields
}
```

**State Machine**:
```go
const (
    menuView state = iota        // Main menu
    chapterListView             // Chapter selection
    readingView                 // Reading interface
    aboutView                   // Information screen
    progressView                // Progress statistics
)
```

**Key Methods**:
- `Init()`: Initialize application state
- `Update()`: Handle user input and state transitions
- `View()`: Render current UI state

### 3. Book Content System

**File**: `book.go`

**Core Types**:
```go
type Book struct {
    Title    string
    Author   string
    Chapters []Chapter
}

type Chapter struct {
    Title   string
    Content string
}
```

**Loading Pipeline**:
1. **Detection**: Try Markdown first, fall back to LaTeX
2. **Parsing**: Extract chapters and metadata
3. **Conversion**: Convert to plain text with formatting
4. **Caching**: Keep parsed content in memory

**Format Support**:
- **Markdown**: Primary format, fast parsing
- **LaTeX**: Source format, converted to plain text
- **Extensible**: Easy to add new formats

### 4. Progress Management

**File**: `progress.go`

**Data Structure**:
```go
type UserProgress struct {
    Username       string            // User identifier
    CurrentChapter int               // Last read chapter
    ScrollOffset   int               // Position in chapter
    LastRead       time.Time         // Last reading session
    ChapterProgress map[int]bool     // Completed chapters
    Bookmarks      []Bookmark        // Saved positions
    ReadingTime    time.Duration     // Total reading time
}
```

**Persistence**:
- **Format**: JSON files per user
- **Location**: `.void_reader_data/`
- **Auto-save**: On chapter change and exit
- **Backup**: Manual backup procedures

## 🔧 Development Tasks

### Adding New Features

#### 1. Add New UI Screen

Example: Adding a search screen

```go
// 1. Add to state enum
const (
    menuView state = iota
    chapterListView
    readingView
    aboutView
    progressView
    searchView        // New state
)

// 2. Add to model struct
type model struct {
    // ... existing fields
    searchQuery string
    searchResults []SearchResult
}

// 3. Add update handler
func (m model) updateSearch(msg tea.KeyMsg) (model, tea.Cmd) {
    switch msg.String() {
    case "enter":
        // Perform search
        m.searchResults = searchInBook(m.book, m.searchQuery)
    case "esc":
        return m.toMenu(), nil
    default:
        // Handle search input
        m.searchQuery += msg.String()
    }
    return m, nil
}

// 4. Add view renderer
func (m model) viewSearch() string {
    header := headerStyle.Width(m.width-2).Render("🔍 SEARCH")
    
    searchBox := fmt.Sprintf("Search: %s", m.searchQuery)
    
    var results []string
    for _, result := range m.searchResults {
        results = append(results, fmt.Sprintf("Chapter %d: %s", 
            result.Chapter, result.Preview))
    }
    
    content := lipgloss.JoinVertical(lipgloss.Left, 
        searchBox, "", lipgloss.JoinVertical(lipgloss.Left, results...))
    
    return lipgloss.JoinVertical(lipgloss.Top, header, "", content)
}

// 5. Add to main Update/View functions
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch m.state {
        // ... existing cases
        case searchView:
            return m.updateSearch(msg)
        }
    }
    return m, nil
}

func (m model) View() string {
    switch m.state {
    // ... existing cases
    case searchView:
        return m.viewSearch()
    }
    return ""
}
```

#### 2. Add New Book Format Support

Example: Adding EPUB support

```go
// 1. Add to book.go
import "github.com/go-shiori/go-epub"

func loadFromEPUB(bookPath string) (*Book, error) {
    epub, err := epub.Open(bookPath)
    if err != nil {
        return nil, err
    }
    defer epub.Close()
    
    var chapters []Chapter
    for _, item := range epub.Spine.ItemRefs {
        content, err := epub.ReadFile(item.IDRef)
        if err != nil {
            continue
        }
        
        chapter := parseHTMLChapter(string(content))
        chapters = append(chapters, chapter)
    }
    
    return &Book{
        Title:    epub.Title,
        Author:   epub.Creator,
        Chapters: chapters,
    }, nil
}

// 2. Update LoadBook function
func LoadBook(bookDir string) (*Book, error) {
    // Try EPUB first
    if epubPath := filepath.Join(bookDir, "book.epub"); 
       fileExists(epubPath) {
        if book, err := loadFromEPUB(epubPath); err == nil {
            return book, nil
        }
    }
    
    // Existing markdown and LaTeX loaders
    // ...
}
```

#### 3. Add Authentication

```go
// 1. Add to main.go imports
import "github.com/charmbracelet/ssh"

// 2. Add public key authentication
func loadAuthorizedKeys(file string) map[string]ssh.PublicKey {
    keys := make(map[string]ssh.PublicKey)
    
    data, err := ioutil.ReadFile(file)
    if err != nil {
        return keys
    }
    
    for len(data) > 0 {
        pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(data)
        if err != nil {
            break
        }
        keys[string(pubKey.Marshal())] = pubKey
        data = rest
    }
    
    return keys
}

// 3. Add to server middleware
wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
    authorizedKeys := loadAuthorizedKeys(".ssh/authorized_keys")
    _, authorized := authorizedKeys[string(key.Marshal())]
    return authorized
}),
```

### Performance Optimizations

#### 1. Book Content Caching

```go
import "sync"

type BookCache struct {
    cache map[string]*Book
    mutex sync.RWMutex
}

func NewBookCache() *BookCache {
    return &BookCache{
        cache: make(map[string]*Book),
    }
}

func (bc *BookCache) Get(path string) (*Book, bool) {
    bc.mutex.RLock()
    defer bc.mutex.RUnlock()
    book, exists := bc.cache[path]
    return book, exists
}

func (bc *BookCache) Set(path string, book *Book) {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()
    bc.cache[path] = book
}

// Global cache instance
var bookCache = NewBookCache()

// Update LoadBook to use cache
func LoadBook(bookDir string) (*Book, error) {
    if cached, exists := bookCache.Get(bookDir); exists {
        return cached, nil
    }
    
    book, err := loadBookFromDisk(bookDir)
    if err != nil {
        return nil, err
    }
    
    bookCache.Set(bookDir, book)
    return book, nil
}
```

#### 2. Progress Auto-save Optimization

```go
type ProgressManager struct {
    dataDir    string
    saveQueue  chan *UserProgress
    batchTimer *time.Timer
}

func NewProgressManager() *ProgressManager {
    pm := &ProgressManager{
        dataDir:   ".void_reader_data",
        saveQueue: make(chan *UserProgress, 100),
    }
    
    // Start background saver
    go pm.batchSaver()
    
    return pm
}

func (pm *ProgressManager) batchSaver() {
    batch := make(map[string]*UserProgress)
    ticker := time.NewTicker(5 * time.Second)
    
    for {
        select {
        case progress := <-pm.saveQueue:
            batch[progress.Username] = progress
            
        case <-ticker.C:
            // Save all batched progress
            for _, progress := range batch {
                pm.saveProgressToDisk(progress)
            }
            batch = make(map[string]*UserProgress)
        }
    }
}

func (pm *ProgressManager) SaveProgressAsync(progress *UserProgress) {
    select {
    case pm.saveQueue <- progress:
        // Queued successfully
    default:
        // Queue full, save immediately
        pm.saveProgressToDisk(progress)
    }
}
```

### Testing

#### Unit Tests

```go
// book_test.go
package main

import (
    "testing"
    "os"
    "path/filepath"
)

func TestLoadBook(t *testing.T) {
    // Create temporary book directory
    tempDir := t.TempDir()
    bookDir := filepath.Join(tempDir, "test_book")
    os.MkdirAll(filepath.Join(bookDir, "markdown"), 0755)
    
    // Create test chapter
    chapterContent := "# Test Chapter\n\nThis is test content."
    err := os.WriteFile(
        filepath.Join(bookDir, "markdown", "chapter01.md"),
        []byte(chapterContent),
        0644,
    )
    if err != nil {
        t.Fatal(err)
    }
    
    // Test loading
    book, err := LoadBook(bookDir)
    if err != nil {
        t.Fatalf("Failed to load book: %v", err)
    }
    
    if len(book.Chapters) != 1 {
        t.Errorf("Expected 1 chapter, got %d", len(book.Chapters))
    }
    
    if book.Chapters[0].Title != "Test Chapter" {
        t.Errorf("Expected 'Test Chapter', got '%s'", book.Chapters[0].Title)
    }
}

func TestProgressManager(t *testing.T) {
    pm := NewProgressManager()
    
    progress := &UserProgress{
        Username:       "testuser",
        CurrentChapter: 5,
        ScrollOffset:   100,
    }
    
    err := pm.SaveProgress(progress)
    if err != nil {
        t.Fatalf("Failed to save progress: %v", err)
    }
    
    loaded, err := pm.LoadProgress("testuser")
    if err != nil {
        t.Fatalf("Failed to load progress: %v", err)
    }
    
    if loaded.CurrentChapter != 5 {
        t.Errorf("Expected chapter 5, got %d", loaded.CurrentChapter)
    }
}
```

#### Integration Tests

```go
// integration_test.go
package main

import (
    "testing"
    "net"
    "time"
    "golang.org/x/crypto/ssh"
)

func TestSSHConnection(t *testing.T) {
    // Start server in background
    go func() {
        main()
    }()
    
    // Wait for server to start
    time.Sleep(2 * time.Second)
    
    // Test SSH connection
    conn, err := net.Dial("tcp", "localhost:23234")
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    
    // Basic SSH handshake
    config := &ssh.ClientConfig{
        User: "testuser",
        Auth: []ssh.AuthMethod{
            ssh.Password(""), // No password required
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    
    sshConn, chans, reqs, err := ssh.NewClientConn(conn, "localhost:23234", config)
    if err != nil {
        t.Fatalf("SSH handshake failed: %v", err)
    }
    defer sshConn.Close()
    
    client := ssh.NewClient(sshConn, chans, reqs)
    defer client.Close()
    
    // Test successful connection
    session, err := client.NewSession()
    if err != nil {
        t.Fatalf("Failed to create session: %v", err)
    }
    defer session.Close()
}
```

#### Benchmark Tests

```go
// benchmark_test.go
package main

import (
    "testing"
    "strings"
)

func BenchmarkTextWrapping(b *testing.B) {
    longText := strings.Repeat("This is a long line of text that needs wrapping. ", 1000)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wrapText(longText, 80)
    }
}

func BenchmarkBookLoading(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        LoadBook("book1_void_reavers")
    }
}

func BenchmarkProgressSaving(b *testing.B) {
    pm := NewProgressManager()
    progress := &UserProgress{
        Username:       "benchuser",
        CurrentChapter: 10,
        ScrollOffset:   500,
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pm.SaveProgress(progress)
    }
}
```

### Debugging

#### Debug Mode

```go
// Add debug flag
var debugMode = os.Getenv("DEBUG") == "1"

func debugLog(format string, args ...interface{}) {
    if debugMode {
        log.Printf("[DEBUG] "+format, args...)
    }
}

// Usage throughout code
debugLog("User %s connected from %s", username, remoteAddr)
debugLog("Loading book from %s", bookPath)
debugLog("Saving progress for %s: chapter %d, offset %d", 
    username, chapter, offset)
```

#### Profiling

```go
// Add profiling endpoint
import _ "net/http/pprof"

func init() {
    if os.Getenv("ENABLE_PPROF") == "1" {
        go func() {
            log.Println("Starting pprof server on :6060")
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
}
```

Use profiling:
```bash
# Enable profiling
export ENABLE_PPROF=1
go run .

# Profile CPU usage
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Profile memory usage
go tool pprof http://localhost:6060/debug/pprof/heap

# Profile goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

#### Live Debugging with Delve

```bash
# Start with debugger
dlv debug . -- 

# Set breakpoints
(dlv) break main.teaHandler
(dlv) break book.go:LoadBook

# Continue execution
(dlv) continue

# Inspect variables
(dlv) print username
(dlv) print m.progress
```

## 🏗️ Build System

### Build Automation

The build system consists of several scripts:

#### build.sh
- Downloads dependencies
- Generates SSH keys if needed
- Builds optimized binary
- Sets up data directories

#### deploy.sh
- Production deployment
- Creates system user
- Installs systemd service
- Sets proper permissions

#### Docker Build
```dockerfile
# Multi-stage build for optimization
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o void-reader

FROM alpine:latest
RUN apk --no-cache add ca-certificates openssh-keygen
WORKDIR /app
COPY --from=builder /app/void-reader .
# ... rest of Dockerfile
```

### Cross-Compilation

```bash
# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o void-reader-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o void-reader-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o void-reader-windows-amd64.exe

# ARM builds
GOOS=linux GOARCH=arm64 go build -o void-reader-linux-arm64
GOOS=linux GOARCH=arm go build -o void-reader-linux-arm
```

### Release Process

```bash
#!/bin/bash
# release.sh

VERSION=${1:-"v1.0.0"}

echo "Building release $VERSION"

# Clean previous builds
rm -rf dist/
mkdir dist/

# Build for multiple platforms
platforms=("linux/amd64" "linux/arm64" "darwin/amd64" "windows/amd64")

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    
    output_name="void-reader-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    
    echo "Building $output_name..."
    
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-s -w -X main.Version=$VERSION" \
        -o dist/$output_name
        
    # Create archive
    if [ $GOOS = "windows" ]; then
        zip -j dist/void-reader-$GOOS-$GOARCH-$VERSION.zip dist/$output_name
    else
        tar -czf dist/void-reader-$GOOS-$GOARCH-$VERSION.tar.gz -C dist $output_name
    fi
done

echo "Release $VERSION built successfully"
ls -la dist/
```

## 🤝 Contributing

### Code Style

#### Go Conventions
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Comment exported functions and types
- Handle errors explicitly
- Use interfaces for testability

#### Project Conventions
```go
// File headers
// Package main provides the SSH reader server for Void Reavers book

// Function documentation
// LoadBook loads book content from the specified directory,
// trying Markdown format first, then falling back to LaTeX.
func LoadBook(bookDir string) (*Book, error) {
    // Implementation
}

// Error handling
if err != nil {
    return nil, fmt.Errorf("failed to load chapter %s: %w", filename, err)
}

// Struct documentation
// UserProgress tracks a user's reading progress including current position,
// completion status, bookmarks, and reading statistics.
type UserProgress struct {
    Username string `json:"username"`
    // ...
}
```

### Git Workflow

#### Branch Strategy
```bash
# Feature development
git checkout -b feature/search-functionality
git commit -am "Add basic search interface"
git push origin feature/search-functionality

# Bug fixes
git checkout -b fix/progress-saving-issue
git commit -am "Fix progress not saving on chapter change"
git push origin fix/progress-saving-issue

# Releases
git checkout -b release/v1.1.0
git tag v1.1.0
git push origin v1.1.0
```

#### Commit Messages
```
feat: add full-text search across all chapters
fix: resolve progress not saving on disconnection
docs: update configuration guide with new options
refactor: extract book loading logic into separate package
test: add integration tests for SSH connection handling
```

### Pull Request Process

1. **Fork and Clone**
   ```bash
   git clone https://github.com/yourusername/void-reavers-reader.git
   cd void-reavers-reader
   ```

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Develop and Test**
   ```bash
   # Make changes
   go test ./...
   go build
   ./void-reader  # Test manually
   ```

4. **Document Changes**
   - Update relevant documentation
   - Add inline code comments
   - Update CHANGELOG if applicable

5. **Submit Pull Request**
   - Clear description of changes
   - Reference any related issues
   - Include test results
   - Screenshots for UI changes

### Development Environment

#### Recommended Tools
- **Editor**: VS Code with Go extension
- **Debugger**: Delve (`dlv`)
- **Testing**: Go test runner
- **Linting**: `golangci-lint`
- **Formatting**: `gofmt`, `goimports`

#### VS Code Configuration
```json
// .vscode/settings.json
{
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v"],
    "go.buildFlags": ["-race"],
    "go.vetFlags": ["-all"]
}
```

#### Pre-commit Hooks
```bash
#!/bin/sh
# .git/hooks/pre-commit

# Format code
gofmt -w .
goimports -w .

# Run tests
go test ./...

# Lint code
golangci-lint run

# Check for security issues
gosec ./...
```

## 📚 Additional Resources

### Go Libraries Used

| Library | Purpose | Documentation |
|---------|---------|---------------|
| `github.com/charmbracelet/bubbletea` | TUI framework | [Docs](https://github.com/charmbracelet/bubbletea) |
| `github.com/charmbracelet/lipgloss` | Styling and layout | [Docs](https://github.com/charmbracelet/lipgloss) |
| `github.com/charmbracelet/wish` | SSH server middleware | [Docs](https://github.com/charmbracelet/wish) |
| `github.com/muesli/reflow` | Text processing utilities | [Docs](https://github.com/muesli/reflow) |

### Learning Resources

- **Bubbletea Tutorial**: [Official Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- **SSH in Go**: [Go SSH Documentation](https://pkg.go.dev/golang.org/x/crypto/ssh)
- **TUI Design**: [Charm Terminal UI Guidelines](https://charm.sh/)
- **Go Best Practices**: [Effective Go](https://golang.org/doc/effective_go.html)

### Architecture Inspiration

- **Glow**: Markdown reader by Charm
- **Soft Serve**: Git server with TUI
- **VHS**: Terminal session recorder
- **Slides**: Terminal presentation tool

---

**Happy Developing!** 🚀✨

This guide should give you everything needed to understand, modify, and contribute to the Void Reavers SSH Reader. The codebase is designed to be modular and extensible - feel free to experiment and improve it!

**Next Steps:**
- Set up your development environment
- Try adding a small feature or fix
- Run the tests and ensure everything works
- Submit your first contribution!