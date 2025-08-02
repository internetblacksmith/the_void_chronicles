# Contributing Guide

Welcome to the Void Reavers SSH Reader project! We're excited that you're interested in contributing. This guide will help you get started with contributing code, documentation, bug reports, and feature requests.

## 🎯 Ways to Contribute

### Code Contributions
- **Bug fixes**: Help resolve issues and improve stability
- **New features**: Add functionality like search, themes, or new book formats
- **Performance improvements**: Optimize loading, rendering, or memory usage
- **Refactoring**: Improve code structure and maintainability

### Documentation Contributions
- **User guides**: Improve installation and usage instructions
- **API documentation**: Document functions and data structures
- **Tutorials**: Create step-by-step guides for common tasks
- **Examples**: Provide code examples and use cases

### Community Contributions
- **Bug reports**: Help identify and document issues
- **Feature requests**: Suggest new functionality
- **Testing**: Help test new releases and features
- **Support**: Help other users in discussions

## 🚀 Getting Started

### Prerequisites

Before contributing, ensure you have:
- **Go 1.21+** installed
- **Git** for version control
- **SSH client** for testing
- **Text editor** or IDE with Go support

### Setting Up Development Environment

1. **Fork the Repository**
   ```bash
   # Go to GitHub and fork the repository
   # Then clone your fork
   git clone https://github.com/YOUR_USERNAME/void-reavers-reader.git
   cd void-reavers-reader
   ```

2. **Add Upstream Remote**
   ```bash
   git remote add upstream https://github.com/ORIGINAL_OWNER/void-reavers-reader.git
   ```

3. **Install Dependencies**
   ```bash
   go mod tidy
   ```

4. **Verify Setup**
   ```bash
   # Build the project
   go build -o void-reader-dev
   
   # Run tests
   go test ./...
   
   # Start development server
   ./void-reader-dev
   ```

5. **Test SSH Connection**
   ```bash
   # In another terminal
   ssh localhost -p 23234
   ```

### Development Workflow

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Write code following project conventions
   - Add tests for new functionality
   - Update documentation as needed

3. **Test Changes**
   ```bash
   # Run tests
   go test ./...
   
   # Manual testing
   go run . &
   ssh localhost -p 23234
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   # Create pull request on GitHub
   ```

## 📝 Code Style Guidelines

### Go Code Standards

#### Formatting
- Use `gofmt` for code formatting
- Use `goimports` for import organization
- Maximum line length: 100 characters
- Use tabs for indentation

#### Naming Conventions
```go
// Exported functions and types use PascalCase
func LoadBook(bookDir string) (*Book, error) {}
type UserProgress struct {}

// Unexported functions and variables use camelCase
func parseChapter(content string) Chapter {}
var currentUser string

// Constants use ALL_CAPS for package-level constants
const MAX_CONNECTIONS = 100

// Interface names end with -er when possible
type BookLoader interface {}
```

#### Code Organization
```go
// Package documentation
// Package main provides the SSH reader server for Void Reavers book.
package main

// Imports grouped and sorted
import (
    "fmt"
    "log"
    "os"
    
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Constants first
const (
    defaultPort = "23234"
    maxUsers    = 50
)

// Types second
type Model struct {
    // Public fields first
    State string
    Book  *Book
    
    // Private fields after
    cursor int
    width  int
}

// Functions last, with public functions before private
func NewModel() *Model {}
func (m *Model) Update() {}

func parseInput() {}
```

#### Error Handling
```go
// Always handle errors explicitly
book, err := LoadBook(bookDir)
if err != nil {
    return fmt.Errorf("failed to load book from %s: %w", bookDir, err)
}

// Use descriptive error messages
if len(chapters) == 0 {
    return errors.New("no chapters found in book directory")
}

// Wrap errors with context
if err := saveProgress(progress); err != nil {
    return fmt.Errorf("save progress for user %s: %w", username, err)
}
```

#### Documentation
```go
// Package-level documentation
// Package main provides the SSH reader server for Void Reavers book.

// Exported function documentation
// LoadBook loads book content from the specified directory,
// trying Markdown format first, then falling back to LaTeX.
// It returns an error if no valid book content is found.
func LoadBook(bookDir string) (*Book, error) {
    // Implementation
}

// Type documentation
// UserProgress tracks a user's reading progress including current position,
// completion status, bookmarks, and reading statistics.
type UserProgress struct {
    Username string `json:"username"` // Unique user identifier
    // ... other fields
}
```

### Project-Specific Conventions

#### File Organization
- **main.go**: SSH server, TUI logic, main application
- **book.go**: Book loading, parsing, content management
- **progress.go**: User progress, bookmarks, statistics
- **_test.go**: Test files alongside implementation

