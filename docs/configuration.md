# Configuration Guide

Advanced configuration options for the Void Reavers SSH Reader to customize your reading experience and server setup.

## 📋 Configuration Overview

The SSH Reader can be configured at multiple levels:
- **Application Level**: Server settings, ports, and behavior
- **User Level**: Individual reading preferences and progress
- **System Level**: Security, logging, and deployment options
- **Book Level**: Content sources and formatting

## ⚙️ Application Configuration

### Server Settings

#### Network Configuration

Edit `main.go` to modify basic network settings:

```go
const (
    host = "localhost"  // Server bind address
    port = "23234"      // Server port
)
```

**Common Configurations:**

**Local Only (Default):**
```go
const (
    host = "localhost"  // Only local connections
    port = "23234"
)
```

**Allow External Connections:**
```go
const (
    host = "0.0.0.0"    // All interfaces
    port = "23234"
)
```

**Custom Port:**
```go
const (
    host = "localhost"
    port = "2222"       // Custom port
)
```

**IPv6 Support:**
```go
const (
    host = "::"         // IPv6 all interfaces
    port = "23234"
)
```

#### SSH Configuration

##### Host Key Management

**Default Behavior:**
- Generates Ed25519 key automatically
- Stored in `.ssh/id_ed25519`
- 256-bit security level

**Custom Host Key:**
```bash
# Generate custom key
ssh-keygen -t rsa -b 4096 -f .ssh/custom_host_key -N ""

# Update main.go
wish.WithHostKeyPath(".ssh/custom_host_key")
```

**Multiple Host Keys:**
```go
// In teaHandler function, add multiple keys
wish.WithHostKeyPath(".ssh/id_ed25519"),
wish.WithHostKeyPath(".ssh/id_rsa"),
```

##### SSH Options

Add custom SSH middleware in `main.go`:

```go
wish.WithMiddleware(
    bubbletea.Middleware(teaHandler),
    logging.Middleware(),
    // Add custom middleware here
),
```

### Application Behavior

#### Reading Settings

**Auto-save Frequency:**
Modify the reading update function to save more/less frequently:

```go
// In updateReading function, add periodic saves
case "any_key":
    // Save progress every 10 scrolls
    if m.scrollOffset%10 == 0 {
        m.progress.CurrentChapter = m.currentChapter
        m.progress.ScrollOffset = m.scrollOffset
        m.progressManager.SaveProgress(m.progress)
    }
```

**Chapter Completion Threshold:**
Set when chapters are marked complete:

```go
// Mark complete when reaching 90% of chapter
func (m model) checkChapterCompletion() {
    maxScroll := m.getMaxScroll()
    if maxScroll > 0 && float64(m.scrollOffset)/float64(maxScroll) >= 0.9 {
        m.progress.MarkChapterComplete(m.currentChapter)
    }
}
```

#### Menu Customization

**Menu Items:**
Modify the menu items in `initialModel()`:

```go
menuItems: []string{
    "📖 Continue Reading",
    "📚 Chapter List", 
    "📊 Progress",
    "🔍 Search",        // New item
    "⚙️  Settings",     // New item
    "ℹ️  About",
    "🚪 Exit"
},
```

**Custom Menu Actions:**
Add new cases in `updateMenu()`:

```go
case 3: // Search
    m.state = searchView
case 4: // Settings  
    m.state = settingsView
```

## 🎨 User Interface Configuration

### Styling and Themes

#### Color Schemes

Edit the lipgloss styles in `main.go`:

**Default Theme:**
```go
titleStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("86")).  // Cyan
    Bold(true)

selectedStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("0")).   // Black on
    Background(lipgloss.Color("86"))   // Cyan background
```

**Dark Theme:**
```go
titleStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("15")).  // White
    Bold(true)

selectedStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("15")).  // White on
    Background(lipgloss.Color("236"))  // Dark gray
```

**High Contrast Theme:**
```go
titleStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("15")).  // White
    Bold(true).
    Background(lipgloss.Color("0"))    // Black background

selectedStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("0")).   // Black on
    Background(lipgloss.Color("15"))   // White background
```

#### Layout Configuration

