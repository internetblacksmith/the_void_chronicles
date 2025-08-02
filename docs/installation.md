# Installation Guide

Complete installation instructions for the Void Reavers SSH Reader across different environments and use cases.

## 📋 System Requirements

### Minimum Requirements
- **Operating System**: Linux, macOS, Windows (with WSL)
- **Go Version**: 1.21 or later
- **RAM**: 50MB available memory
- **Storage**: 100MB free disk space
- **Network**: 1 available port (default: 23234)
- **Terminal**: UTF-8 compatible terminal emulator

### Recommended Requirements
- **CPU**: 1 core at 1GHz or better
- **RAM**: 100MB available memory
- **Storage**: 500MB free disk space (for multiple books and user data)
- **Network**: Stable network connection for SSH access
- **Terminal**: Modern terminal with color support

## 🛠️ Installation Methods

### Method 1: Local Development (Recommended for Testing)

#### Step 1: Install Go

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install golang-go
```

**CentOS/RHEL/Fedora:**
```bash
sudo dnf install golang  # Fedora
sudo yum install golang  # CentOS/RHEL
```

**macOS:**
```bash
# Using Homebrew
brew install go

# Or download from https://golang.org/dl/
```

**Windows:**
- Download installer from [golang.org](https://golang.org/dl/)
- Or use WSL with Linux instructions

#### Step 2: Verify Go Installation

```bash
go version
# Should output: go version go1.21.x linux/amd64 (or similar)
```

#### Step 3: Get the Source Code

```bash
# Navigate to your development directory
cd ~/projects  # or wherever you keep projects

# If you have the source already
cd void-reavers-reader

# Verify book content exists
ls book1_void_reavers/
# Should show chapter files and markdown directory
```

#### Step 4: Build and Install

```bash
# Make build script executable (if needed)
chmod +x build.sh run.sh

# Build the application
./build.sh
```

#### Step 5: Test the Installation

```bash
# Start the server
./run.sh

# In another terminal, test connection
ssh localhost -p 23234
```

### Method 2: Production Deployment

#### Step 1: Prepare the System

**Create dedicated user:**
```bash
sudo useradd --system --home-dir /opt/void-reader --create-home voidreader
```

**Install dependencies:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go openssh-client

# CentOS/RHEL/Fedora
sudo dnf install golang openssh-clients
```

#### Step 2: Deploy Application

```bash
# Run the deployment script (requires sudo)
sudo ./deploy.sh
```

This script will:
- Build the application
- Create service user and directories
- Install systemd service
- Generate SSH keys
- Set proper permissions
- Start the service

#### Step 3: Verify Production Installation

```bash
# Check service status
sudo systemctl status void-reader

# Test connection
ssh localhost -p 23234

# View logs
sudo journalctl -u void-reader -f
```

### Method 3: Docker Installation

#### Step 1: Install Docker

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install docker.io docker-compose
sudo systemctl enable --now docker
sudo usermod -aG docker $USER
```

**CentOS/RHEL/Fedora:**
```bash
sudo dnf install docker docker-compose
sudo systemctl enable --now docker
sudo usermod -aG docker $USER
```

#### Step 2: Build and Run with Docker

```bash
# Build and start with Docker Compose
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f void-reader
```

#### Step 3: Test Docker Installation

```bash
# Connect to the containerized service
ssh localhost -p 23234
```

### Method 4: Manual Installation

If the automated scripts don't work for your environment:

#### Step 1: Download Dependencies

```bash
go mod download
```

#### Step 2: Generate SSH Keys

```bash
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ""
```

#### Step 3: Create Data Directory

```bash
mkdir -p .void_reader_data
```

#### Step 4: Build Application

```bash
go build -ldflags="-s -w" -o void-reader
chmod +x void-reader
```

#### Step 5: Start Application

```bash
./void-reader
```

## 🔧 Post-Installation Configuration

### Configure SSH Access

#### Allow External Connections

By default, the server only accepts local connections. To allow remote access:

1. Edit `main.go`:
```go
const (
    host = "0.0.0.0"    // Change from "localhost"
    port = "23234"
)
```

2. Rebuild:
```bash
./build.sh
```

3. Configure firewall:
```bash
# Ubuntu/Debian (ufw)
sudo ufw allow 23234/tcp

# CentOS/RHEL/Fedora (firewalld)
sudo firewall-cmd --permanent --add-port=23234/tcp
sudo firewall-cmd --reload

