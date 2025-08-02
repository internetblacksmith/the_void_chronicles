# API Reference

Technical reference for the Void Reavers SSH Reader internal APIs, data structures, and extension points.

## 📋 Overview

This document covers:
- **Core Data Types**: Structures and interfaces
- **Public APIs**: Functions available for extension
- **Internal APIs**: Implementation details
- **Extension Points**: How to add functionality
- **Configuration**: Programmatic configuration options

## 🏗️ Core Data Types

### Book Management

#### Book Structure
```go
type Book struct {
    Title    string    `json:"title"`
    Author   string    `json:"author"`
    Chapters []Chapter `json:"chapters"`
}
```

**Fields:**
- `Title`: Book title displayed in UI
- `Author`: Author name for metadata
- `Chapters`: Ordered list of book chapters

**Methods:**
```go
func (b *Book) GetChapter(index int) (*Chapter, error)
func (b *Book) ChapterCount() int
func (b *Book) FindChapterByTitle(title string) (*Chapter, int, error)
```

#### Chapter Structure
```go
type Chapter struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}
```

**Fields:**
- `Title`: Chapter title for navigation
- `Content`: Full chapter text content

**Methods:**
```go
func (c *Chapter) WordCount() int
func (c *Chapter) EstimatedReadingTime() time.Duration
func (c *Chapter) Preview(maxLength int) string
```

### Progress Management

#### UserProgress Structure
```go
type UserProgress struct {
    Username       string            `json:"username"`
    CurrentChapter int               `json:"current_chapter"`
    ScrollOffset   int               `json:"scroll_offset"`
    LastRead       time.Time         `json:"last_read"`
    ChapterProgress map[int]bool     `json:"chapter_progress"`
    Bookmarks      []Bookmark        `json:"bookmarks"`
    ReadingTime    time.Duration     `json:"reading_time"`
    SessionStart   time.Time         `json:"-"`
}
```

**Fields:**
- `Username`: Unique user identifier
- `CurrentChapter`: Last read chapter (0-indexed)
- `ScrollOffset`: Position within current chapter
- `LastRead`: Timestamp of last reading session
- `ChapterProgress`: Completion status per chapter
- `Bookmarks`: Saved reading positions
- `ReadingTime`: Total accumulated reading time
- `SessionStart`: Current session start time (not persisted)

**Methods:**
```go
func (p *UserProgress) AddBookmark(chapter, scrollOffset int, note string)
func (p *UserProgress) RemoveBookmark(index int)
func (p *UserProgress) MarkChapterComplete(chapter int)
func (p *UserProgress) IsChapterComplete(chapter int) bool
func (p *UserProgress) GetCompletionPercentage(totalChapters int) float64
func (p *UserProgress) GetReadingStats() map[string]interface{}
```

#### Bookmark Structure
```go
type Bookmark struct {
    Chapter      int       `json:"chapter"`
    ScrollOffset int       `json:"scroll_offset"`
    Note         string    `json:"note"`
    Created      time.Time `json:"created"`
}
```

**Fields:**
- `Chapter`: Chapter index (0-based)
- `ScrollOffset`: Position within chapter
- `Note`: User annotation (future feature)
- `Created`: Bookmark creation timestamp

### UI State Management

#### Model Structure
```go
type model struct {
    state           state
    book            *Book
    menuCursor      int
    chapterCursor   int
    currentChapter  int
    scrollOffset    int
    width           int
    height          int
    menuItems       []string
    quitting        bool
    progress        *UserProgress
    progressManager *ProgressManager
    username        string
}
```

**State Enumeration:**
```go
type state int

const (
    menuView state = iota    // Main menu screen
    chapterListView         // Chapter selection screen
    readingView            // Reading interface
    aboutView              // Information screen
    progressView           // Progress statistics
)
```

**Bubbletea Interface Implementation:**
```go
func (m model) Init() tea.Cmd
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m model) View() string
```

## 🔧 Public APIs

### Book Loading API

#### LoadBook Function
```go
func LoadBook(bookDir string) (*Book, error)
```

**Parameters:**
- `bookDir`: Path to book directory

**Returns:**
- `*Book`: Loaded book structure
- `error`: Loading error if any

**Behavior:**
1. Attempts to load from `bookDir/markdown/` directory
2. Falls back to LaTeX files in `bookDir/`
3. Returns error if no valid content found

**Example Usage:**
```go
book, err := LoadBook("book1_void_reavers")
if err != nil {
    log.Fatalf("Failed to load book: %v", err)
}

fmt.Printf("Loaded: %s by %s (%d chapters)\n", 
    book.Title, book.Author, len(book.Chapters))
```

