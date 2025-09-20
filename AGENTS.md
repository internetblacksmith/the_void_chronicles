# AGENTS.md

This document provides guidelines for agentic coding assistants working in this repository.

## Build, Lint, and Test Commands

### Build the Application
```bash
make build
# or
./build.sh
```

### Run the Application Locally
```bash
make run
# or
./run.sh
```

### Run Tests
- Run all tests:
  ```bash
  make test
  ```
- Run tests with coverage:
  ```bash
  make test-coverage
  ```
- Run tests with verbose output:
  ```bash
  make test-verbose
  ```
- Run a single test:
  ```bash
  cd ssh-reader && go test -run TestName
  ```

### Lint the Code
```bash
make lint
```

## Code Style Guidelines

### Imports
- Group imports into standard library, third-party, and local packages.
- Use parentheses for grouped imports.

### Formatting
- Use `go fmt` to format code.
- Ensure `go vet` passes without errors.

### Types and Naming Conventions
- Use `CamelCase` for types and structs.
- Use `snake_case` for JSON keys.
- Prefix private functions with lowercase letters.

### Error Handling
- Return errors using `fmt.Errorf` with `%w` for wrapping.
- Log errors with context using `log.Printf`.

### General Practices
- Avoid global variables.
- Write unit tests for all new functions.
- Follow the Markdown style guide for book content.

### Markdown Style Highlights
- Use British English spelling.
- Italicize ship names (e.g., *Crimson Revenge*).
- Use `* * *` for scene breaks.

## Documentation-First Commit Policy

**📝 Documentation-First Commit Policy:**

You MAY commit changes autonomously, but ONLY after:
1. **Documentation is FULLY updated** to reflect all code changes
2. **README files** match the current implementation exactly
3. **Deployment guides** reflect current configuration
4. **File paths in docs** are verified and correct

**The workflow MUST be:**
1. Make code changes
2. Update ALL relevant documentation
3. Verify docs match code state
4. Then commit with a clear message

Documentation drift is unacceptable - docs must ALWAYS match the code state before any commit.

For more details, refer to `MARKDOWN_STYLE_GUIDE.md`.