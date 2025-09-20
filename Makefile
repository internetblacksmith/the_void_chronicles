.PHONY: test test-coverage test-verbose build run clean docker-build docker-run lint security-scan menu

# Interactive menu as default
menu:
	@echo "Select an option:"
	@echo "1) Run tests"
	@echo "2) Run tests with coverage"
	@echo "3) Run tests with verbose output"
	@echo "4) Build the application"
	@echo "5) Run the application locally"
	@echo "6) Clean build artifacts"
	@echo "7) Build Docker image"
	@echo "8) Run Docker container"
	@echo "9) Run linting"
	@echo "10) Run security scan"
	@echo "11) Generate PDF and EPUB with build metadata"
	@echo "12) Publish PDF and EPUB without build metadata"
	@echo "Enter your choice:"
	@read choice; \
	case $$choice in \
		1) $(MAKE) test ;; \
		2) $(MAKE) test-coverage ;; \
		3) $(MAKE) test-verbose ;; \
		4) $(MAKE) build ;; \
		5) $(MAKE) run ;; \
		6) $(MAKE) clean ;; \
		7) $(MAKE) docker-build ;; \
		8) $(MAKE) docker-run ;; \
		9) $(MAKE) lint ;; \
		10) $(MAKE) security-scan ;; \
		11) ruby markdown_to_kdp_pdf.rb && ruby markdown_to_epub.rb ;; \
		12) ruby markdown_to_publish.rb ;; \
		*) echo "Invalid choice" ;; \
	esac

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
	cd ssh-reader && go run .

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
	docker run -p 8080:8080 -p 2222:2222 void-chronicles

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