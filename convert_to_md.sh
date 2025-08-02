#!/bin/bash

# Function to convert a book directory
convert_book() {
    local book_dir="$1"
    echo "Converting $book_dir..."
    echo "============================="
    
    # Create markdown directory if it doesn't exist
    mkdir -p "${book_dir}/markdown"
    
    # Convert each chapter file
    for i in {01..20}; do
        if [ -f "${book_dir}/chapter${i}.tex" ]; then
            echo "Converting chapter${i}.tex..."
            
            # Basic conversion using sed
            sed -e 's/\\chapter{\([^}]*\)}/# \1/g' \
                -e 's/\\textit{\([^}]*\)}/*\1*/g' \
                -e 's/\\textbf{\([^}]*\)}/**\1**/g' \
                -e "s/\`\`\([^']*\)''/\"\1\"/g" \
                -e "s/\`\([^']*\)'/'\1'/g" \
                -e 's/---/—/g' \
                -e 's/\\%/%/g' \
                -e 's/\\\$/$$/g' \
                -e 's/\\&/\&/g' \
                -e 's/\\_/_/g' \
                -e 's/\\#/#/g' \
                "${book_dir}/chapter${i}.tex" > "${book_dir}/markdown/chapter${i}.md"
        fi
    done
    
    # Create complete book for this directory
    create_complete_book "$book_dir"
}

# Function to create complete book
create_complete_book() {
    local book_dir="$1"
    echo "Creating complete book for $book_dir..."
    {
        echo "# Void Reavers"
        echo ""
        echo "## A Tale of Space Pirates and Cosmic Plunder"
        echo ""
        echo "By Captain J. Starwind"
        echo ""
        echo "---"
        echo ""
        
        for i in {01..20}; do
            if [ -f "${book_dir}/markdown/chapter${i}.md" ]; then
                cat "${book_dir}/markdown/chapter${i}.md"
                echo ""
                echo "---"
                echo ""
            fi
        done
    } > "${book_dir}/markdown/void_reavers_complete.md"
    
    # Create README for this book
    cat > "${book_dir}/markdown/README.md" << EOF
# Void Reavers - Markdown Version

This directory contains the markdown version of the book "Void Reavers".

## Files

- \`void_reavers_complete.md\` - The complete book in a single file
- \`chapter01.md\` through \`chapter20.md\` - Individual chapter files

## About

This book tells the story of Captain Zara "Bloodhawk" Vega and the transformation of space piracy in a universe where humanity must prove itself worthy of the stars.

## Reading Order

The chapters should be read in numerical order from 01 to 20. The complete book file contains all chapters in the correct order.
EOF
}

echo "Converting LaTeX files to Markdown..."
echo "===================================="

# Check if we're in a book directory or the main directory
if [ -f "book.tex" ]; then
    # We're in a book directory
    convert_book "."
elif [ -d "book1_void_reavers" ]; then
    # We're in the main directory, convert all book directories
    for book_dir in book*; do
        if [ -d "$book_dir" ] && [ -f "${book_dir}/book.tex" ]; then
            convert_book "$book_dir"
        fi
    done
else
    echo "No book.tex found and no book directories detected."
    echo "Please run this script from a book directory or the main series directory."
    exit 1
fi

# Create complete book
echo "Creating complete book..."
{
    echo "# Void Reavers"
    echo ""
    echo "## A Tale of Space Pirates and Cosmic Plunder"
    echo ""
    echo "By Captain J. Starwind"
    echo ""
    echo "---"
    echo ""
    
    for i in {01..20}; do
        if [ -f "markdown/chapter${i}.md" ]; then
            cat "markdown/chapter${i}.md"
            echo ""
            echo "---"
            echo ""
        fi
    done
} > "markdown/void_reavers_complete.md"

# Create README
cat > "markdown/README.md" << EOF
# Void Reavers - Markdown Version

This directory contains the markdown version of the book "Void Reavers".

## Files

- \`void_reavers_complete.md\` - The complete book in a single file
- \`chapter01.md\` through \`chapter20.md\` - Individual chapter files

## About

This book tells the story of Captain Zara "Bloodhawk" Vega and the transformation of space piracy in a universe where humanity must prove itself worthy of the stars.

## Reading Order

The chapters should be read in numerical order from 01 to 20. The complete book file contains all chapters in the correct order.
EOF

echo ""
echo "Conversion complete!"
echo "Markdown files saved in: ./markdown/"