#### Import Organization
```go
import (
    // Standard library first
    "fmt"
    "log"
    "os"
    "time"
    
    // Third-party packages
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/wish"
    
    // Local packages (if any)
    "./internal/utils"
)
```

#### Variable Naming
```go
// UI-related variables
var (
    titleStyle    = lipgloss.NewStyle()    // UI styles
    menuCursor    int                      // UI state
    currentChapter int                     // Content state
)

// Business logic variables
var (
    progressManager *ProgressManager       // Core functionality
    bookCache      map[string]*Book       // Data management
)
```

## 🧪 Testing Guidelines

### Unit Tests

#### Test File Structure
```go
// book_test.go
package main

import (
    "os"
    "path/filepath"
    "testing"
    "time"
)

func TestLoadBook(t *testing.T) {
    // Test setup
    tempDir := t.TempDir()
    bookDir := createTestBook(t, tempDir)
    
    // Test execution
    book, err := LoadBook(bookDir)
    
    // Assertions
    if err != nil {
        t.Fatalf("LoadBook failed: %v", err)
    }
    
    if book.Title != "Test Book" {
        t.Errorf("Expected title 'Test Book', got '%s'", book.Title)
    }
    
    if len(book.Chapters) != 2 {
        t.Errorf("Expected 2 chapters, got %d", len(book.Chapters))
    }
}

// Test helper functions
func createTestBook(t *testing.T, tempDir string) string {
    bookDir := filepath.Join(tempDir, "test_book")
    os.MkdirAll(filepath.Join(bookDir, "markdown"), 0755)
    
    // Create test chapters
    chapters := []struct {
        filename string
        content  string
    }{
        {"chapter01.md", "# Chapter 1\n\nTest content 1"},
        {"chapter02.md", "# Chapter 2\n\nTest content 2"},
    }
    
    for _, ch := range chapters {
        err := os.WriteFile(
            filepath.Join(bookDir, "markdown", ch.filename),
            []byte(ch.content),
            0644,
        )
        if err != nil {
            t.Fatal(err)
        }
    }
    
    return bookDir
}
```

#### Test Naming
```go
func TestLoadBook(t *testing.T)                    // Basic functionality
func TestLoadBookFromMarkdown(t *testing.T)       // Specific format
func TestLoadBookWithMissingFiles(t *testing.T)   // Error cases
func TestLoadBookFromEmptyDirectory(t *testing.T) // Edge cases
```

#### Table-Driven Tests
```go
func TestTextWrapping(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        width    int
        expected []string
    }{
        {
            name:     "basic wrapping",
            input:    "This is a test line",
            width:    10,
            expected: []string{"This is a", "test line"},
        },
        {
            name:     "no wrapping needed",
            input:    "Short",
            width:    10,
            expected: []string{"Short"},
        },
        {
            name:     "empty input",
            input:    "",
            width:    10,
            expected: []string{},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := wrapText(tt.input, tt.width)
            if !slicesEqual(result, tt.expected) {
                t.Errorf("wrapText() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Integration Tests

#### SSH Connection Tests
```go
func TestSSHConnection(t *testing.T) {
    // Skip if SSH not available
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Start server in background
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    go func() {
        startServer(ctx)
    }()
    
    // Wait for server to start
    time.Sleep(2 * time.Second)
    
    // Test connection
    conn, err := net.Dial("tcp", "localhost:23234")
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    
    // Verify SSH handshake
    config := &ssh.ClientConfig{
        User: "testuser",
        Auth: []ssh.AuthMethod{},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout: 5 * time.Second,
    }
    
    sshConn, chans, reqs, err := ssh.NewClientConn(conn, "localhost:23234", config)
    if err != nil {
        t.Fatalf("SSH handshake failed: %v", err)
    }
    defer sshConn.Close()
    
    client := ssh.NewClient(sshConn, chans, reqs)
    defer client.Close()
    
    // Test session creation
    session, err := client.NewSession()
    if err != nil {
        t.Fatalf("Failed to create session: %v", err)
    }
    defer session.Close()
}
```

### Benchmark Tests
```go
func BenchmarkLoadBook(b *testing.B) {
    tempDir := b.TempDir()
    bookDir := createLargeTestBook(b, tempDir)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := LoadBook(bookDir)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkTextWrapping(b *testing.B) {
    longText := strings.Repeat("This is a long line that needs wrapping. ", 1000)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wrapText(longText, 80)
    }
}
```

## 📋 Pull Request Process

### Before Submitting

1. **Sync with Upstream**
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   git push origin main
   ```

2. **Rebase Feature Branch**
   ```bash
   git checkout feature/your-feature
   git rebase main
   ```

3. **Run Tests**
   ```bash
   go test ./...
   go vet ./...
   golangci-lint run
   ```

4. **Build and Manual Test**
   ```bash
   go build -o void-reader-test
   ./void-reader-test &
   ssh localhost -p 23234
   # Test your changes manually
   ```

