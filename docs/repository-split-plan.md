# Repository Split Plan: Void Reader & Void Chronicles

**Status:** Planning  
**Created:** 2025-10-25  
**Goal:** Split monorepo into generic SSH book reader and book content repository

---

## Executive Summary

Split the current monorepo into:
1. **void-reader** - Generic, reusable SSH terminal book reader (AGPL-3.0)
2. **void-chronicles** - Book content for Void Chronicles series (CC BY-NC-SA 4.0)

**Benefits:**
- Clean separation of concerns (reader vs content)
- void-reader can be used by other authors
- Proper license separation
- Independent versioning and deployment
- Clearer contribution guidelines

---

## Current State Analysis

### Hardcoded Dependencies

**In `ssh-reader/main.go`:**
```go
Line 732: baseDir := "book1_void_reavers_source"
Line 735: metadataPath := filepath.Join(baseDir, "metadata.yaml")
```

**In `ssh-reader/book.go`:**
- `LoadBook()` expects specific directory structure: `chapters/*.md`
- Metadata parsing assumes specific YAML schema

**In `ssh-reader/series.json`:**
- Hardcoded Void Chronicles series information
- Specific book titles, descriptions, status

**In Docker/Deployment:**
- `Dockerfile` copies `book1_void_reavers_source` directly
- `config/deploy.yml` builds image with embedded content
- No runtime content injection

**In Documentation:**
- README.md mixes reader and book information
- AGENTS.md references specific book paths
- Build scripts (markdown_to_epub.rb, markdown_to_kdp_pdf.rb) are book-specific

---

## Target Architecture

### Repository 1: void-reader (Generic Reader)

**Purpose:** Reusable SSH terminal application for reading books

**Key Features:**
- Environment variable configuration
- Runtime book loading (not compile-time)
- Generic series.json schema
- Pluggable content sources (local filesystem, Git, HTTP)
- Clear API for book authors

**Configuration via Environment:**
```bash
BOOK_SOURCE_TYPE=filesystem|git|http
BOOK_BASE_DIR=/books
BOOK_SERIES_FILE=/books/series.json
BOOK_METADATA_FILE=/books/metadata.yaml
DATA_DIR=/data/reader_data
```

**Directory Structure:**
```
void-reader/
├── main.go
├── book.go
├── progress.go
├── storage.go
├── ratelimit.go
├── go.mod
├── go.sum
├── Dockerfile (generic, expects volume mount)
├── docker-compose.yml (with volume examples)
├── README.md (for reader users & developers)
├── docs/
│   ├── book-author-guide.md (NEW)
│   ├── deployment.md
│   ├── configuration.md (NEW)
│   └── schema/
│       ├── series.json.schema (NEW)
│       └── metadata.yaml.schema (NEW)
├── examples/
│   ├── series.json.example
│   ├── metadata.yaml.example
│   └── sample-book/
│       ├── chapters/
│       │   └── chapter-01.md
│       ├── metadata.yaml
│       └── book.md
└── LICENSE (AGPL-3.0)
```

### Repository 2: void-chronicles (Book Content)

**Purpose:** Source content for Void Chronicles series

**Directory Structure:**
```
void-chronicles/
├── book1-void-reavers/
│   ├── chapters/
│   │   ├── chapter-01.md
│   │   └── ...
│   ├── metadata.yaml
│   └── book.md
├── book2-stellar-tomb/ (future)
├── book3-omega-directive/ (future)
├── series.json
├── scripts/
│   ├── markdown_to_epub.rb
│   └── markdown_to_kdp_pdf.rb
├── docs/
│   └── writing-guide.md
├── Dockerfile (uses void-reader base image)
├── config/
│   └── deploy.yml (Kamal config for void-chronicles deployment)
├── LICENSE (CC BY-NC-SA 4.0)
├── README.md (about the series)
└── void_chronicles_series_bible.md
```

---

## Implementation Phases

### Phase 1: Prepare void-reader for Extraction (Current Repo)

**Step 1.1: Make Reader Configuration-Driven**

- [ ] Add environment variable loading for book paths
- [ ] Replace hardcoded `book1_void_reavers_source` with `BOOK_BASE_DIR`
- [ ] Create config.go for centralized configuration
- [ ] Add validation for required environment variables
- [ ] Test with existing book content (no behavior change)