#### Format-Specific Loaders
```go
func loadFromMarkdown(bookDir string) (*Book, error)
func loadFromLaTeX(bookDir string) (*Book, error)
```

**Extension Point**: Add new format loaders:
```go
func loadFromEPUB(bookPath string) (*Book, error) {
    // Implementation for EPUB format
}

// Register in LoadBook function
func LoadBook(bookDir string) (*Book, error) {
    // Try EPUB first
    if epubFile := findEPUBFile(bookDir); epubFile != "" {
        if book, err := loadFromEPUB(epubFile); err == nil {
            return book, nil
        }
    }
    
    // Existing loaders...
}
```

### Progress Management API

#### ProgressManager Interface
```go
type ProgressManager struct {
    dataDir string
}

func NewProgressManager() *ProgressManager
func (pm *ProgressManager) LoadProgress(username string) (*UserProgress, error)
func (pm *ProgressManager) SaveProgress(progress *UserProgress) error
```

**Usage Example:**
```go
pm := NewProgressManager()

// Load user progress
progress, err := pm.LoadProgress("alice")
if err != nil {
    // Create new progress for new user
    progress = &UserProgress{
        Username: "alice",
        CurrentChapter: 0,
        ChapterProgress: make(map[int]bool),
        Bookmarks: []Bookmark{},
    }
}

// Update progress
progress.CurrentChapter = 5
progress.ScrollOffset = 150
progress.MarkChapterComplete(4)

// Save progress
err = pm.SaveProgress(progress)
```

#### Custom Storage Backends

**Interface for Storage Backends:**
```go
type ProgressStorage interface {
    Load(username string) (*UserProgress, error)
    Save(progress *UserProgress) error
    Delete(username string) error
    List() ([]string, error)
}
```

**File-based Implementation:**
```go
type FileProgressStorage struct {
    dataDir string
}

func (fps *FileProgressStorage) Load(username string) (*UserProgress, error) {
    filename := filepath.Join(fps.dataDir, username+".json")
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    var progress UserProgress
    err = json.Unmarshal(data, &progress)
    return &progress, err
}
```

**Redis Implementation Example:**
```go
type RedisProgressStorage struct {
    client *redis.Client
}

func (rps *RedisProgressStorage) Load(username string) (*UserProgress, error) {
    data, err := rps.client.Get(ctx, "progress:"+username).Result()
    if err != nil {
        return nil, err
    }
    
    var progress UserProgress
    err = json.Unmarshal([]byte(data), &progress)
    return &progress, err
}

func (rps *RedisProgressStorage) Save(progress *UserProgress) error {
    data, err := json.Marshal(progress)
    if err != nil {
        return err
    }
    
    return rps.client.Set(ctx, "progress:"+progress.Username, data, 0).Err()
}
```

### Text Processing API

#### Text Wrapping
```go
func wrapText(text string, width int) []string
```

**Parameters:**
- `text`: Input text to wrap
- `width`: Maximum line width

**Returns:**
- `[]string`: Array of wrapped lines

**Example:**
```go
lines := wrapText("This is a very long line that needs to be wrapped.", 20)
// Returns: ["This is a very long", "line that needs to be", "wrapped."]
```

#### Content Conversion
```go
func convertLaTeXToPlainText(content string) string
func parseMarkdownChapter(content string) Chapter
func parseLaTeXChapter(content string) Chapter
```

**Extension Point**: Add custom converters:
```go
func convertCustomFormat(content string) string {
    // Custom format conversion logic
    content = regexp.MustCompile(`\[ship:([^\]]+)\]`).
        ReplaceAllString(content, "*$1*")
    return content
}
```

## 🎨 UI Extension API

### Styling System

#### Style Definitions
```go
var (
    titleStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("86")).
        Bold(true).
        Align(lipgloss.Center)

    headerStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("39")).
        Bold(true).
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("63")).
        Padding(0, 1)

    selectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("0")).
        Background(lipgloss.Color("86")).
        Bold(true)

    normalStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("252"))
)
```