**Reading Area Padding:**
```go
contentStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("252")).
    Padding(2, 4)  // Increase padding: top/bottom, left/right
```

**Border Styles:**
```go
headerStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).     // Rounded corners
    BorderForeground(lipgloss.Color("63"))
```

#### Typography

**Font Weight:**
```go
titleStyle = lipgloss.NewStyle().
    Bold(true).        // Bold
    Italic(false).     // Not italic
    Underline(false)   // Not underlined
```

**Text Alignment:**
```go
titleStyle = lipgloss.NewStyle().
    Align(lipgloss.Center)    // Center, Left, or Right
```

### Terminal Compatibility

#### Color Support Detection

Add automatic color detection:

```go
func detectColorSupport() int {
    term := os.Getenv("TERM")
    colorterm := os.Getenv("COLORTERM")
    
    if colorterm == "truecolor" || colorterm == "24bit" {
        return 16777216 // 24-bit color
    }
    if strings.Contains(term, "256color") {
        return 256      // 256 colors
    }
    return 16           // 16 colors
}
```

#### Responsive Layout

```go
func (m model) adaptToTerminalSize() {
    if m.width < 80 {
        // Compact mode for narrow terminals
        m.menuItems = []string{"📖 Read", "📚 List", "📊 Stats", "🚪 Exit"}
    } else {
        // Full mode for wide terminals  
        m.menuItems = []string{
            "📖 Continue Reading",
            "📚 Chapter List", 
            "📊 Progress",
            "ℹ️  About",
            "🚪 Exit"
        }
    }
}
```

## 📚 Book Content Configuration

### Multiple Book Support

#### Book Directory Structure

Organize multiple books:
```
books/
├── book1_void_reavers/
│   ├── markdown/
│   └── *.tex
├── book2_shadow_dancers/
│   ├── markdown/
│   └── *.tex
└── book3_quantum_academy/
    ├── markdown/
    └── *.tex
```

#### Book Selection Menu

Add book selection to the main menu:

```go
type model struct {
    // ... existing fields
    books        []BookInfo
    selectedBook int
}

type BookInfo struct {
    Path        string
    Title       string
    Author      string
    Description string
}

func loadAvailableBooks() []BookInfo {
    var books []BookInfo
    
    dirs, _ := filepath.Glob("book*")
    for _, dir := range dirs {
        if info, err := loadBookInfo(dir); err == nil {
            books = append(books, info)
        }
    }
    
    return books
}
```

### Content Format Support

#### Markdown Configuration

**Custom Markdown Parser:**
```go
func parseMarkdownChapter(content string) Chapter {
    // Add support for custom markdown extensions
    content = convertCustomMarkdown(content)
    
    // Process standard markdown
    return processStandardMarkdown(content)
}

func convertCustomMarkdown(content string) string {
    // Custom formatting rules
    content = regexp.MustCompile(`\[ship:([^\]]+)\]`).
        ReplaceAllString(content, "*$1*")  // Ship names in italics
    
    content = regexp.MustCompile(`\[comm:([^\]]+)\]`).
        ReplaceAllString(content, `"$1"`)  // Communications in quotes
        
    return content
}
```

#### LaTeX Configuration

**Custom LaTeX Commands:**
```go
func convertLaTeXToPlainText(content string) string {
    // Standard conversions...
    
    // Custom command support
    customCommands := map[string]string{
        `\\ship{([^}]+)}`:     "*$1*",      // Ship names
        `\\comm{([^}]+)}`:     `"$1"`,      // Communications  
        `\\location{([^}]+)}`: "[$1]",      // Locations
        `\\tech{([^}]+)}`:     "[$1]",      // Technology terms
    }
    
    for pattern, replacement := range customCommands {
        re := regexp.MustCompile(pattern)
        content = re.ReplaceAllString(content, replacement)
    }
    
    return content
}
```

### Content Preprocessing

#### Text Processing Pipeline

```go
type ContentProcessor struct {
    filters []ContentFilter
}

type ContentFilter func(string) string

func NewContentProcessor() *ContentProcessor {
    return &ContentProcessor{
        filters: []ContentFilter{
            normalizeWhitespace,
            convertQuotes,
            formatDialogue,
            processEmphasis,
        },
    }
}

func (cp *ContentProcessor) Process(content string) string {
    for _, filter := range cp.filters {
        content = filter(content)
    }
    return content
}
```