**Code Changes:**
```go
// ssh-reader/config.go (NEW)
type Config struct {
    BookSourceType  string // "filesystem", "git", "http"
    BookBaseDir     string
    SeriesFile      string
    MetadataFile    string
    DataDir         string
    HTTPPort        int
    HTTPSPort       int
    SSHPort         int
}

func LoadConfig() (*Config, error) {
    return &Config{
        BookSourceType: getEnv("BOOK_SOURCE_TYPE", "filesystem"),
        BookBaseDir:    getEnv("BOOK_BASE_DIR", "book1_void_reavers_source"),
        SeriesFile:     getEnv("BOOK_SERIES_FILE", "ssh-reader/series.json"),
        MetadataFile:   getEnv("BOOK_METADATA_FILE", "metadata.yaml"),
        DataDir:        getEnv("DATA_DIR", ".void_reader_data"),
        HTTPPort:       getEnvInt("HTTP_PORT", 8080),
        HTTPSPort:      getEnvInt("HTTPS_PORT", 8443),
        SSHPort:        getEnvInt("SSH_PORT", 2222),
    }, nil
}
```

**Step 1.2: Document Generic Schema**

- [ ] Create `docs/schema/series.json.schema`
- [ ] Create `docs/schema/metadata.yaml.schema`
- [ ] Create `docs/book-author-guide.md`
- [ ] Add examples/ directory with sample book

**Step 1.3: Update Dockerfile for Volume Mounting**

- [ ] Modify Dockerfile to expect `/books` volume mount
- [ ] Add BOOK_BASE_DIR=/books environment variable
- [ ] Update docker-compose.yml with volume example
- [ ] Test local deployment with mounted content

**Step 1.4: Test & Validate**

- [ ] Run all tests: `make test`
- [ ] Test local dev with env vars: `BOOK_BASE_DIR=book1_void_reavers_source ./run.sh`
- [ ] Test Docker with volume mount
- [ ] Verify no regressions

**Commit:** "refactor: make reader configuration-driven for repository split preparation"

---

### Phase 2: Create void-reader Repository

**Step 2.1: Initialize New Repository**

```bash
mkdir void-reader
cd void-reader
git init
gh repo create internetblacksmith/void-reader --public --description "SSH terminal book reader - read books in your terminal"
```

**Step 2.2: Copy Reader Code**

- [ ] Copy `ssh-reader/*` (Go code) to void-reader root
- [ ] Copy `docs/deployment-alternatives.md` → `docs/deployment.md`
- [ ] Copy relevant deployment scripts (run.sh, Dockerfile)
- [ ] Create new README.md focused on reader (not Void Chronicles)
- [ ] Copy LICENSE (AGPL-3.0)
- [ ] Create .gitignore for Go project

**Step 2.3: Create Documentation**

- [ ] `README.md` - Generic reader introduction
- [ ] `docs/book-author-guide.md` - How to use void-reader for your books
- [ ] `docs/configuration.md` - Environment variables and config
- [ ] `docs/deployment.md` - Deployment options (Docker, Kamal, Railway)
- [ ] `examples/` - Sample book structure
- [ ] `CONTRIBUTING.md` - How to contribute to reader

**Step 2.4: Setup CI/CD**

- [ ] Copy `.github/workflows/test.yml`
- [ ] Create `.github/workflows/docker-publish.yml` (publish to ghcr.io)
- [ ] Update workflow to use example book for tests
- [ ] Test GitHub Actions

**Step 2.5: First Release**

- [ ] Tag v0.1.0
- [ ] Create GitHub release with usage instructions
- [ ] Publish Docker image: `ghcr.io/internetblacksmith/void-reader:latest`

---

### Phase 3: Create void-chronicles Repository

**Step 3.1: Initialize New Repository**

```bash
mkdir void-chronicles
cd void-chronicles
git init
gh repo create internetblacksmith/void-chronicles --public --description "Void Chronicles - A space pirate adventure series"
```

**Step 3.2: Copy Book Content**

- [ ] Copy `book1_void_reavers_source/` → `book1-void-reavers/`
- [ ] Copy `ssh-reader/series.json` → `series.json`
- [ ] Copy `void_chronicles_series_bible.md`
- [ ] Copy `markdown_to_epub.rb`, `markdown_to_kdp_pdf.rb` → `scripts/`
- [ ] Copy `LICENSE-CONTENT.md` → `LICENSE`
- [ ] Copy `MARKDOWN_STYLE_GUIDE.md` → `docs/`
- [ ] Create new README.md focused on book series

**Step 3.3: Create Deployment Configuration**

**Dockerfile:**
```dockerfile
FROM ghcr.io/internetblacksmith/void-reader:latest

COPY book1-void-reavers /books/book1-void-reavers
COPY series.json /books/series.json

ENV BOOK_BASE_DIR=/books/book1-void-reavers
ENV BOOK_SERIES_FILE=/books/series.json

CMD ["doppler", "run", "--", "/app/void-reader"]
```

