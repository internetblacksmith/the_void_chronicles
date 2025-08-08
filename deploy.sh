#!/bin/bash

set -e

INSTALL_DIR="/opt/void-reader"
SERVICE_USER="voidreader"
SERVICE_FILE="ssh-reader/systemd/void-reader.service"

echo "🚀 Deploying Void Reavers SSH Reader..."
echo "======================================"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "❌ Please run as root (use sudo)"
    exit 1
fi

# Build the application first
echo "🔨 Building application..."
sudo -u $USER ./build.sh

# Create service user if it doesn't exist
if ! id "$SERVICE_USER" &>/dev/null; then
    echo "👤 Creating service user: $SERVICE_USER"
    useradd --system --home-dir $INSTALL_DIR --create-home --shell /bin/false $SERVICE_USER
fi

# Create installation directory
echo "📁 Creating installation directory: $INSTALL_DIR"
mkdir -p $INSTALL_DIR
mkdir -p $INSTALL_DIR/.ssh
mkdir -p $INSTALL_DIR/.void_reader_data

# Copy application files
echo "📋 Copying application files..."
cp ssh-reader/void-reader $INSTALL_DIR/
cp -r book1_void_reavers $INSTALL_DIR/

# Copy or generate SSH key
if [ -f ".ssh/id_ed25519" ]; then
    cp .ssh/id_ed25519* $INSTALL_DIR/.ssh/
else
    echo "🔑 Generating SSH host key..."
    ssh-keygen -t ed25519 -f $INSTALL_DIR/.ssh/id_ed25519 -N "" -C "void-reader-production"
fi

# Set proper ownership
echo "🔒 Setting file permissions..."
chown -R $SERVICE_USER:$SERVICE_USER $INSTALL_DIR
chmod 755 $INSTALL_DIR/void-reader
chmod 600 $INSTALL_DIR/.ssh/id_ed25519
chmod 644 $INSTALL_DIR/.ssh/id_ed25519.pub

# Install systemd service
if [ -f "$SERVICE_FILE" ]; then
    echo "⚙️  Installing systemd service..."
    cp $SERVICE_FILE /etc/systemd/system/
    systemctl daemon-reload
    systemctl enable void-reader
    
    # Stop service if running
    if systemctl is-active --quiet void-reader; then
        echo "🛑 Stopping existing service..."
        systemctl stop void-reader
        sleep 2
    fi
    
    # Start service
    echo "▶️  Starting service..."
    systemctl start void-reader
    
    # Check status
    sleep 2
    if systemctl is-active --quiet void-reader; then
        echo "✅ Service started successfully"
        echo ""
        echo "📊 Service status:"
        systemctl status void-reader --no-pager -l
    else
        echo "❌ Service failed to start"
        echo "📋 Check logs with: journalctl -u void-reader -f"
        exit 1
    fi
else
    echo "⚠️  Systemd service file not found. Running manually..."
    echo "To start manually: sudo -u $SERVICE_USER $INSTALL_DIR/void-reader"
fi

echo ""
echo "🎉 Deployment complete!"
echo ""
echo "📡 Server is running on localhost:23234"
echo "🔌 Connect with: ssh localhost -p 23234"
echo "📋 View logs: journalctl -u void-reader -f"
echo "⚙️  Manage service:"
echo "   sudo systemctl start void-reader"
echo "   sudo systemctl stop void-reader"
echo "   sudo systemctl restart void-reader"
echo "   sudo systemctl status void-reader"
echo ""
echo "📁 Installation directory: $INSTALL_DIR"
echo "💾 User data: $INSTALL_DIR/.void_reader_data"
echo ""