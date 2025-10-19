# Agent Guidelines for Void Chronicles

## Project Structure
Dual-component project: (1) Science fiction book series source in Markdown, (2) SSH-based terminal reader application in Go. Markdown is canonical; PDF/EPUB generated from it. SSH reader loads books directly from `book1_void_reavers_source/chapters/`.

## Build/Test Commands
- **Test**: `cd ssh-reader && go test ./...` or `make test`
- **Single test**: `cd ssh-reader && go test -run TestName`
- **Coverage**: `make test-coverage`
- **Build**: `cd ssh-reader && go build` or `make build`
- **Lint**: `cd ssh-reader && go fmt ./... && go vet ./...` or `make lint`
- **Local dev**: `./run.sh` (HTTP:8080, SSH:2222, password: Amigos4Life!)
- **Deploy Kamal**: `kamal deploy` (requires Doppler token and VPS setup per KAMAL_CONFIG_INSTRUCTIONS.md)
- **Generate PDF**: `./markdown_to_kdp_pdf.rb book1_void_reavers_source void_reavers.pdf`
- **Generate EPUB**: `./markdown_to_epub.rb book1_void_reavers_source void_reavers.epub`

## Code Style
- **License**: All Go files start with AGPL-3.0 copyright header (see existing files)
- **Imports**: Standard library first, blank line, then external packages (Bubbletea, Lipgloss, Wish, godotenv)
- **Naming**: camelCase private, PascalCase exported; descriptive names (`LoadBook`, `UserProgress`)
- **Errors**: Always wrap with context: `fmt.Errorf("failed to load book: %w", err)`
- **Types**: Explicit struct tags: `json:"field_name"`, use `-` for non-persisted fields
- **Comments**: Document exported functions only; no inline comments unless critical
- **Paths**: Use `filepath.Join()` for cross-platform compatibility

## Key Architecture
- Dual servers: HTTP (90s homepage disguise) + SSH (reading interface) in `main.go`
- TUI states: Main menu (split-view library), chapter list, reading view, progress, about
- Progress tracking: JSON persistence in `.void_reader_data/username.json`
- Book loading: Markdown parser in `book.go` reads from `chapters/*.md`
- Environment: Variables loaded via `godotenv` with fallback defaults
- Deployment: Kamal orchestration with Traefik SSL (HTTP), direct SSH port publishing (2222), Doppler secrets, persistent volumes for progress and SSH keys

## Critical Commit Policy
**Documentation-First**: Before ANY commit, verify ALL documentation matches code (README, DEPLOYMENT.md, guides, file paths). Documentation drift is unacceptable. Workflow: (1) Code changes, (2) Update docs, (3) Verify accuracy, (4) Commit.
