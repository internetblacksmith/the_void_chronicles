.PHONY: menu help setup test test-coverage test-verbose build run clean docker-build docker-run lint security-scan pre-commit
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
	@echo "  make setup               - Setup dev environment (install all dependencies)"
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

# Setup development environment
setup:
	@echo "🚀 Setting up development environment..."
	@echo ""
	@echo "📦 Installing Go dependencies..."
	cd ssh-reader && go mod download
	cd ssh-reader && go mod tidy
	@echo ""
	@echo "💎 Installing Ruby and Bundler for Kamal deployment..."
	@if ! command -v ruby > /dev/null; then \
		echo "❌ Ruby not found. Please install Ruby first."; \
		exit 1; \
	fi
	@if ! command -v bundle > /dev/null; then \
		echo "Installing bundler..."; \
		gem install bundler; \
	fi
	bundle install
	@echo ""
	@echo "🔧 Generating .kamal/secrets file..."
	@$(MAKE) kamal-secrets-setup
	@echo ""
	@echo "🧪 Running tests to verify setup..."
	@$(MAKE) test
	@echo ""
	@echo "✅ Development environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Ensure Doppler CLI is installed and configured"
	@echo "  2. Run 'make run' to start the local development server"
	@echo "  3. Run 'make test' to run tests"
	@echo "  4. Run 'make deploy' to deploy to production (requires Doppler secrets)"

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
deploy:
	@echo "🚀 Deploying to production with Doppler secrets..."
	@$(MAKE) deploy-cleanup
	doppler run --project void-reader --config prd --command='bash -c "export KAMAL_REGISTRY_PASSWORD && export DOPPLER_TOKEN && kamal deploy"'

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
	@echo "# Doppler injects the actual values at runtime" >> .kamal/secrets
	@echo "" >> .kamal/secrets
	@echo "KAMAL_REGISTRY_PASSWORD=\$$KAMAL_REGISTRY_PASSWORD" >> .kamal/secrets
	@echo "DOPPLER_TOKEN=\$$DOPPLER_TOKEN" >> .kamal/secrets
	@echo "" >> .kamal/secrets
	@echo "✅ .kamal/secrets file created successfully"
	@echo ""
	@echo "⚠️  This file uses variable substitution (\$$VAR_NAME) so Doppler can inject the actual secrets."
	@echo "   Make sure you have the following secrets configured in Doppler:"
	@echo "   - KAMAL_REGISTRY_PASSWORD (GitHub Personal Access Token)"
	@echo "   - DOPPLER_TOKEN (Service token for container runtime)"