### PR Description Template

```markdown
## Description
Brief description of changes and motivation.

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Changes Made
- List specific changes made
- Include any new files or modified files
- Mention any dependencies added/removed

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed
- [ ] All tests pass

## Screenshots (if applicable)
Include screenshots for UI changes.

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added for new functionality
- [ ] No breaking changes or breaking changes documented
```

### Review Process

1. **Automated Checks**
   - Tests must pass
   - Linting must pass
   - Build must succeed

2. **Code Review**
   - At least one maintainer review required
   - Address all feedback before merge
   - Maintain discussion in PR comments

3. **Final Steps**
   - Squash commits if requested
   - Update documentation
   - Merge when approved

## 🐛 Bug Reports

### Before Reporting

1. **Search Existing Issues**
   - Check if the bug is already reported
   - Look for similar issues or feature requests

2. **Reproduce the Issue**
   - Try to reproduce consistently
   - Test with minimal configuration
   - Note system information

3. **Gather Information**
   - Version information
   - System details
   - Configuration files
   - Log output

### Bug Report Template

```markdown
## Bug Description
Clear and concise description of the bug.

## To Reproduce
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

## Expected Behavior
What you expected to happen.

## Actual Behavior
What actually happened.

## Environment
- OS: [e.g. Ubuntu 20.04]
- Go Version: [e.g. 1.21.1]
- Application Version: [e.g. v1.0.0]
- Terminal: [e.g. gnome-terminal, iTerm2]

## Configuration
Relevant configuration files or environment variables.

## Logs
```
Relevant log output
```

## Additional Context
Any additional context, screenshots, or related issues.
```

## 💡 Feature Requests

### Before Requesting

1. **Check Existing Issues**
   - Look for similar requests
   - Check if feature is already planned

2. **Consider Scope**
   - Is this a core feature or plugin candidate?
   - How does it fit with project goals?
   - Would others find it useful?

### Feature Request Template

```markdown
## Feature Description
Clear and concise description of the feature.

## Problem Statement
What problem does this feature solve?

## Proposed Solution
Describe how you envision this feature working.

## Alternatives Considered
Other solutions you've considered.

## Use Cases
Specific scenarios where this would be useful.

## Implementation Ideas
Any thoughts on how this could be implemented.

## Additional Context
Mockups, examples, or related features.
```

## 🤝 Community Guidelines

### Code of Conduct

We are committed to providing a welcoming and inclusive environment:

- **Be respectful**: Treat everyone with respect and kindness
- **Be constructive**: Provide helpful feedback and suggestions
- **Be collaborative**: Work together towards common goals
- **Be patient**: Help newcomers learn and grow

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and community chat
- **Pull Requests**: Code review and technical discussion

### Getting Help

If you need help:

1. **Check Documentation**: Start with the docs in `/docs`
2. **Search Issues**: Look for similar questions
3. **Ask in Discussions**: Community help and questions
4. **Create Issue**: For specific bugs or feature requests

## 🏆 Recognition

### Contributors

We recognize contributors in several ways:

- **Contributors List**: Maintained in README
- **Release Notes**: Contributors credited in releases
- **Special Recognition**: Outstanding contributions highlighted

### Becoming a Maintainer

Active contributors may be invited to become maintainers:

- Consistent high-quality contributions
- Good understanding of codebase
- Helpful to community members
- Aligned with project goals

## 📚 Development Resources

### Learning Resources

- **Go Documentation**: [golang.org/doc](https://golang.org/doc/)
- **Bubbletea Tutorial**: [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- **SSH in Go**: [pkg.go.dev/golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh)
- **Project Architecture**: [Development Guide](development.md)

### Tools

- **Go**: Programming language
- **Bubbletea**: TUI framework
- **Lipgloss**: Styling library
- **Wish**: SSH server library
- **Delve**: Go debugger

### Related Projects

- **Glow**: Markdown reader by Charm
- **Soft Serve**: Git server with TUI
- **VHS**: Terminal session recorder

## 🔄 Release Process

### Versioning

We use semantic versioning (semver):
- **Major**: Breaking changes
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes (backward compatible)

### Release Schedule

- **Patch releases**: As needed for bug fixes
- **Minor releases**: Monthly for new features
- **Major releases**: When significant changes accumulated

### Contributing to Releases

- Test release candidates
- Update documentation
- Help with changelog
- Verify compatibility

---

**Thank You for Contributing!** 🎉

Your contributions make the Void Reavers SSH Reader better for everyone. Whether you're fixing bugs, adding features, improving documentation, or helping other users, every contribution is valuable.

*"In the void between stars, every contribution lights the way for others."* ✨

---

**Next Steps:**
- Set up your development environment
- Look for "good first issue" labels
- Join the community discussions
- Start contributing!