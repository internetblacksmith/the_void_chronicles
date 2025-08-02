#!/bin/bash

# Convert LaTeX book to EPUB format for Kindle

echo "📚 Void Reavers EPUB Converter"
echo "============================="
echo ""

# Check if Ruby is installed
if ! command -v ruby &> /dev/null; then
    echo "❌ Ruby is not installed. Please install Ruby to use this converter."
    echo "   Ubuntu/Debian: sudo apt install ruby"
    echo "   macOS: brew install ruby"
    exit 1
fi

# Check if zip is installed
if ! command -v zip &> /dev/null; then
    echo "❌ zip is not installed. Please install zip to create EPUB files."
    echo "   Ubuntu/Debian: sudo apt install zip"
    echo "   macOS: Should be pre-installed"
    exit 1
fi

# Run the Ruby converter
ruby convert_to_epub.rb

echo ""
echo "📖 EPUB Conversion Tips:"
echo "========================"
echo ""
echo "📧 Send to Kindle via Email:"
echo "   1. Attach the .epub file to an email"
echo "   2. Send to your Kindle email address (yourname@kindle.com)"
echo "   3. Use 'Convert' as the subject line for older Kindles"
echo ""
echo "📱 Send via Kindle App:"
echo "   1. Open the Kindle app on your phone/tablet"
echo "   2. Use 'Send to Kindle' feature"
echo "   3. Select the EPUB file"
echo ""
echo "🔌 Transfer via USB:"
echo "   1. Connect Kindle to computer via USB"
echo "   2. Copy EPUB to 'Documents' folder"
echo "   3. Safely eject Kindle"
echo ""
echo "💡 Note: Newer Kindles (2022+) support EPUB directly."
echo "   Older models may need conversion to MOBI format."