## 💾 Data Storage Configuration

### Progress Storage

#### Custom Storage Location

```go
func NewProgressManager(dataDir string) *ProgressManager {
    if dataDir == "" {
        dataDir = ".void_reader_data"
    }
    
    // Ensure directory exists
    os.MkdirAll(dataDir, 0755)
    
    return &ProgressManager{dataDir: dataDir}
}
```

#### Database Backend (Advanced)

For high-volume deployments, replace JSON with database:

```go
import "database/sql"
import _ "github.com/mattn/go-sqlite3"

type DatabaseProgressManager struct {
    db *sql.DB
}

func NewDatabaseProgressManager(dbPath string) *DatabaseProgressManager {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create tables
    createTables(db)
    
    return &DatabaseProgressManager{db: db}
}
```

### User Data Encryption

```go
import "crypto/aes"
import "crypto/cipher"

func (pm *ProgressManager) SaveProgressEncrypted(progress *UserProgress, key []byte) error {
    data, err := json.Marshal(progress)
    if err != nil {
        return err
    }
    
    // Encrypt data
    encrypted, err := encrypt(data, key)
    if err != nil {
        return err
    }
    
    filename := filepath.Join(pm.dataDir, progress.Username+".enc")
    return os.WriteFile(filename, encrypted, 0644)
}
```

## 🔐 Security Configuration

### Authentication

#### SSH Key Authentication

```go
// Add public key authentication
wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
    // Load authorized keys
    authorizedKeys := loadAuthorizedKeys(".ssh/authorized_keys")
    
    // Check if key is authorized
    return isKeyAuthorized(key, authorizedKeys)
}),
```

#### User Identification

```go
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
    // Get user from SSH session
    username := s.User()
    if username == "" {
        username = "anonymous"
    }
    
    // Initialize model with user context
    m := initialModelWithUser(pty.Window.Width, pty.Window.Height, username)
    return m, []tea.ProgramOption{tea.WithAltScreen()}
}
```

### Access Control

#### Rate Limiting

```go
import "golang.org/x/time/rate"

type RateLimiter struct {
    limiter *rate.Limiter
}

func NewRateLimiter() *RateLimiter {
    // Allow 10 connections per minute
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Every(time.Minute/10), 1),
    }
}

func (rl *RateLimiter) Allow() bool {
    return rl.limiter.Allow()
}
```

#### Connection Limits

```go
var activeConnections int32
var maxConnections int32 = 50

func checkConnectionLimit() bool {
    current := atomic.LoadInt32(&activeConnections)
    return current < maxConnections
}
```

### Data Protection

#### File Permissions

```bash
# Secure the data directory
chmod 700 .void_reader_data
chmod 600 .void_reader_data/*.json

# Secure SSH keys
chmod 600 .ssh/id_ed25519
chmod 644 .ssh/id_ed25519.pub
```

#### Backup Configuration

```go
func (pm *ProgressManager) BackupUserData() error {
    backupDir := filepath.Join(pm.dataDir, "backups")
    os.MkdirAll(backupDir, 0755)
    
    timestamp := time.Now().Format("2006-01-02_15-04-05")
    backupPath := filepath.Join(backupDir, "backup_"+timestamp+".tar.gz")
    
    return createTarGzBackup(pm.dataDir, backupPath)
}
```

## 🚀 Performance Configuration

### Memory Management

#### Buffer Settings

```go
const (
    MaxChapterSize  = 1024 * 1024    // 1MB per chapter
    MaxBookSize     = 50 * 1024 * 1024 // 50MB per book
    UserDataCache   = 100             // Cache 100 users
)
```

#### Connection Pooling

```go
type ConnectionPool struct {
    connections chan *ssh.Session
    maxSize     int
}

func NewConnectionPool(maxSize int) *ConnectionPool {
    return &ConnectionPool{
        connections: make(chan *ssh.Session, maxSize),
        maxSize:     maxSize,
    }
}
```

### Logging Configuration

#### Log Levels

