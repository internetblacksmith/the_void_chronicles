# Contributing to The Void Chronicles

First off, thank you for considering contributing to The Void Chronicles! It's people like you that make this project such a great tool for the community.

## Code of Conduct

By participating in this project, you are expected to uphold our [Code of Conduct](CODE_OF_CONDUCT.md).

## How Can I Contribute?

### 🐛 Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples**
- **Describe the behavior you observed and expected**
- **Include screenshots if relevant**
- **Include your environment details** (OS, Go version, etc.)

### 💡 Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description of the proposed enhancement**
- **Explain why this enhancement would be useful**
- **List any alternative solutions you've considered**

### 📝 Contributing to Documentation

- Fix typos, improve clarity, add examples
- Translate documentation to other languages
- Write tutorials or blog posts about the project

### 🎨 Contributing to Book Content

The narrative content is under CC-BY-SA 4.0. Contributions could include:

- Fixing typos and grammar
- Suggesting plot improvements (open an issue first)
- Creating fan art or illustrations
- Writing additional short stories in the universe

### 💻 Contributing Code

#### Development Setup

1. **Fork the repository**
   ```bash
   git clone https://github.com/yourusername/void-chronicles.git
   cd void-chronicles
   ```

2. **Set up development environment**
   ```bash
   cd ssh-reader
   go mod download
   ```

3. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

4. **Make your changes**
   - Write clean, readable code
   - Add comments for complex logic
   - Follow existing code style
   - Add tests for new features

5. **Test your changes**
   ```bash
   go test ./...
   ./run.sh  # Test manually
   ```

6. **Commit your changes**
   ```bash
   git commit -m "feat: add amazing feature"
   ```

#### Commit Message Guidelines

We use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Maintenance tasks

Examples:
```
feat: add bookmark export functionality
fix: correct chapter navigation in vim mode
docs: update Railway deployment guide
```

#### Pull Request Process

1. **Update documentation** - Update README.md and other docs if needed
2. **Update CLAUDE.md** - If you changed architecture or commands
3. **Add tests** - Ensure test coverage for new features
4. **Update version** - Follow semantic versioning if applicable
5. **Create Pull Request** - Use the PR template

### 🧪 Testing Guidelines

- Write unit tests for new functions
- Test edge cases and error conditions
- Ensure all tests pass: `go test ./...`
- Test manually with different terminal sizes
- Test SSH connectivity and authentication

### 📚 Style Guides

#### Go Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Use meaningful variable names
- Add comments for exported functions
- Handle errors explicitly

Example:
```go
// LoadBook reads and parses book content from the specified directory
func LoadBook(bookDir string) (*Book, error) {
    if bookDir == "" {
        return nil, fmt.Errorf("book directory cannot be empty")
    }
    // Implementation...
}
```

#### Markdown Style (for book content)

- Use single asterisks for italics: `*Ship Name*`
- Use double asterisks for bold: `**important**`
- Use `# Chapter Title` for chapter headings
- Use `* * *` for scene breaks
- Keep lines under 80 characters when possible

### 🚀 Areas We Need Help With

- **Terminal UI improvements** - Better layouts, themes, animations
- **Performance optimization** - Faster book loading, reduced memory usage
- **Platform support** - Windows compatibility, mobile SSH clients
- **Accessibility** - Screen reader support, high contrast modes
- **Localization** - Translate UI to other languages
- **Testing** - Unit tests, integration tests, CI/CD setup
- **Documentation** - Tutorials, videos, examples

## Recognition

Contributors will be:
- Listed in the project README
- Mentioned in release notes
- Given credit in commit messages

## Questions?

Feel free to:
- Open an issue for questions
- Start a discussion in GitHub Discussions
- Contact the maintainers

## License

By contributing, you agree that your contributions will be licensed under:
- **Code**: GNU Affero General Public License v3.0
- **Book Content**: Creative Commons CC-BY-SA 4.0

Thank you for contributing to The Void Chronicles! 🚀📚