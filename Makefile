.PHONY: menu help setup setup-dev setup-deploy test test-coverage test-verbose build run clean docker-build docker-run lint security-scan pre-commit
.PHONY: deploy deploy-build deploy-logs deploy-restart deploy-rollback deploy-stop deploy-shell deploy-status deploy-env deploy-setup deploy-cleanup
.PHONY: kamal-secrets-setup

.DEFAULT_GOAL := menu

menu:
	@bash scripts/menu.sh

help:
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  Void Chronicles - Makefile Commands"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "📦 Development Commands:"
	@echo "  make setup-dev           - Setup dev environment (Go dependencies)"
	@echo "  make setup-deploy        - Setup deployment environment (Ruby/Kamal/Doppler)"
	@echo "  make test                - Run all tests"
	@echo "  make test-coverage       - Run tests with coverage report"
	@echo "  make test-verbose        - Run tests with verbose output"
	@echo "  make build               - Build the Go binary"
	@echo "  make run                 - Run the application locally (./run.sh)"
	@echo "  make lint                - Format and lint Go code"
	@echo "  make security-scan       - Run security vulnerability scan"
	@echo "  make pre-commit          - Run all checks before committing"
	@echo "  make clean               - Clean build artifacts"
	@echo "  make kamal-secrets-setup - Generate .kamal/secrets file for development"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  make docker-build    - Build Docker image locally"
	@echo "  make docker-run      - Run Docker container locally"
	@echo ""
	@echo "🚀 Deployment Commands (Kamal + Doppler):"
	@echo "  make deploy          - Deploy to production (auto-cleanup)"
	@echo "  make deploy-cleanup  - Stop old containers to free ports"
	@echo "  make deploy-build    - Build and push image only"
	@echo "  make deploy-logs     - Stream production logs"
	@echo "  make deploy-restart  - Restart production containers"
	@echo "  make deploy-rollback - Rollback to previous version"
	@echo "  make deploy-stop     - Stop production containers"
	@echo "  make deploy-shell    - Open shell in production container"
	@echo "  make deploy-status   - Show deployment status"
	@echo "  make deploy-env      - Show production environment variables"
	@echo "  make deploy-setup    - Setup Kamal on new server"
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Setup development environment (Go dependencies only)
setup-dev:
	@echo "🚀 Setting up development environment..."
	@echo ""
	@echo "📋 Checking Go installation..."
	@if ! command -v go > /dev/null; then \
		echo "❌ Go not found. Please install Go first."; \
		echo "   Visit: https://golang.org/doc/install"; \
		exit 1; \
	fi
	@echo "✅ Go $(shell go version | awk '{print $$3}') found"
	@echo ""
	@echo "📦 Installing Go dependencies..."
	cd ssh-reader && go mod download
	cd ssh-reader && go mod tidy
	@echo ""
	@echo "🧪 Running tests to verify setup..."
	@$(MAKE) test
	@echo ""
	@echo "✅ Development environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Run 'make run' to start the local development server"
	@echo "     (No Doppler needed for dev - app uses sensible defaults)"
	@echo "  2. Run 'make test' to run tests"
	@echo "  3. Run 'make setup-deploy' to configure deployment (optional)"

# Setup deployment environment (Ruby, Kamal, Doppler)
setup-deploy:
	@echo "🚀 Setting up deployment environment..."
	@echo ""
	@echo "📋 Checking prerequisites..."
	@if ! command -v ruby > /dev/null; then \
		echo "❌ Ruby not found. Please install Ruby first."; \
		echo "   On macOS: brew install ruby"; \
		echo "   On Ubuntu: sudo apt install ruby-full"; \
		exit 1; \
	fi
	@echo "✅ Ruby $(shell ruby --version | awk '{print $$2}') found"
	@if ! command -v docker > /dev/null; then \
		echo "❌ Docker not found. Please install Docker first."; \
		echo "   Visit: https://docs.docker.com/get-docker/"; \
		exit 1; \
	fi
	@echo "✅ Docker $(shell docker --version | awk '{print $$3}' | tr -d ',') found"
	@if ! command -v doppler > /dev/null; then \
		echo "❌ Doppler CLI not found. Please install Doppler first."; \
		echo "   On macOS: brew install dopplerhq/cli/doppler"; \
		echo "   On Linux: curl -sLf https://cli.doppler.com/install.sh | sh"; \
		exit 1; \
	fi
	@echo "✅ Doppler CLI found"
	@echo ""
	@echo "💎 Installing deployment gems (Kamal)..."
	@if ! command -v bundle > /dev/null; then \
		echo "Installing bundler..."; \
		gem install bundler; \
	fi
	bundle install --only deployment
	@echo "✅ Kamal $(shell bundle exec kamal version 2>/dev/null || echo 'installed')"
	@echo ""
	@echo "🔐 Checking Doppler configuration..."
	@if ! doppler configure get project --plain 2>/dev/null | grep -q "void-reader"; then \
		echo "⚠️  Doppler not configured. You'll need to set up:"; \
		echo "   - prd config (for deployment): doppler setup --project void-reader --config prd"; \
	else \
		echo "✅ Doppler project: void-reader"; \
	fi
	@echo ""
	@echo "🔧 Generating .kamal/secrets file..."
	@$(MAKE) kamal-secrets-setup
	@echo ""
	@echo "✅ Deployment environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Configure Doppler for deployment:"
	@echo "     doppler setup --project void-reader --config prd"
	@echo "  2. Set deployment secret in Doppler prd config:"
	@echo "     - KAMAL_REGISTRY_PASSWORD (GitHub Personal Access Token for ghcr.io)"
	@echo "  3. Test VPS connection: ssh digitalocean-deploy"
	@echo "  4. Deploy: make deploy"