**config/deploy.yml:**
```yaml
service: void-chronicles
image: internetblacksmith/void-chronicles

servers:
  web:
    hosts:
      - 161.35.165.206

registry:
  server: ghcr.io
  username: internetblacksmith
  password:
    - KAMAL_REGISTRY_PASSWORD

env:
  secret:
    - DOPPLER_TOKEN
  clear:
    BOOK_BASE_DIR: /books/book1-void-reavers
    BOOK_SERIES_FILE: /books/series.json
    DATA_DIR: /data/void_reader_data

volumes:
  - "void-data:/data/void_reader_data"
  - "void-ssl:/ssl"

builder:
  multiarch: false

healthcheck:
  path: /health
  port: 8080
  interval: 10s

accessories:
  files:
    port: 80:8080
    port: 443:8443
    port: 22:2222
```

**Step 3.4: Test Deployment**

- [ ] Build Docker image locally
- [ ] Test with docker-compose
- [ ] Deploy to staging (if available)
- [ ] Verify all functionality works

---

### Phase 4: Migrate Production Deployment

**Step 4.1: Prepare Migration**

- [ ] Document current deployment state
- [ ] Backup current user progress data: `/data/void_reader_data/*.json`
- [ ] Create rollback plan
- [ ] Schedule maintenance window (low-traffic time)

**Step 4.2: Update Deployment**

```bash
# On VPS (ssh -p 1447 root@161.35.165.206)
cd /root/void-chronicles
git pull origin main

# Setup Kamal secrets (if needed)
make kamal-secrets-setup

# Deploy new version
doppler run --project void-reader --config prd -- kamal deploy
```

**Step 4.3: Verify Migration**

- [ ] Test HTTP endpoint: http://161.35.165.206
- [ ] Test HTTPS endpoint: https://voidreader.fun
- [ ] Test SSH: `ssh read@voidreader.fun`
- [ ] Verify user progress persists
- [ ] Check all menu options (library, progress, about, license)
- [ ] Monitor logs for errors

**Step 4.4: Update DNS (if needed)**

- [ ] Verify voidreader.fun points to 161.35.165.206
- [ ] No changes needed unless domain changes

---

### Phase 5: Update Original Repository

**Step 5.1: Archive or Archive Marker**

Options:
1. **Archive entire repo** (recommended if no other use)
2. **Add deprecation notice** to README
3. **Keep as monorepo** with submodules pointing to new repos

**Recommended: Add Deprecation Notice**

Update `README.md`:
```markdown
# ⚠️ This Repository Has Been Split

This monorepo has been split into two focused repositories:

- **[void-reader](https://github.com/internetblacksmith/void-reader)** - Generic SSH terminal book reader (AGPL-3.0)
- **[void-chronicles](https://github.com/internetblacksmith/void-chronicles)** - Void Chronicles book series content (CC BY-NC-SA 4.0)

**For book readers:** Visit https://voidreader.fun or `ssh read@voidreader.fun`

**For developers/authors:** See [void-reader](https://github.com/internetblacksmith/void-reader) to use the reader for your own books

**For book content:** See [void-chronicles](https://github.com/internetblacksmith/void-chronicles) for the Void Chronicles series source

This repository is kept for historical reference but is no longer actively maintained.

---

## Original README (Archived)
[... rest of original README ...]
```

**Step 5.2: Update Documentation**

- [ ] Add MIGRATION.md documenting the split
- [ ] Update all links to point to new repos
- [ ] Add badges linking to new repositories
- [ ] Update CONTRIBUTING.md to point to appropriate repo

---

## Configuration Schema

### series.json Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Book Series Metadata",
  "type": "object",
  "required": ["series", "books"],
  "properties": {
    "series": {
      "type": "object",
      "required": ["title", "author"],
      "properties": {
        "title": {"type": "string"},
        "author": {"type": "string"},
        "description": {"type": "string"},
        "genre": {"type": "array", "items": {"type": "string"}},
        "website": {"type": "string", "format": "uri"}
      }
    },
    "books": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "title", "status"],
        "properties": {
          "id": {"type": "string"},
          "title": {"type": "string"},
          "subtitle": {"type": "string"},
          "status": {
            "type": "string",
            "enum": ["published", "in-progress", "planned"]
          },
          "description": {"type": "string"},
          "publishDate": {"type": "string", "format": "date"},
          "chapters": {"type": "integer", "minimum": 1},
          "baseDir": {"type": "string"}
        }
      }
    }
  }
}
```

### metadata.yaml Schema

```yaml
# Required fields
title: "Book Title"
author: "Author Name"
language: "en-US"

# Optional fields
subtitle: "Book Subtitle"
publisher: "Publisher Name"
publishDate: "2025-01-01"
isbn: "978-1-234567-89-0"
genre:
  - "Science Fiction"
  - "Space Opera"
description: "Book description for retailers and readers"
keywords:
  - "keyword1"
  - "keyword2"
