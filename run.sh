#!/bin/bash

set -e

echo "🚀 Starting Void Reavers SSH Reader..."
echo "====================================="

# Check if binary exists
if [ ! -f "./void-reader" ]; then
    echo "❌ Binary not found. Building first..."
    ./build.sh
fi

# Check if book content exists
if [ ! -d "./book1_void_reavers" ]; then
    echo "❌ Book content not found at ./book1_void_reavers/"
    echo "Please ensure the book directory exists and contains chapter files."
    exit 1
fi

# Check if SSH key exists
if [ ! -f "./.ssh/id_ed25519" ]; then
    echo "❌ SSH host key not found. Building first..."
    ./build.sh
fi

echo "📚 Book: Void Reavers"
echo "🌐 Server: localhost:23234"  
echo "🔑 SSH Key: .ssh/id_ed25519"
echo "💾 Data Dir: .void_reader_data/"
echo ""
echo "🎯 To connect: ssh localhost -p 23234"
echo ""
echo "Starting server..."
echo ""

# Start the server
./void-reader