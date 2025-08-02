#!/bin/bash

# Simple PDF conversion wrapper for Void Reavers

echo "Converting Void Reavers to PDF..."

# Check if markdown directory exists
if [ ! -d "markdown" ]; then
    echo "Error: markdown directory not found!"
    echo "Please run convert_to_md.rb first."
    exit 1
fi

# Create pdf directory
mkdir -p pdf

# Try pandoc first
if command -v pandoc &> /dev/null; then
    echo "Using pandoc for conversion..."
    
    # Convert complete book
    pandoc markdown/void_reavers_complete.md \
        -o pdf/void_reavers_complete.pdf \
        --toc \
        --standalone \
        -V geometry:margin=1in \
        -V fontsize=11pt \
        -V title="Void Reavers" \
        -V author="Captain J. Starwind" \
        -V subtitle="A Tale of Space Pirates and Cosmic Plunder"
    
    echo "PDF created: pdf/void_reavers_complete.pdf"
    
# Try wkhtmltopdf as fallback
elif command -v wkhtmltopdf &> /dev/null; then
    echo "Using wkhtmltopdf for conversion..."
    
    # First convert to HTML
    pandoc markdown/void_reavers_complete.md -o /tmp/void_reavers.html --standalone
    
    # Then to PDF
    wkhtmltopdf --title "Void Reavers" /tmp/void_reavers.html pdf/void_reavers_complete.pdf
    
    rm /tmp/void_reavers.html
    echo "PDF created: pdf/void_reavers_complete.pdf"
    
else
    echo "Error: No PDF converter found!"
    echo "Please install pandoc or wkhtmltopdf"
    exit 1
fi
