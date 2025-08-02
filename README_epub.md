# EPUB Converter for Kindle 📚

Convert your LaTeX book "Void Reavers" to EPUB format for reading on Kindle devices and apps.

## Features

- 📖 **Direct LaTeX to EPUB conversion** - No intermediate steps needed
- 🎨 **Kindle-optimized formatting** - Proper styling for e-readers
- 📱 **Complete metadata** - Title, author, description included
- 🔖 **Full table of contents** - Easy chapter navigation
- ✨ **Clean HTML output** - Properly formatted text with italics and bold

## Quick Start

```bash
# Run the converter
./convert_to_epub.sh

# Or directly with Ruby
ruby convert_to_epub.rb
```

The EPUB file will be created at: `book1_void_reavers/void_reavers.epub`

## Transfer to Kindle

### Method 1: Email (Recommended)
1. Find your Kindle email address (Settings → Your Account → Send-to-Kindle Email)
2. Email the EPUB file as an attachment to your Kindle address
3. For older Kindles, use "Convert" as the subject line

### Method 2: Send to Kindle App
1. Download the "Send to Kindle" app on your computer
2. Drag and drop the EPUB file into the app
3. Select your device and send

### Method 3: USB Transfer
1. Connect your Kindle via USB cable
2. Copy `void_reavers.epub` to the `Documents` folder
3. Safely eject your Kindle

### Method 4: Kindle Mobile App
1. Save the EPUB to your phone/tablet
2. Open the file and select "Open with Kindle"
3. The book will sync to all your devices

## Kindle Compatibility

- **Kindle (2022 and newer)**: Native EPUB support ✅
- **Kindle Paperwhite (2021+)**: Native EPUB support ✅
- **Kindle Oasis**: Native EPUB support ✅
- **Older Kindles**: Email with "Convert" subject line to convert to MOBI

## EPUB Structure

The converter creates a complete EPUB 3.0 file with:
```
void_reavers.epub
├── mimetype
├── META-INF/
│   └── container.xml
└── OEBPS/
    ├── content.opf      # Book metadata and manifest
    ├── toc.ncx          # Navigation (EPUB 2.0 compatibility)
    ├── toc.xhtml        # Navigation (EPUB 3.0)
    ├── cover.xhtml      # Title page
    ├── chapter01-20.xhtml # All 20 chapters
    └── css/
        └── style.css    # Kindle-optimized styles
```

## Customization

### Modify Book Metadata
Edit the `load_book_info` method in `convert_to_epub.rb`:
```ruby
def load_book_info(book_dir)
  {
    title: "Your Book Title",
    author: "Your Name",
    description: "Your book description",
    # ... other metadata
  }
end
```

### Adjust Styling
Edit the CSS in the `generate_css` method for different formatting.

## Troubleshooting

### EPUB won't open
- Ensure you have `zip` installed: `sudo apt install zip`
- Check that all chapter files exist in the source directory

### Formatting issues on Kindle
- The converter uses Kindle-compatible CSS
- Complex formatting may be simplified for e-reader compatibility

### Missing chapters
- Verify all chapter files (chapter01.tex through chapter20.tex) exist
- Check file permissions

## Technical Details

The converter:
1. Reads LaTeX files directly (no Markdown conversion needed)
2. Converts LaTeX commands to HTML:
   - `\textit{}` → `<em></em>`
   - `\textbf{}` → `<strong></strong>`
   - Smart quotes and em-dashes preserved
3. Creates proper EPUB structure with all required files
4. Generates navigation files for both EPUB 2.0 and 3.0
5. Includes Kindle-optimized CSS for best reading experience

## Requirements

- Ruby (any recent version)
- `zip` command-line tool
- The LaTeX source files

---

Enjoy reading "Void Reavers" on your Kindle! 🚀📱