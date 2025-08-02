#!/bin/bash

# Amazon KDP PDF Generator for Void Reavers

echo "📚 Amazon KDP PDF Generator"
echo "==========================="
echo ""
echo "This tool creates a print-ready PDF for Amazon KDP self-publishing."
echo ""

# Check if Ruby is installed
if ! command -v ruby &> /dev/null; then
    echo "❌ Ruby is not installed. Please install Ruby first."
    echo "   Ubuntu/Debian: sudo apt install ruby"
    echo "   macOS: brew install ruby"
    exit 1
fi

# Check if pdflatex is installed
if ! command -v pdflatex &> /dev/null; then
    echo "❌ pdflatex is not installed. Please install LaTeX first."
    echo ""
    echo "Installation instructions:"
    echo ""
    echo "🐧 Ubuntu/Debian:"
    echo "   sudo apt update"
    echo "   sudo apt install texlive-full"
    echo "   (Large download, ~5GB, but includes everything needed)"
    echo ""
    echo "🐧 Fedora/RHEL:"
    echo "   sudo dnf install texlive-scheme-full"
    echo ""
    echo "🍎 macOS:"
    echo "   brew install --cask mactex"
    echo "   (Or download from https://www.tug.org/mactex/)"
    echo ""
    echo "🪟 Windows:"
    echo "   Download MiKTeX from https://miktex.org/"
    echo ""
    exit 1
fi

# Run the Ruby generator
ruby create_kdp_pdf.rb

# If successful, provide additional tips
if [ $? -eq 0 ]; then
    echo ""
    echo "📖 Amazon KDP Publishing Tips:"
    echo "=============================="
    echo ""
    echo "1️⃣ Interior File:"
    echo "   - Review the generated PDF carefully"
    echo "   - Check page numbers, margins, and formatting"
    echo "   - Ensure no blank pages unless intentional"
    echo "   - Order a printed proof before publishing"
    echo ""
    echo "2️⃣ Cover Design:"
    echo "   - Use the dimensions in cover_specs.txt"
    echo "   - Design in CMYK color mode for best results"
    echo "   - Leave space for ISBN barcode (added by Amazon)"
    echo "   - Consider hiring a designer on Fiverr/99designs"
    echo ""
    echo "3️⃣ KDP Upload Process:"
    echo "   - Go to kdp.amazon.com"
    echo "   - Create a new paperback book"
    echo "   - Upload interior PDF"
    echo "   - Use Cover Creator or upload your own"
    echo "   - Set pricing and distribution"
    echo ""
    echo "4️⃣ ISBN Options:"
    echo "   - Free ISBN from Amazon (limited to Amazon)"
    echo "   - Buy your own ISBN for wider distribution"
    echo "   - ISBNs from Bowker (US) or your country's agency"
    echo ""
    echo "5️⃣ Pricing Strategy:"
    echo "   - Check comparable books in your genre"
    echo "   - Factor in printing costs (KDP shows this)"
    echo "   - Consider starting lower to gain reviews"
    echo "   - 35% or 70% royalty options available"
    echo ""
    echo "🚀 Good luck with your publishing journey!"
fi