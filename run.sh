#!/bin/bash

set -e

echo "🚀 Starting Void Reavers SSH Reader..."
echo "====================================="

# Check if binary exists
if [ ! -f "./ssh-reader/void-reader" ]; then
    echo "❌ Binary not found. Building first..."
    ./build.sh
fi

# Check if book content exists
if [ ! -d "./book1_void_reavers_source" ]; then
    echo "❌ Book content not found at ./book1_void_reavers_source/"
    echo "Please ensure the book directory exists and contains chapter files."
    exit 1
fi

# Check if SSH key exists
if [ ! -f "./.ssh/id_ed25519" ]; then
    echo "❌ SSH host key not found. Building first..."
    ./build.sh
fi

echo "📚 Book: Void Reavers"
echo "🔑 SSH Key: .ssh/id_ed25519"
echo "💾 Data Dir: .void_reader_data/"
echo ""

# Set ports for local development
export PORT=8080  # HTTP port for local dev (can't use 80 without sudo)
export SSH_PORT=23234

echo "🌐 HTTP Server: http://localhost:8080"
echo "🚀 SSH Server: localhost:23234"
echo ""
echo "🎯 To connect: ssh localhost -p 23234"
echo "🔑 Password: Amigos4Life!"
echo ""
echo "Starting servers..."
echo ""

# Start the server from the project root so it can find book files
cd ssh-reader && ./void-reader