#### Theme System Extension
```go
type Theme struct {
    Primary     lipgloss.Color
    Secondary   lipgloss.Color
    Background  lipgloss.Color
    Text        lipgloss.Color
    Accent      lipgloss.Color
}

var themes = map[string]Theme{
    "default": {
        Primary:    lipgloss.Color("86"),
        Secondary:  lipgloss.Color("39"),
        Background: lipgloss.Color("0"),
        Text:       lipgloss.Color("252"),
        Accent:     lipgloss.Color("63"),
    },
    "dark": {
        Primary:    lipgloss.Color("15"),
        Secondary:  lipgloss.Color("244"),
        Background: lipgloss.Color("0"),
        Text:       lipgloss.Color("252"),
        Accent:     lipgloss.Color("236"),
    },
}

func ApplyTheme(themeName string) {
    theme := themes[themeName]
    
    titleStyle = titleStyle.Foreground(theme.Primary)
    headerStyle = headerStyle.Foreground(theme.Secondary)
    selectedStyle = selectedStyle.Background(theme.Primary)
    // ... apply to other styles
}
```

### Custom View Components

#### Component Interface
```go
type ViewComponent interface {
    Render(width, height int) string
    HandleInput(msg tea.KeyMsg) tea.Cmd
    Update(msg tea.Msg) tea.Cmd
}
```

#### Progress Bar Component
```go
type ProgressBar struct {
    width     int
    current   int
    total     int
    style     lipgloss.Style
}

func NewProgressBar(width int) *ProgressBar {
    return &ProgressBar{
        width: width,
        style: lipgloss.NewStyle().Foreground(lipgloss.Color("86")),
    }
}

func (pb *ProgressBar) SetProgress(current, total int) {
    pb.current = current
    pb.total = total
}

func (pb *ProgressBar) Render(width, height int) string {
    if pb.total == 0 {
        return ""
    }
    
    filled := int(float64(pb.current) / float64(pb.total) * float64(width))
    bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
    
    return pb.style.Render(bar)
}
```

## 🔌 Plugin System API

### Plugin Interface (Future Feature)
```go
type Plugin interface {
    Name() string
    Version() string
    Init(config map[string]interface{}) error
    Shutdown() error
}

type ContentPlugin interface {
    Plugin
    ProcessContent(content string) string
    SupportedFormats() []string
}

type UIPlugin interface {
    Plugin
    RegisterViews() map[string]ViewComponent
    RegisterCommands() map[string]func(args []string) error
}
```

### Plugin Manager
```go
type PluginManager struct {
    plugins map[string]Plugin
    config  map[string]map[string]interface{}
}

func NewPluginManager() *PluginManager {
    return &PluginManager{
        plugins: make(map[string]Plugin),
        config:  make(map[string]map[string]interface{}),
    }
}

func (pm *PluginManager) LoadPlugin(name string, plugin Plugin) error {
    config := pm.config[name]
    if config == nil {
        config = make(map[string]interface{})
    }
    
    err := plugin.Init(config)
    if err != nil {
        return fmt.Errorf("failed to initialize plugin %s: %w", name, err)
    }
    
    pm.plugins[name] = plugin
    return nil
}

func (pm *PluginManager) GetPlugin(name string) (Plugin, bool) {
    plugin, exists := pm.plugins[name]
    return plugin, exists
}
```

## 🌐 Network API

### SSH Server Configuration
```go
type ServerConfig struct {
    Host            string        `json:"host"`
    Port            string        `json:"port"`
    HostKeyPath     string        `json:"host_key_path"`
    MaxConnections  int           `json:"max_connections"`
    ReadTimeout     time.Duration `json:"read_timeout"`
    WriteTimeout    time.Duration `json:"write_timeout"`
    IdleTimeout     time.Duration `json:"idle_timeout"`
}

func NewServer(config ServerConfig) (*ssh.Server, error) {
    s, err := wish.NewServer(
        wish.WithAddress(net.JoinHostPort(config.Host, config.Port)),
        wish.WithHostKeyPath(config.HostKeyPath),
        wish.WithIdleTimeout(config.IdleTimeout),
        wish.WithMaxTimeout(config.ReadTimeout),
        wish.WithMiddleware(
            bubbletea.Middleware(teaHandler),
            logging.Middleware(),
        ),
    )
    
    return s, err
}
```

### Middleware API
```go
type Middleware func(ssh.Handler) ssh.Handler

func AuthMiddleware(authorizedKeys []ssh.PublicKey) Middleware {
    return func(next ssh.Handler) ssh.Handler {
        return func(s ssh.Session) {
            // Authentication logic
            if !isAuthorized(s.PublicKey(), authorizedKeys) {
                s.Exit(1)
                return
            }
            next(s)
        }
    }
}

func RateLimitMiddleware(limit rate.Limit) Middleware {
    limiter := rate.NewLimiter(limit, 1)
    
    return func(next ssh.Handler) ssh.Handler {
        return func(s ssh.Session) {
            if !limiter.Allow() {
                s.Write([]byte("Rate limit exceeded\n"))
                s.Exit(1)
                return
            }
            next(s)
        }
    }
}

func MetricsMiddleware() Middleware {
    return func(next ssh.Handler) ssh.Handler {
        return func(s ssh.Session) {
            start := time.Now()
            connectionsTotal.WithLabelValues("started").Inc()
            activeConnections.Inc()
            
            defer func() {
                duration := time.Since(start)
                activeConnections.Dec()
                connectionDuration.Observe(duration.Seconds())
            }()
            
            next(s)
        }
    }
}
```

