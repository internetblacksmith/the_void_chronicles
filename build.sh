#!/bin/bash

set -e

echo "🚀 Building Void Reavers SSH Reader..."
echo "=================================="

# Create .ssh directory if it doesn't exist
mkdir -p .ssh

# Generate SSH host key if it doesn't exist
if [ ! -f .ssh/id_ed25519 ]; then
    echo "🔑 Generating SSH host key..."
    ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N "" -C "void-reader-host-key"
    echo "✅ SSH host key generated"
fi

# Create data directory
mkdir -p .void_reader_data
echo "📁 Created data directory for user progress"

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod tidy

# Run tests if any exist
if ls *_test.go 1> /dev/null 2>&1; then
    echo "🧪 Running tests..."
    go test -v ./...
fi

# Build the application
echo "🔨 Building application..."
go build -ldflags="-s -w" -o void-reader

# Make sure the binary is executable
chmod +x void-reader

echo ""
echo "✅ Build complete!"
echo ""
echo "📚 Book content location: $(pwd)/book1_void_reavers/"
echo "🔑 SSH host key: $(pwd)/.ssh/id_ed25519"
echo "💾 User data: $(pwd)/.void_reader_data/"
echo "🚀 Binary: $(pwd)/void-reader"
echo ""
echo "To start the server, run:"
echo "  ./void-reader"
echo ""
echo "To connect, run:"
echo "  ssh localhost -p 23234"
echo ""