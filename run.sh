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

# Load environment variables from .env if it exists
if [ -f ".env" ]; then
    echo "📋 Loading configuration from .env"
    export $(grep -v '^#' .env | xargs)
else
    echo "⚠️  No .env file found, using defaults"
    echo "   (Copy .env.example to .env to customize)"
fi

# Display connection info (read from env or defaults)
HTTP_PORT=${HTTP_PORT:-8080}
SSH_PORT=${SSH_PORT:-2222}

echo "🌐 HTTP Server: http://localhost:${HTTP_PORT}"
echo "🚀 SSH Server: localhost:${SSH_PORT}"
echo ""
echo "🎯 To connect: ssh localhost -p ${SSH_PORT}"
echo "🔑 Password: ${SSH_PASSWORD:-Amigos4Life!}"
echo ""
echo "Starting servers..."
echo ""

# Start the server from the project root so it can find book files
cd ssh-reader && ./void-reader