copyright: "© 2025 Author Name. Licensed under CC BY-NC-SA 4.0"
coverImage: "path/to/cover.jpg"
```

---

## Testing Strategy

### void-reader Testing

**Unit Tests:**
- [ ] Config loading with env vars
- [ ] Book loading from different base directories
- [ ] Series.json parsing (valid & invalid)
- [ ] Metadata.yaml parsing (valid & invalid)
- [ ] Progress tracking with different data directories

**Integration Tests:**
- [ ] Full reader with example book
- [ ] Multi-book series navigation
- [ ] Progress persistence across restarts
- [ ] SSH authentication and session management

**Docker Tests:**
- [ ] Build generic image
- [ ] Mount volume and verify book loads
- [ ] Test with multiple book sources

### void-chronicles Testing

**Content Tests:**
- [ ] Markdown validation (all chapters)
- [ ] Metadata validation (series.json, metadata.yaml)
- [ ] Link validation (internal references)
- [ ] EPUB generation
- [ ] PDF generation

**Deployment Tests:**
- [ ] Docker build with void-reader base
- [ ] Kamal deployment to staging
- [ ] End-to-end reader functionality

---

## Rollback Plan

If production deployment fails:

```bash
# On VPS
cd /root/space_pirate  # Old monorepo location
doppler run --project void-reader --config prd -- kamal deploy

# Restore user data if needed
docker cp void-data:/data/void_reader_data /backup/
```

**Pre-deployment checklist:**
1. Backup user progress data
2. Test new deployment in local Docker
3. Verify all environment variables are set
4. Document current container state
5. Have SSH access to VPS ready

---

## Timeline Estimate

| Phase | Duration | Dependencies |
|-------|----------|--------------|
| Phase 1: Refactor current repo | 4-6 hours | None |
| Phase 2: Create void-reader | 3-4 hours | Phase 1 complete |
| Phase 3: Create void-chronicles | 2-3 hours | Phase 2 complete |
| Phase 4: Migrate production | 1-2 hours | Phase 2 & 3 complete, tested |
| Phase 5: Update original repo | 1 hour | Phase 4 complete |
| **Total** | **11-16 hours** | Sequential execution |

**Recommended schedule:**
- **Day 1:** Phase 1 (refactor and test)
- **Day 2:** Phase 2 & 3 (create new repos)
- **Day 3:** Phase 4 (production migration during low-traffic window)
- **Day 4:** Phase 5 (cleanup and documentation)

---

## Success Criteria

- [ ] void-reader works standalone with example book
- [ ] void-reader Docker image published to ghcr.io
- [ ] void-chronicles deploys successfully with void-reader base
- [ ] Production deployment maintains all functionality
- [ ] User progress data persists through migration
- [ ] No downtime > 5 minutes during migration
- [ ] All documentation updated and accurate
- [ ] Both repositories have clear README and contribution guidelines
- [ ] CI/CD pipelines working for both repos
- [ ] License separation properly documented

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Data loss during migration | High | Backup all user progress before deployment |
| Deployment downtime | Medium | Deploy during low-traffic hours, have rollback ready |
| Breaking changes in reader | High | Thorough testing with existing content before split |
| Missing dependencies | Medium | Document all env vars, test in clean environment |
| DNS propagation delay | Low | Keep old deployment running until DNS stable |
| Docker registry issues | Medium | Pre-push images, verify access before deployment |

---

## Future Enhancements (Post-Split)

### void-reader
- [ ] Git-based book sources (clone repos at runtime)
- [ ] HTTP-based book sources (fetch from CDN)
- [ ] Multi-series support (browse multiple series)
- [ ] Search functionality across books
- [ ] Bookmarks and annotations
- [ ] Export progress/annotations
- [ ] Plugin system for custom themes
- [ ] Analytics dashboard (privacy-focused)

### void-chronicles
- [ ] Book 2: Stellar Tomb (when ready)
- [ ] Book 3: Omega Directive (when ready)
- [ ] Automated EPUB/PDF builds via GitHub Actions
- [ ] Preview deployment for draft chapters
- [ ] Reader feedback integration
- [ ] Translations (if demand exists)

---

## Questions to Resolve

1. **Docker Registry:** Keep both images under `internetblacksmith/*` or create separate org?
2. **Versioning:** Semantic versioning for reader? How to version book content?
3. **Data Migration:** Keep shared Doppler project or create separate ones?
4. **Domain:** Use voidreader.fun for both? Or separate domains?
5. **CI/CD:** Use same GitHub Actions workflow or separate?
6. **Support:** Create separate Discord/forum for reader vs books?

---

## References

- Current deployment: https://voidreader.fun
- VPS: 161.35.165.206 (SSH port 1447)
- Docker registry: ghcr.io/internetblacksmith
- Doppler project: void-reader (config: prd)
- Current monorepo: https://github.com/internetblacksmith/space_pirate (assumed)

---

**Next Steps:**
1. Review this plan
2. Answer open questions
3. Begin Phase 1 implementation
4. Test thoroughly before proceeding to Phase 2

**Last Updated:** 2025-10-25