# iptables
sudo iptables -A INPUT -p tcp --dport 23234 -j ACCEPT
```

#### Change Default Port

1. Edit `main.go`:
```go
const (
    host = "localhost"
    port = "2222"       // Your preferred port
)
```

2. Rebuild and restart

### Configure Book Content

#### Add Multiple Books

1. Create book directories:
```bash
mkdir book2_shadow_dancers
mkdir book3_quantum_academy
```

2. Add content in expected format:
```
book2_shadow_dancers/
├── markdown/
│   ├── chapter01.md
│   ├── chapter02.md
│   └── ...
└── chapter01.tex  # LaTeX fallback
```

3. Update the book loader in `book.go` if needed

#### Customize Book Loading

The application automatically loads from:
1. `book1_void_reavers/markdown/` (preferred)
2. `book1_void_reavers/*.tex` (fallback)

To add books or change paths, modify the `LoadBook()` function in `book.go`.

## 🔒 Security Configuration

### Production Security

#### SSH Host Key Security

```bash
# Backup your SSH key
cp .ssh/id_ed25519 /secure/backup/location/

# Set restrictive permissions
chmod 600 .ssh/id_ed25519
chmod 644 .ssh/id_ed25519.pub
```

#### User Data Protection

```bash
# Ensure user data is protected
chmod 755 .void_reader_data
chown -R voidreader:voidreader .void_reader_data  # Production
```

#### Firewall Configuration

Restrict access to trusted networks:
```bash
# Allow only specific IP range
sudo ufw allow from 192.168.1.0/24 to any port 23234

# Or specific IPs
sudo ufw allow from 203.0.113.10 to any port 23234
```

### SSL/TLS (Advanced)

For additional security, you can put the SSH server behind a reverse proxy with SSL:

```nginx
# Nginx stream configuration
stream {
    upstream void_reader {
        server localhost:23234;
    }
    
    server {
        listen 443;
        proxy_pass void_reader;
        proxy_timeout 1s;
        proxy_responses 1;
    }
}
```

## 🧪 Testing Installation

### Automated Tests

```bash
# Run Go tests (if any exist)
go test -v ./...

# Test SSH connection
timeout 10 ssh -o ConnectTimeout=5 localhost -p 23234 echo "Connection test"
```

### Manual Testing Checklist

- [ ] Server starts without errors
- [ ] SSH connection succeeds
- [ ] Main menu displays correctly
- [ ] Can navigate to different chapters
- [ ] Progress is saved and restored
- [ ] Bookmarks work correctly
- [ ] Multiple users maintain separate progress
- [ ] Server handles disconnections gracefully

## 🔄 Updating

### Update Application

```bash
# Pull latest changes (if using git)
git pull origin main

# Rebuild
./build.sh

# For production, redeploy
sudo ./deploy.sh
```

### Update Dependencies

```bash
# Update Go modules
go get -u ./...
go mod tidy

# Rebuild
./build.sh
```

## 🗑️ Uninstallation

### Remove Local Installation

```bash
# Stop the application (Ctrl+C if running)

# Remove built files
rm -f void-reader
rm -rf .ssh
rm -rf .void_reader_data

# Remove Go modules cache (optional)
go clean -modcache
```

### Remove Production Installation

```bash
# Stop and disable service
sudo systemctl stop void-reader
sudo systemctl disable void-reader

# Remove service file
sudo rm /etc/systemd/system/void-reader.service
sudo systemctl daemon-reload

# Remove installation directory
sudo rm -rf /opt/void-reader

# Remove service user
sudo userdel voidreader
```

### Remove Docker Installation

```bash
# Stop and remove containers
docker-compose down

# Remove images
docker-compose down --rmi all

# Remove volumes (warning: this deletes user data)
docker-compose down --volumes
```

## 🆘 Troubleshooting Installation

### Common Issues

#### "Go not found"
```bash
# Check Go installation
which go
go version

# Add Go to PATH if needed
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

#### "Permission denied" on build
```bash
# Make scripts executable
chmod +x build.sh run.sh deploy.sh

# Check file ownership
ls -la
```

#### "Port already in use"
```bash
# Find what's using port 23234
sudo netstat -tlnp | grep 23234
sudo lsof -i :23234

# Kill the process or change port in main.go
```

#### "Book content not found"
```bash
# Verify book directory structure
ls -la book1_void_reavers/
ls -la book1_void_reavers/markdown/

# Check file permissions
chmod -R 644 book1_void_reavers/
chmod 755 book1_void_reavers/
```

### Getting Help

If you encounter issues:

1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Review server logs:
   ```bash
   # Local development
   ./void-reader  # Check console output
   
   # Production
   sudo journalctl -u void-reader -f
   
   # Docker
   docker-compose logs -f
   ```
3. Test with minimal configuration
4. Check system requirements
5. Verify network connectivity

---

**Next Steps:**
- Installation complete? See the [User Guide](user-guide.md)
- Need to configure advanced features? Check [Configuration Guide](configuration.md)
- Ready for production? Read the [Deployment Guide](deployment.md)