```go
import "github.com/charmbracelet/log"

func setupLogging() {
    log.SetLevel(log.InfoLevel)  // Debug, Info, Warn, Error
    log.SetReportCaller(true)
    log.SetTimeFormat("2006-01-02 15:04:05")
}
```

#### Custom Log Output

```go
func setupCustomLogging() {
    logFile, err := os.OpenFile("void-reader.log", 
        os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }
    
    log.SetOutput(logFile)
}
```

## 🔧 Environment Variables

### Configuration via Environment

```go
func loadConfig() Config {
    return Config{
        Host:        getEnv("VOID_READER_HOST", "localhost"),
        Port:        getEnv("VOID_READER_PORT", "23234"),
        DataDir:     getEnv("VOID_READER_DATA", ".void_reader_data"),
        LogLevel:    getEnv("VOID_READER_LOG_LEVEL", "info"),
        MaxUsers:    parseInt(getEnv("VOID_READER_MAX_USERS", "50")),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### Environment Variables Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `VOID_READER_HOST` | `localhost` | Server bind address |
| `VOID_READER_PORT` | `23234` | Server port |
| `VOID_READER_DATA` | `.void_reader_data` | User data directory |
| `VOID_READER_LOG_LEVEL` | `info` | Logging level |
| `VOID_READER_MAX_USERS` | `50` | Maximum concurrent users |
| `VOID_READER_BOOK_DIR` | `book1_void_reavers` | Default book directory |
| `VOID_READER_SSH_KEY` | `.ssh/id_ed25519` | SSH host key path |

## 📝 Configuration Files

### JSON Configuration

Create `config.json`:

```json
{
  "server": {
    "host": "localhost",
    "port": "23234",
    "max_connections": 50
  },
  "books": {
    "default": "book1_void_reavers",
    "directory": "./books",
    "auto_detect": true
  },
  "ui": {
    "theme": "default",
    "animations": true,
    "compact_mode": false
  },
  "storage": {
    "data_directory": ".void_reader_data",
    "backup_enabled": true,
    "encryption": false
  },
  "security": {
    "require_auth": false,
    "rate_limit": true,
    "log_connections": true
  }
}
```

### YAML Configuration

Create `config.yaml`:

```yaml
server:
  host: localhost
  port: 23234
  max_connections: 50

books:
  default: book1_void_reavers
  directory: ./books
  auto_detect: true

ui:
  theme: default
  animations: true
  compact_mode: false

storage:
  data_directory: .void_reader_data
  backup_enabled: true
  encryption: false

security:
  require_auth: false
  rate_limit: true
  log_connections: true
```

## 🔄 Advanced Customization

### Plugin System (Future)

```go
type Plugin interface {
    Name() string
    Init(config map[string]interface{}) error
    ProcessContent(content string) string
    HandleCommand(cmd string, args []string) error
}

type PluginManager struct {
    plugins map[string]Plugin
}

func (pm *PluginManager) LoadPlugin(name string, plugin Plugin) error {
    pm.plugins[name] = plugin
    return plugin.Init(nil)
}
```

### Custom Commands

```go
func (m model) handleCustomCommand(cmd string) (model, tea.Cmd) {
    switch cmd {
    case "stats":
        return m.showDetailedStats(), nil
    case "search":
        return m.enterSearchMode(), nil
    case "bookmark-list":
        return m.showBookmarkList(), nil
    default:
        return m, nil
    }
}
```

### Integration Hooks

```go
type EventHook func(event string, data interface{})

type HookManager struct {
    hooks map[string][]EventHook
}

func (hm *HookManager) RegisterHook(event string, hook EventHook) {
    hm.hooks[event] = append(hm.hooks[event], hook)
}

func (hm *HookManager) TriggerHooks(event string, data interface{}) {
    for _, hook := range hm.hooks[event] {
        hook(event, data)
    }
}
```

---

**Configuration Complete!** ⚙️✨

This guide covers most configuration scenarios. For more advanced customization, see the [Development Guide](development.md) to modify the source code directly.

**Next Steps:**
- Ready to deploy? See [Deployment Guide](deployment.md)
- Need help with issues? Check [Troubleshooting Guide](troubleshooting.md)
- Want to contribute? Read [Development Guide](development.md)