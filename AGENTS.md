# Agent Guidelines for Void Chronicles

## Project Structure
Dual-component project: (1) Science fiction book series source in Markdown, (2) SSH-based terminal reader application in Go. Markdown is canonical; PDF/EPUB generated from it. SSH reader loads books directly from `book1_void_reavers_source/chapters/`.

## Build/Test Commands
- **Interactive Menu**: `make` or `make menu` (launches interactive menu with all commands)
- **Test**: `cd ssh-reader && go test ./...` or `make test`
- **Single test**: `cd ssh-reader && go test -run TestName`
- **Coverage**: `make test-coverage`
- **Build**: `cd ssh-reader && go build` or `make build`
- **Lint**: `cd ssh-reader && go fmt ./... && go vet ./...` or `make lint`
- **Local dev**: `./run.sh` or `make run` (HTTP:8080, HTTPS:8443, SSH:2222, password: Amigos4Life!)
- **Setup Kamal secrets**: `make kamal-secrets-setup` (generates `.kamal/secrets` file with variable substitution)
- **Deploy Kamal**: `make deploy` or `doppler run --project void-reader --config prd -- kamal deploy`
- **SSL renewal**: `./renew-ssl-certs.sh` (Let's Encrypt certificate renewal and Docker volume copy)
- **Generate PDF**: `./markdown_to_kdp_pdf.rb book1_void_reavers_source void_reavers.pdf`
- **Generate EPUB**: `./markdown_to_epub.rb book1_void_reavers_source void_reavers.epub`

## Code Style
- **License Headers**: All Go files start with AGPL-3.0 copyright header (see existing files)
- **Book Content License**: Books are CC BY-NC-SA 4.0 (see LICENSE-CONTENT.md)
- **Imports**: Standard library first, blank line, then external packages (Bubbletea, Lipgloss, Wish, godotenv)
- **Naming**: camelCase private, PascalCase exported; descriptive names (`LoadBook`, `UserProgress`)
- **Errors**: Always wrap with context: `fmt.Errorf("failed to load book: %w", err)`
- **Types**: Explicit struct tags: `json:"field_name"`, use `-` for non-persisted fields
- **Comments**: Document exported functions only; no inline comments unless critical
- **Paths**: Use `filepath.Join()` for cross-platform compatibility
- **UI Consistency**: All content views must use consistent dimensions matching menu panels

## Key Architecture
- Triple servers: HTTP (8080), HTTPS (8443), SSH (2222) all in `main.go`
- TUI states: Main menu (split-view library), chapter list, reading view, progress, about, license
- Progress tracking: JSON persistence in `/data/void_reader_data/username.json` (production) or `.void_reader_data/username.json` (local dev)
- Book loading: Markdown parser in `book.go` reads from `chapters/*.md`
- Environment: Variables loaded via `godotenv` with fallback defaults, Doppler in production
- HTTPS: Native TLS support with graceful fallback if certificates not found
- Deployment: Kamal orchestration with direct port mapping (80→8080 HTTP, 443→8443 HTTPS, 22→2222 SSH), Doppler secrets, persistent volumes (void-data for progress, void-ssl for certificates)
- UI Layout: All views use consistent dimensions (width - 6, height - 8) with rounded borders, padding (1, 2), and centered alignment

## Critical Rules

### Rule 1: Always Use Doppler for Secrets
**NEVER** remove Doppler integration from this project. All production secrets MUST be managed via Doppler. The Dockerfile MUST include Doppler CLI installation and `CMD ["doppler", "run", "--", "./void-reader"]`. The `config/deploy.yml` MUST have Doppler env configuration with `DOPPLER_TOKEN` as a secret.

### Rule 2: SSH Port is 22
The application SSH server listens on container port 2222, mapped to host port **22** (not 2222). System SSH runs on port 1447, so port 22 is available. Port mapping in `config/deploy.yml` MUST be `"22:2222"`.

### Rule 3: Kamal Secrets File Required
Kamal requires a `.kamal/secrets` file even when using Doppler environment variables. This file MUST use variable substitution format (`$VAR_NAME`) so Doppler can inject actual values at runtime. Use `make kamal-secrets-setup` to generate this file. The secrets file contains:
- `KAMAL_REGISTRY_PASSWORD=$KAMAL_REGISTRY_PASSWORD` (GitHub PAT)
- `DOPPLER_TOKEN=$DOPPLER_TOKEN` (service token for container runtime)

## Licensing Structure
This project uses a dual-license approach:
- **Book Content**: Creative Commons BY-NC-SA 4.0 (see LICENSE-CONTENT.md, metadata.yaml, series.json)
- **SSH Reader Code**: GNU AGPL-3.0 (see LICENSE, Go file headers)
- License info is displayed in the SSH reader's License view (accessible from main menu)
- Both licenses must be maintained and documented in all relevant locations

## Critical Commit Policy
**Documentation-First**: Before ANY commit, verify ALL documentation matches code (README, DEPLOYMENT.md, AGENTS.md, guides, file paths). Documentation drift is unacceptable. Workflow: (1) Code changes, (2) Update docs, (3) Verify accuracy, (4) Commit.
