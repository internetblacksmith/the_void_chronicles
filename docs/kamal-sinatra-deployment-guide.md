# Kamal Deployment Guide for Sinatra Applications

This guide provides instructions for deploying a Sinatra (Ruby) application using Kamal, based on the patterns established in the Void Chronicles project.

## Prerequisites

- Ruby application (Sinatra framework)
- Docker installed locally
- Kamal gem installed (`gem install kamal`)
- Doppler CLI installed (for secrets management)
- VPS server with Docker installed
- Domain name configured

## Project Structure

Your Sinatra project should follow this structure:

```
your-sinatra-app/
├── .kamal/
│   └── secrets              # Generated, not committed
├── config/
│   └── deploy.yml          # Kamal configuration
├── .dockerignore
├── .env.example
├── .gitignore
├── Dockerfile
├── Gemfile
├── Gemfile.lock
├── Makefile                # Build automation
├── config.ru               # Rack config
└── app.rb                  # Your Sinatra app
```

## Step 1: Create Dockerfile

Create a `Dockerfile` for your Sinatra application:

```dockerfile
# Use official Ruby image
FROM ruby:3.2-slim

# Install dependencies
RUN apt-get update -qq && \
    apt-get install -y build-essential curl && \
    rm -rf /var/lib/apt/lists/*

# Install Doppler CLI for secrets management
RUN curl -sLf --retry 3 --tlsv1.2 --proto "=https" 'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' | gpg --dearmor -o /usr/share/keyrings/doppler-archive-keyring.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/doppler-archive-keyring.gpg] https://packages.doppler.com/public/cli/deb/debian any-version main" | tee /etc/apt/sources.list.d/doppler-cli.list && \
    apt-get update && \
    apt-get -y install doppler

# Set working directory
WORKDIR /app

# Copy Gemfile and install dependencies
COPY Gemfile Gemfile.lock ./
RUN bundle install --without development test

# Copy application code
COPY . .

# Expose port (adjust as needed)
EXPOSE 4567

# Use Doppler to inject secrets and run the application
CMD ["doppler", "run", "--", "bundle", "exec", "rackup", "-o", "0.0.0.0", "-p", "4567"]
```

## Step 2: Create .dockerignore

```
.git
.kamal
.env
.env.*
tmp/
log/
coverage/
*.md
Makefile
```

## Step 3: Create Kamal Configuration

Create `config/deploy.yml`:

```yaml
service: your-app-name
image: your-dockerhub-username/your-app-name

servers:
  web:
    hosts:
      - YOUR_SERVER_IP
    labels:
      traefik.http.routers.your-app.rule: Host(`your-domain.com`)
      traefik.http.routers.your-app.tls: true
      traefik.http.routers.your-app.tls.certresolver: letsencrypt
      traefik.http.routers.your-app.entrypoints: websecure
      traefik.http.services.your-app.loadbalancer.server.port: 4567
    options:
      network: "private"

registry:
  server: ghcr.io
  username: your-github-username
  password:
    - KAMAL_REGISTRY_PASSWORD

env:
  clear:
    RACK_ENV: production
  secret:
    - DOPPLER_TOKEN

traefik:
  options:
    publish:
      - "80:80"
      - "443:443"
    volume:
      - "/letsencrypt/acme.json:/letsencrypt/acme.json"
    network: "private"
  args:
    entryPoints.web.address: ":80"
    entryPoints.websecure.address: ":443"
    entryPoints.web.http.redirections.entryPoint.to: websecure
    entryPoints.web.http.redirections.entryPoint.scheme: https
    certificatesResolvers.letsencrypt.acme.email: "your-email@example.com"
    certificatesResolvers.letsencrypt.acme.storage: "/letsencrypt/acme.json"
    certificatesResolvers.letsencrypt.acme.httpchallenge: true
    certificatesResolvers.letsencrypt.acme.httpchallenge.entrypoint: web

volumes:
  - "app-data:/app/data"

# Healthcheck endpoint
healthcheck:
  path: /health
  port: 4567
  interval: 10s
  timeout: 5s
```

## Step 4: Create Makefile

Create a `Makefile` with the same format as Void Chronicles:

```makefile
.PHONY: help menu test build run deploy kamal-secrets-setup docker-build docker-run clean

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
MAGENTA := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

help: menu

menu:
	@clear
	@echo "$(CYAN)╔════════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(CYAN)║         YOUR SINATRA APP - Development Commands                ║$(NC)"
	@echo "$(CYAN)╚════════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@echo "$(GREEN)Development:$(NC)"
	@echo "  $(YELLOW)make test$(NC)              - Run tests"
	@echo "  $(YELLOW)make run$(NC)               - Run app locally (port 4567)"
	@echo "  $(YELLOW)make build$(NC)             - Build for production"
	@echo ""
	@echo "$(GREEN)Docker:$(NC)"
	@echo "  $(YELLOW)make docker-build$(NC)      - Build Docker image"
	@echo "  $(YELLOW)make docker-run$(NC)        - Run Docker container locally"
	@echo ""
	@echo "$(GREEN)Deployment:$(NC)"
	@echo "  $(YELLOW)make kamal-secrets-setup$(NC) - Generate .kamal/secrets file"
	@echo "  $(YELLOW)make deploy$(NC)            - Deploy to production with Kamal"
	@echo ""
	@echo "$(GREEN)Utilities:$(NC)"
	@echo "  $(YELLOW)make clean$(NC)             - Clean temporary files"
	@echo "  $(YELLOW)make help$(NC)              - Show this menu"
	@echo ""
	@read -p "Press Enter to continue..."

test:
	@echo "$(BLUE)→ Running tests...$(NC)"
	bundle exec rspec

build:
	@echo "$(BLUE)→ Building application...$(NC)"
	bundle install --without development test

run:
	@echo "$(BLUE)→ Starting development server on http://localhost:4567$(NC)"
	bundle exec rackup -o 0.0.0.0 -p 4567

docker-build:
	@echo "$(BLUE)→ Building Docker image...$(NC)"
	docker build -t your-app-name .

docker-run:
	@echo "$(BLUE)→ Running Docker container on http://localhost:4567$(NC)"
	docker run -it -p 4567:4567 --env-file .env your-app-name bundle exec rackup -o 0.0.0.0 -p 4567

kamal-secrets-setup:
	@echo "$(BLUE)→ Generating .kamal/secrets file...$(NC)"
	@mkdir -p .kamal
	@echo "KAMAL_REGISTRY_PASSWORD=\$$KAMAL_REGISTRY_PASSWORD" > .kamal/secrets
	@echo "DOPPLER_TOKEN=\$$DOPPLER_TOKEN" >> .kamal/secrets
	@chmod 600 .kamal/secrets
	@echo "$(GREEN)✓ .kamal/secrets file created$(NC)"
	@echo "$(YELLOW)Note: Doppler will substitute variables at runtime$(NC)"

deploy:
	@echo "$(BLUE)→ Deploying to production with Kamal...$(NC)"
	@echo "$(YELLOW)Checking if containers need cleanup...$(NC)"
	@if ssh root@YOUR_SERVER_IP 'docker ps -q --filter name=your-app-name' | grep -q .; then \
		echo "$(YELLOW)→ Stopping existing containers to free ports...$(NC)"; \
		ssh root@YOUR_SERVER_IP 'docker ps -q --filter name=your-app-name | xargs -r docker stop'; \
		sleep 2; \
	fi
	@echo "$(BLUE)→ Running Kamal deployment...$(NC)"
	doppler run --project your-project --config prd -- kamal deploy
	@echo "$(GREEN)✓ Deployment complete!$(NC)"

clean:
	@echo "$(BLUE)→ Cleaning temporary files...$(NC)"
	rm -rf tmp/
	rm -rf log/*.log
	@echo "$(GREEN)✓ Clean complete$(NC)"
```

**Important Customizations:**
- Replace `YOUR_SERVER_IP` with your actual server IP
- Replace `your-app-name` with your application name
- Replace `your-dockerhub-username` with your Docker Hub username
- Replace `your-project` with your Doppler project name
- Adjust port numbers if not using 4567

## Step 5: Setup Doppler Secrets

1. **Create Doppler Project:**
   ```bash
   doppler projects create your-project
   ```

2. **Create Production Config:**
   ```bash
   doppler configs create prd --project your-project
   ```

3. **Add Secrets:**
   ```bash
   # GitHub Container Registry PAT (for pulling/pushing images)
   doppler secrets set KAMAL_REGISTRY_PASSWORD="ghp_your_github_token" --project your-project --config prd

   # Doppler service token (for container runtime)
   doppler secrets set DOPPLER_TOKEN="dp.st.your_service_token" --project your-project --config prd

   # Add your application secrets
   doppler secrets set DATABASE_URL="your_database_url" --project your-project --config prd
   doppler secrets set SECRET_KEY_BASE="your_secret_key" --project your-project --config prd
   ```

## Step 6: Create .env.example

Create `.env.example` for local development reference:

```bash
# Application
RACK_ENV=development

# Database
DATABASE_URL=postgres://localhost/yourapp_dev

# Secrets (use Doppler in production)
SECRET_KEY_BASE=your-secret-key-here

# Doppler (for local testing with Doppler)
DOPPLER_TOKEN=dp.st.dev.your_token_here
```

## Step 7: Update .gitignore

```
.env
.env.*
!.env.example
.kamal/secrets
tmp/
log/
coverage/
.bundle/
vendor/bundle/
```

## Step 8: Initial Kamal Setup

1. **Generate secrets file:**
   ```bash
   make kamal-secrets-setup
   ```

2. **Initialize Kamal on server:**
   ```bash
   doppler run --project your-project --config prd -- kamal setup
   ```

3. **Deploy application:**
   ```bash
   make deploy
   ```

## Step 9: Add Health Check Endpoint

Add a health check endpoint to your Sinatra app (`app.rb`):

```ruby
require 'sinatra'

get '/health' do
  content_type :json
  { status: 'ok', timestamp: Time.now.to_i }.to_json
end

# Your other routes here
get '/' do
  "Hello World!"
end
```

## Common Commands

```bash
# Show menu
make

# Local development
make run

# Run tests
make test

# Build and test Docker image locally
make docker-build
make docker-run

# Deploy to production
make deploy

# View logs
doppler run -- kamal app logs --tail

# SSH into container
doppler run -- kamal app exec -i bash

# Rollback deployment
doppler run -- kamal rollback
```

## Deployment Workflow

1. **Development:**
   ```bash
   make run
   # Make changes, test locally
   ```

2. **Test in Docker:**
   ```bash
   make docker-build
   make docker-run
   ```

3. **Commit changes:**
   ```bash
   git add .
   git commit -m "Add feature"
   ```

4. **Deploy:**
   ```bash
   make deploy
   ```

## Critical Rules for Kamal + Doppler

### Rule 1: Always Use Doppler for Secrets
- **NEVER** commit secrets to git
- All production secrets MUST be in Doppler
- The Dockerfile MUST include Doppler CLI installation
- The CMD MUST use `doppler run --`

### Rule 2: Secrets File Required
- Kamal requires `.kamal/secrets` even with Doppler
- Use variable substitution format: `$VAR_NAME`
- Generate with `make kamal-secrets-setup`
- Never commit `.kamal/secrets` to git

### Rule 3: Port Mapping
- Container port (e.g., 4567) must match application port
- Traefik routes external traffic (80/443) to container
- Update `config/deploy.yml` if changing ports

### Rule 4: Cleanup Before Deploy
- The Makefile includes automatic cleanup
- Stops old containers before deploying
- Prevents port conflicts

## Troubleshooting

**Container won't start:**
```bash
doppler run -- kamal app logs
```

**Port already in use:**
```bash
ssh root@YOUR_SERVER_IP 'docker ps -a'
ssh root@YOUR_SERVER_IP 'docker stop $(docker ps -q --filter name=your-app-name)'
```

**Secrets not loading:**
```bash
# Verify Doppler token
doppler run -- echo "Token works"

# Check secrets file
cat .kamal/secrets
```

**Health check failing:**
- Ensure `/health` endpoint returns 200 OK
- Check container logs: `doppler run -- kamal app logs`
- Verify port mapping in `config/deploy.yml`

## File Checklist

Before deploying, verify these files exist:

- [ ] `Dockerfile` with Doppler integration
- [ ] `config/deploy.yml` with correct IPs and domains
- [ ] `Makefile` with deployment commands
- [ ] `.kamal/secrets` generated (not committed)
- [ ] `.dockerignore` to exclude unnecessary files
- [ ] `.gitignore` with `.kamal/secrets` excluded
- [ ] `.env.example` documenting required variables
- [ ] Health check endpoint (`/health`)
- [ ] Doppler project and secrets configured

## Next Steps

1. Customize the configurations for your specific app
2. Set up Doppler project and secrets
3. Test Docker build locally
4. Run initial `kamal setup`
5. Deploy with `make deploy`
6. Monitor logs and verify deployment
7. Set up monitoring/alerting as needed

## Additional Resources

- [Kamal Documentation](https://kamal-deploy.org/)
- [Doppler Documentation](https://docs.doppler.com/)
- [Sinatra Documentation](http://sinatrarb.com/)
- Void Chronicles project as reference implementation

---

**Note:** This guide is based on the deployment patterns from the Void Chronicles project. Adapt server IPs, domains, ports, and application-specific settings as needed for your Sinatra application.