# Full setup (both dev and deploy)
setup: setup-dev setup-deploy

# Run tests
test:
	cd ssh-reader && go test ./...

# Run tests with coverage
test-coverage:
	cd ssh-reader && go test -coverprofile=coverage.out ./...
	cd ssh-reader && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: ssh-reader/coverage.html"

# Run tests with verbose output
test-verbose:
	cd ssh-reader && go test -v ./...

# Build the application
build:
	cd ssh-reader && go build -o void-reader

# Run the application locally
run:
	./run.sh

# Clean build artifacts
clean:
	rm -f ssh-reader/void-reader
	rm -f ssh-reader/coverage.out
	rm -f ssh-reader/coverage.html
	rm -f ssh-reader/gosec-report.json

# Build Docker image
docker-build:
	docker build -t void-chronicles .

# Run Docker container
docker-run:
	docker run -it -p 8080:8080 -p 8443:8443 -p 2222:2222 --env-file .env void-chronicles ./void-reader

# Run linting
lint:
	cd ssh-reader && go fmt ./...
	cd ssh-reader && go vet ./...
	cd ssh-reader && go mod tidy

# Run security scan
security-scan:
	@echo "Installing gosec..."
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	cd ssh-reader && gosec -fmt json -out gosec-report.json ./... || true
	@echo "Security scan complete. Check ssh-reader/gosec-report.json for results"

# Run all checks before committing
pre-commit: lint test security-scan
	@echo "All checks passed!"

# Stop old containers before deployment to avoid port conflicts
deploy-cleanup:
	@echo "🧹 Stopping old containers to free port 22..."
	ssh root@161.35.165.206 -p 1447 "docker stop \$$(docker ps -q --filter 'name=void-chronicles-web') 2>/dev/null || true"
	@echo "✅ Cleanup complete"

# Deploy to production using Doppler for secrets
deploy: pre-commit
	@echo ""
	@echo "✅ All pre-deployment checks passed!"
	@echo ""
	@echo "🚀 Deploying to production with Doppler secrets..."
	@echo "🔐 Using Doppler prd environment..."
	@$(MAKE) deploy-cleanup
	doppler run --project void-reader --config prd -- kamal deploy

# Build and push Docker image only
deploy-build:
	@echo "🔨 Building and pushing Docker image..."
	doppler run --project void-reader --config prd -- kamal build push

# Stream production logs
deploy-logs:
	doppler run --project void-reader --config prd -- kamal app logs -f

# Restart production containers
deploy-restart:
	doppler run --project void-reader --config prd -- kamal app boot

# Rollback to previous version
deploy-rollback:
	doppler run --project void-reader --config prd -- kamal rollback

# Stop production containers
deploy-stop:
	doppler run --project void-reader --config prd -- kamal app stop

# Open shell in production container
deploy-shell:
	doppler run --project void-reader --config prd -- kamal app exec -i bash

# Show deployment status
deploy-status:
	doppler run --project void-reader --config prd -- kamal details

# Show production environment variables
deploy-env:
	doppler run --project void-reader --config prd -- kamal app exec env | grep -v PASSWORD | grep -v TOKEN | sort

# Setup Kamal on new server
deploy-setup:
	doppler run --project void-reader --config prd -- kamal setup

# Generate .kamal/secrets file for development
kamal-secrets-setup:
	@echo "📝 Generating .kamal/secrets file..."
	@mkdir -p .kamal
	@echo "# Kamal secrets file - uses variable substitution with Doppler" > .kamal/secrets
	@echo "# This file is required by Kamal even when using environment variables" >> .kamal/secrets
	@echo "# Doppler injects the actual values during deployment (not at runtime)" >> .kamal/secrets
	@echo "" >> .kamal/secrets
	@echo "KAMAL_REGISTRY_PASSWORD=\$$KAMAL_REGISTRY_PASSWORD" >> .kamal/secrets
	@echo "" >> .kamal/secrets
	@echo "✅ .kamal/secrets file created successfully"
	@echo ""
	@echo "⚠️  This file uses variable substitution (\$$VAR_NAME) so Doppler can inject the actual secrets."
	@echo "   Make sure you have the following secret configured in Doppler prd:"
	@echo "   - KAMAL_REGISTRY_PASSWORD (GitHub Personal Access Token for ghcr.io)"