## 📊 Metrics and Monitoring API

### Prometheus Metrics
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    connectionsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "void_reader_connections_total",
            Help: "Total number of SSH connections",
        },
        []string{"status"},
    )
    
    activeConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "void_reader_active_connections",
            Help: "Number of currently active connections",
        },
    )
    
    chapterReads = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "void_reader_chapter_reads_total",
            Help: "Total number of chapter reads",
        },
        []string{"chapter"},
    )
    
    readingTime = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name: "void_reader_session_duration_seconds",
            Help: "Time spent in reading sessions",
            Buckets: prometheus.DefBuckets,
        },
    )
)

func init() {
    prometheus.MustRegister(connectionsTotal)
    prometheus.MustRegister(activeConnections)
    prometheus.MustRegister(chapterReads)
    prometheus.MustRegister(readingTime)
}
```

### Health Check API
```go
type HealthChecker struct {
    server          *ssh.Server
    progressManager *ProgressManager
    bookCache       *BookCache
}

type HealthStatus struct {
    Status      string                 `json:"status"`
    Timestamp   time.Time             `json:"timestamp"`
    Version     string                `json:"version"`
    Uptime      time.Duration         `json:"uptime"`
    Connections int                   `json:"active_connections"`
    Memory      MemoryStats           `json:"memory"`
    Checks      map[string]CheckResult `json:"checks"`
}

type CheckResult struct {
    Status  string        `json:"status"`
    Message string        `json:"message,omitempty"`
    Latency time.Duration `json:"latency"`
}

func (hc *HealthChecker) GetHealth() HealthStatus {
    status := HealthStatus{
        Status:      "healthy",
        Timestamp:   time.Now(),
        Version:     Version,
        Uptime:      time.Since(startTime),
        Connections: getActiveConnections(),
        Memory:      getMemoryStats(),
        Checks:      make(map[string]CheckResult),
    }
    
    // Check individual components
    status.Checks["ssh_server"] = hc.checkSSHServer()
    status.Checks["progress_storage"] = hc.checkProgressStorage()
    status.Checks["book_cache"] = hc.checkBookCache()
    
    // Determine overall status
    for _, check := range status.Checks {
        if check.Status != "healthy" {
            status.Status = "unhealthy"
            break
        }
    }
    
    return status
}
```

## 🔧 Configuration API

### Configuration Structure
```go
type Config struct {
    Server   ServerConfig   `json:"server"`
    Books    BooksConfig    `json:"books"`
    UI       UIConfig       `json:"ui"`
    Storage  StorageConfig  `json:"storage"`
    Security SecurityConfig `json:"security"`
    Logging  LoggingConfig  `json:"logging"`
}

type ServerConfig struct {
    Host           string        `json:"host"`
    Port           string        `json:"port"`
    MaxConnections int           `json:"max_connections"`
    ReadTimeout    time.Duration `json:"read_timeout"`
    WriteTimeout   time.Duration `json:"write_timeout"`
}

type BooksConfig struct {
    DefaultBook   string   `json:"default_book"`
    Directory     string   `json:"directory"`
    AutoDetect    bool     `json:"auto_detect"`
    SupportedFormats []string `json:"supported_formats"`
}

type UIConfig struct {
    Theme         string `json:"theme"`
    Animations    bool   `json:"animations"`
    CompactMode   bool   `json:"compact_mode"`
    ShowProgress  bool   `json:"show_progress"`
}
```

### Configuration Loading
```go
func LoadConfig(configPath string) (*Config, error) {
    // Default configuration
    config := &Config{
        Server: ServerConfig{
            Host:           "localhost",
            Port:           "23234",
            MaxConnections: 50,
            ReadTimeout:    30 * time.Second,
            WriteTimeout:   30 * time.Second,
        },
        Books: BooksConfig{
            DefaultBook:      "book1_void_reavers",
            Directory:        "./books",
            AutoDetect:       true,
            SupportedFormats: []string{"markdown", "latex", "epub"},
        },
        // ... other defaults
    }
    
    // Load from file if exists
    if _, err := os.Stat(configPath); err == nil {
        data, err := os.ReadFile(configPath)
        if err != nil {
            return nil, fmt.Errorf("failed to read config file: %w", err)
        }
        
        if err := json.Unmarshal(data, config); err != nil {
            return nil, fmt.Errorf("failed to parse config file: %w", err)
        }
    }
    
    // Override with environment variables
    if host := os.Getenv("VOID_READER_HOST"); host != "" {
        config.Server.Host = host
    }
    if port := os.Getenv("VOID_READER_PORT"); port != "" {
        config.Server.Port = port
    }
    
    return config, nil
}
```

## 🧪 Testing API

### Test Utilities
```go
// Test helper functions
func CreateTestBook(tempDir string) (*Book, error) {
    bookDir := filepath.Join(tempDir, "test_book")
    os.MkdirAll(filepath.Join(bookDir, "markdown"), 0755)
    
    chapters := []struct {
        filename string
        title    string
        content  string
    }{
        {"chapter01.md", "Test Chapter 1", "This is test content for chapter 1."},
        {"chapter02.md", "Test Chapter 2", "This is test content for chapter 2."},
    }
    
    for _, ch := range chapters {
        content := fmt.Sprintf("# %s\n\n%s", ch.title, ch.content)
        err := os.WriteFile(
            filepath.Join(bookDir, "markdown", ch.filename),
            []byte(content),
            0644,
        )
        if err != nil {
            return nil, err
        }
    }
    
    return LoadBook(bookDir)
}

func CreateTestProgress(username string) *UserProgress {
    return &UserProgress{
        Username:        username,
        CurrentChapter:  1,
        ScrollOffset:    50,
        LastRead:        time.Now(),
        ChapterProgress: map[int]bool{0: true},
        Bookmarks: []Bookmark{
            {
                Chapter:      1,
                ScrollOffset: 25,
                Note:         "Test bookmark",
                Created:      time.Now(),
            },
        },
        ReadingTime: 30 * time.Minute,
    }
}
```

### Mock Implementations
```go
type MockProgressManager struct {
    data map[string]*UserProgress
    mu   sync.RWMutex
}

func NewMockProgressManager() *MockProgressManager {
    return &MockProgressManager{
        data: make(map[string]*UserProgress),
    }
}

func (mpm *MockProgressManager) LoadProgress(username string) (*UserProgress, error) {
    mpm.mu.RLock()
    defer mpm.mu.RUnlock()
    
    if progress, exists := mpm.data[username]; exists {
        return progress, nil
    }
    
    return nil, fmt.Errorf("progress not found for user: %s", username)
}

func (mpm *MockProgressManager) SaveProgress(progress *UserProgress) error {
    mpm.mu.Lock()
    defer mpm.mu.Unlock()
    
    mpm.data[progress.Username] = progress
    return nil
}
```

## 🔍 Error Handling

### Error Types
```go
type BookError struct {
    Path    string
    Message string
    Cause   error
}

func (e *BookError) Error() string {
    return fmt.Sprintf("book error in %s: %s", e.Path, e.Message)
}

func (e *BookError) Unwrap() error {
    return e.Cause
}

type ProgressError struct {
    Username string
    Message  string
    Cause    error
}

func (e *ProgressError) Error() string {
    return fmt.Sprintf("progress error for %s: %s", e.Username, e.Message)
}

func (e *ProgressError) Unwrap() error {
    return e.Cause
}
```

### Error Handling Patterns
```go
// Wrap errors with context
func LoadBook(bookDir string) (*Book, error) {
    book, err := loadFromMarkdown(bookDir)
    if err != nil {
        return nil, &BookError{
            Path:    bookDir,
            Message: "failed to load from markdown",
            Cause:   err,
        }
    }
    return book, nil
}

// Check error types
func handleError(err error) {
    var bookErr *BookError
    var progressErr *ProgressError
    
    switch {
    case errors.As(err, &bookErr):
        log.Printf("Book loading failed: %v", bookErr)
    case errors.As(err, &progressErr):
        log.Printf("Progress saving failed: %v", progressErr)
    default:
        log.Printf("Unknown error: %v", err)
    }
}
```

---

**API Reference Complete!** 📚✨

This reference provides comprehensive documentation for all public APIs and extension points in the Void Reavers SSH Reader. Use this as a guide for extending functionality, integrating with other systems, or contributing to the project.

**Related Documentation:**
- [Development Guide](development.md) - For implementation details
- [Configuration Guide](configuration.md) - For configuration options
- [User Guide](user-guide.md) - For end-user functionality