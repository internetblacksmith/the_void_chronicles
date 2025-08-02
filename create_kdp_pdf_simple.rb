#!/usr/bin/env ruby

require 'fileutils'
require 'erb'

class SimpleKDPPDFGenerator
  def initialize(base_dir = Dir.pwd)
    @base_dir = base_dir
    @book_dir = find_book_directory
  end

  def generate
    unless @book_dir
      puts "❌ No book directory found!"
      return
    end

    puts "\n📚 Simple Amazon KDP PDF Generator"
    puts "=================================="
    puts "\nThis creates KDP specifications and instructions for your book."
    puts ""
    
    # Get page count estimate
    page_count = estimate_page_count
    
    # Create output directory
    output_dir = File.join(@book_dir, 'kdp_output')
    FileUtils.mkdir_p(output_dir)
    
    # Generate cover specifications
    generate_cover_specs(page_count)
    
    # Generate formatted HTML for conversion
    generate_kdp_html(page_count)
    
    # Generate instructions
    generate_instructions(page_count)
    
    puts "\n✅ KDP preparation complete!"
    puts "\n📋 Files created in #{output_dir}:"
    puts "   - cover_specs.txt: Exact dimensions for your cover"
    puts "   - kdp_instructions.txt: Step-by-step publishing guide"
    puts "   - void_reavers_kdp.html: Formatted HTML (convert to PDF)"
    puts "\n🔄 To create PDF from HTML:"
    puts "   Option 1: Open HTML in Chrome/Firefox → Print → Save as PDF"
    puts "   Option 2: Use online converter like ilovepdf.com"
    puts "   Option 3: Use pandoc if installed: pandoc void_reavers_kdp.html -o void_reavers_kdp.pdf"
  end

  private

  def find_book_directory
    Dir.glob(File.join(@base_dir, 'book*')).select { |d| File.directory?(d) }.first
  end

  def estimate_page_count
    word_count = 0
    
    # Try markdown files first
    markdown_dir = File.join(@book_dir, 'markdown')
    if Dir.exist?(markdown_dir)
      Dir.glob(File.join(markdown_dir, 'chapter*.md')).each do |file|
        content = File.read(file)
        word_count += content.split.length
      end
    else
      # Fall back to LaTeX files
      Dir.glob(File.join(@book_dir, 'chapter*.tex')).each do |file|
        content = File.read(file)
        text = content.gsub(/\\[a-zA-Z]+\{[^}]*\}/, ' ').gsub(/\\[a-zA-Z]+/, ' ')
        word_count += text.split.length
      end
    end
    
    # Estimate pages (250 words per page for 6x9)
    estimated_pages = (word_count / 250.0).ceil
    
    puts "\n📊 Book Statistics:"
    puts "   Word count: ~#{word_count.to_s.reverse.gsub(/(\d{3})(?=\d)/, '\\1,').reverse}"
    puts "   Estimated pages: ~#{estimated_pages}"
    
    print "\nEnter actual/expected page count [#{estimated_pages}]: "
    input = gets.chomp
    input.empty? ? estimated_pages : input.to_i
  end

  def generate_cover_specs(page_count)
    # Standard 6x9 dimensions
    trim_width = 6.0
    trim_height = 9.0
    
    # Calculate spine width
    spine_width_white = page_count * 0.0025
    spine_width_cream = page_count * 0.002739
    
    # Calculate full cover dimensions with bleed
    bleed = 0.125
    cover_width_white = (trim_width * 2) + spine_width_white + (bleed * 2)
    cover_height = trim_height + (bleed * 2)
    cover_width_cream = (trim_width * 2) + spine_width_cream + (bleed * 2)
    
    specs = <<~SPECS
      Amazon KDP Cover Specifications
      ================================
      
      Book: Void Reavers
      Trim Size: 6" × 9" (15.24 × 22.86 cm)
      Page Count: #{page_count}
      
      COVER DIMENSIONS - WHITE PAPER:
      --------------------------------
      Total Width: #{'%.3f' % cover_width_white}" (#{'%.1f' % (cover_width_white * 25.4)} mm)
      Total Height: #{'%.3f' % cover_height}" (#{'%.1f' % (cover_height * 25.4)} mm)
      Spine Width: #{'%.3f' % spine_width_white}" (#{'%.1f' % (spine_width_white * 25.4)} mm)
      
      COVER DIMENSIONS - CREAM PAPER:
      --------------------------------
      Total Width: #{'%.3f' % cover_width_cream}" (#{'%.1f' % (cover_width_cream * 25.4)} mm)
      Total Height: #{'%.3f' % cover_height}" (#{'%.1f' % (cover_height * 25.4)} mm)
      Spine Width: #{'%.3f' % spine_width_cream}" (#{'%.1f' % (spine_width_cream * 25.4)} mm)
      
      COVER LAYOUT (left to right):
      1. Back Cover: 6" × 9"
      2. Spine: #{'%.3f' % spine_width_white}" × 9"
      3. Front Cover: 6" × 9"
      
      Plus 0.125" (3.2mm) bleed on all outer edges!
      
      TECHNICAL REQUIREMENTS:
      - Resolution: 300 DPI minimum (1800 × 2700 pixels for 6×9")
      - Color mode: CMYK preferred, RGB accepted
      - File format: PDF (flattened, no layers)
      - File size: Maximum 650 MB
      
      DESIGN REQUIREMENTS:
      - Bleed: 0.125" (3.2 mm) on all outer edges
      - Safe zone: Keep text 0.25" (6.4 mm) from trim edges
      - Spine text: Only if 79+ pages (you have #{page_count})
      - Spine margins: 0.0625" (1.6 mm) from spine edges
      - Barcode area: Reserve 2" × 1.2" on back cover (Amazon adds this)
      
      QUICK REFERENCE:
      Front cover safe area: 5.5" × 8.5"
      Back cover safe area: 5.5" × 8.5"
      #{page_count >= 79 ? "Spine text area: #{'%.3f' % (spine_width_white - 0.125)}\" × 8.5\"" : "Spine too narrow for text (needs 79+ pages)"}
      
      FREE COVER DESIGN TOOLS:
      1. Canva.com - Search "book cover 6x9"
      2. BookCoverMaker.com - Free templates
      3. Amazon Cover Creator - Built into KDP
      
      PROFESSIONAL DESIGNERS:
      1. 99designs.com - Contest ~$299+
      2. Fiverr.com - Search "book cover" $5-200
      3. GetCovers.com - Premade covers $45-300
      4. Reedsy.com - Pro designers $300-800
      
      COVER TEMPLATE:
      Download exact template from:
      https://kdp.amazon.com/cover-templates
    SPECS
    
    File.write(File.join(@book_dir, 'kdp_output', 'cover_specs.txt'), specs)
    puts "\n📐 Cover specifications saved!"
  end

  def generate_kdp_html(page_count)
    html = <<~HTML
      <!DOCTYPE html>
      <html>
      <head>
        <meta charset="UTF-8">
        <title>Void Reavers - KDP Edition</title>
        <style>
          /* KDP-optimized print styles */
          @page {
            size: 6in 9in;
            margin: 0.75in 0.5in;
          }
          
          @media print {
            body {
              margin: 0;
              font-size: 12pt;
            }
            .page-break {
              page-break-after: always;
            }
            .chapter {
              page-break-before: always;
            }
          }
          
          body {
            font-family: Georgia, 'Times New Roman', serif;
            font-size: 12pt;
            line-height: 1.6;
            text-align: justify;
            margin: 0 auto;
            max-width: 4.5in;
            padding: 1in;
          }
          
          h1 {
            font-size: 24pt;
            text-align: center;
            margin: 2em 0 1em 0;
            page-break-before: always;
          }
          
          h2 {
            font-size: 18pt;
            text-align: center;
            margin: 1.5em 0 1em 0;
          }
          
          p {
            text-indent: 0.25in;
            margin: 0 0 0.5em 0;
          }
          
          p.first {
            text-indent: 0;
          }
          
          .title-page {
            text-align: center;
            page-break-after: always;
            margin-top: 2in;
          }
          
          .copyright-page {
            font-size: 10pt;
            text-align: center;
            page-break-after: always;
          }
          
          .toc {
            page-break-after: always;
          }
          
          .toc-entry {
            text-indent: 0;
            margin: 0.5em 0;
          }
          
          em {
            font-style: italic;
          }
          
          strong {
            font-weight: bold;
          }
        </style>
      </head>
      <body>
        <!-- Title Page -->
        <div class="title-page">
          <h1 style="font-size: 36pt; margin-top: 0;">VOID REAVERS</h1>
          <h2 style="font-size: 20pt;">A Tale of Space Pirates and Cosmic Plunder</h2>
          <div style="margin-top: 3in;">
            <p style="font-size: 16pt; text-indent: 0;">Captain J. Starwind</p>
          </div>
          <div style="margin-top: 2in;">
            <p style="text-indent: 0;">Book One of The Void Chronicles</p>
          </div>
        </div>
        
        <!-- Copyright Page -->
        <div class="copyright-page">
          <p style="margin-top: 3in; text-indent: 0;">
            Copyright © 2024 by [Your Name]<br><br>
            All rights reserved.<br><br>
            ISBN: [Your ISBN]<br>
            First Edition<br><br>
            This is a work of fiction. Names, characters, places, and incidents
            either are the product of the author's imagination or are used
            fictitiously, and any resemblance to actual persons, living or dead,
            business establishments, events, or locales is entirely coincidental.<br><br>
            Published in the United States of America
          </p>
        </div>
        
        <!-- Table of Contents -->
        <div class="toc">
          <h1 style="page-break-before: avoid;">Table of Contents</h1>
    HTML
    
    # Add TOC entries
    (1..20).each do |num|
      html += "      <p class=\"toc-entry\">Chapter #{num}</p>\n"
    end
    
    html += <<~HTML
        </div>
        
        <!-- Chapters -->
    HTML
    
    # Add chapter content
    markdown_dir = File.join(@book_dir, 'markdown')
    if Dir.exist?(markdown_dir)
      (1..20).each do |num|
        filename = sprintf("chapter%02d.md", num)
        filepath = File.join(markdown_dir, filename)
        
        if File.exist?(filepath)
          content = File.read(filepath)
          
          # Extract title
          title_match = content.match(/^#\s+(.+)/)
          title = title_match ? title_match[1] : "Chapter #{num}"
          
          # Process content
          text = content.sub(/^#\s+.+\n/, '') # Remove title
          
          # Convert markdown to HTML
          paragraphs = text.split(/\n\n+/).map do |para|
            para = para.strip
            next if para.empty?
            
            # Convert emphasis
            para = para.gsub(/\*([^*]+)\*/, '<em>\1</em>')
            para = para.gsub(/\*\*([^*]+)\*\*/, '<strong>\1</strong>')
            
            # First paragraph of chapter has no indent
            if para == text.split(/\n\n+/).first
              "<p class=\"first\">#{para}</p>"
            else
              "<p>#{para}</p>"
            end
          end.compact
          
          html += <<~HTML
        <div class="chapter">
          <h1>#{title}</h1>
          #{paragraphs.join("\n      ")}
        </div>
        
          HTML
        end
      end
    end
    
    html += <<~HTML
        <!-- About the Author -->
        <div class="chapter">
          <h1 style="page-break-before: always;">About the Author</h1>
          <p class="first">
            Captain J. Starwind is a pseudonymous author who has spent years exploring
            the cosmic seas of imagination. When not charting new galaxies or negotiating
            with alien civilizations, the Captain enjoys studying the real physics of
            space travel and the sociology of frontier societies.
          </p>
          <p>
            <em>Void Reavers</em> is the first book in The Void Chronicles series,
            with nine more adventures planned in this epic space opera saga.
          </p>
        </div>
        
        <!-- End Page -->
        <div class="chapter">
          <div style="text-align: center; margin-top: 3in;">
            <h2>Thank you for reading!</h2>
            <p style="text-indent: 0;">
              If you enjoyed this book, please consider leaving a review.
            </p>
            <p style="text-indent: 0; margin-top: 2em;">
              <strong>Coming Next in The Void Chronicles:</strong><br>
              <em>Book Two: The Entropy Dancers</em>
            </p>
          </div>
        </div>
      </body>
      </html>
    HTML
    
    File.write(File.join(@book_dir, 'kdp_output', 'void_reavers_kdp.html'), html)
    puts "📄 KDP-formatted HTML saved!"
  end

  def generate_instructions(page_count)
    instructions = <<~INSTRUCTIONS
      Amazon KDP Publishing Instructions
      ==================================
      
      STEP 1: CREATE YOUR PDF
      -----------------------
      Option A: Browser Method (Easiest)
      1. Open void_reavers_kdp.html in Chrome or Firefox
      2. Press Ctrl+P (or Cmd+P on Mac)
      3. Settings:
         - Destination: Save as PDF
         - Paper size: Custom 6×9 inches
         - Margins: Default
         - Scale: 100%
      4. Click Save
      
      Option B: Online Converter
      1. Go to ilovepdf.com/html-pdf
      2. Upload void_reavers_kdp.html
      3. Choose Page Size: Letter
      4. Download PDF
      
      Option C: Professional Tools
      - Adobe Acrobat
      - Microsoft Word (import HTML)
      - LibreOffice Writer (import HTML)
      
      STEP 2: CREATE YOUR COVER
      -------------------------
      Use dimensions from cover_specs.txt
      
      Free Design Options:
      1. Amazon Cover Creator (in KDP dashboard)
      2. Canva.com (search "book cover 6x9")
      3. BookCoverMaker.com
      
      Design Tips:
      - Use high-contrast text
      - Ensure title is readable as thumbnail
      - Leave space for Amazon's barcode
      - Export as PDF at 300 DPI
      
      STEP 3: SET UP KDP ACCOUNT
      --------------------------
      1. Go to kdp.amazon.com
      2. Sign in with Amazon account
      3. Complete tax interview (required)
      4. Add bank account for royalties
      
      STEP 4: CREATE NEW TITLE
      ------------------------
      Book Details:
      - Title: Void Reavers
      - Subtitle: A Tale of Space Pirates and Cosmic Plunder
      - Series: The Void Chronicles (Book 1)
      - Edition: 1
      - Author: [Your name or pen name]
      - Contributors: (optional)
      - Description: [Your book blurb]
      - Publishing Rights: I own the copyright
      - Keywords: space opera, space pirates, science fiction, adventure
      - Categories: Fiction > Science Fiction > Space Opera
                   Fiction > Science Fiction > Adventure
      - Adult Content: No
      
      STEP 5: ISBN DECISION
      ---------------------
      Free KDP ISBN:
      ✓ No cost
      ✗ Only for Amazon distribution
      ✗ Amazon listed as publisher
      
      Your Own ISBN ($125 from Bowker):
      ✓ Wider distribution options
      ✓ You're listed as publisher
      ✓ More professional
      
      STEP 6: PRINT OPTIONS
      ---------------------
      - Print Type: Paperback
      - Interior: Black & white
      - Paper: White or Cream (your choice)
      - Trim Size: 6" × 9"
      - Bleed: No bleed (unless you have images)
      - Cover: Matte or Glossy (your choice)
      
      STEP 7: UPLOAD FILES
      --------------------
      1. Upload interior PDF
      2. Upload cover PDF or use Cover Creator
      3. Preview your book carefully
      4. Fix any issues KDP reports
      
      STEP 8: PRICING
      ---------------
      Printing Cost (#{page_count} pages): ~$#{'%.2f' % (2.15 + (page_count * 0.012))}
      
      Suggested Retail Prices:
      - US: $14.99-16.99
      - UK: £11.99-13.99
      - EU: €13.99-15.99
      
      You earn 60% royalty minus printing costs
      Example at $14.99: ~$#{'%.2f' % (14.99 * 0.6 - (2.15 + page_count * 0.012))} per book
      
      STEP 9: REVIEW & PUBLISH
      ------------------------
      1. Order a proof copy ($#{'%.2f' % (2.15 + page_count * 0.012)} + shipping)
      2. Review physical book carefully
      3. Make any needed corrections
      4. Click "Publish Your Paperback"
      5. Available within 72 hours
      
      MARKETING QUICK START
      ---------------------
      Week Before Launch:
      - Set up author page on Amazon
      - Create Goodreads author account
      - Announce on social media
      
      Launch Week:
      - Price competitively
      - Ask beta readers for reviews
      - Share in relevant groups/forums
      
      Ongoing:
      - Run Amazon ads ($50-100/month)
      - Build email list
      - Write Book 2!
      
      COMMON ISSUES
      -------------
      "Margins too small"
      → Increase margins in HTML/CSS
      
      "Fonts not embedded"
      → Save as PDF/A format
      
      "Low resolution images"
      → Ensure 300 DPI
      
      SUPPORT
      -------
      KDP Help: kdp.amazon.com/help
      Community: facebook.com/groups/20Booksto50k
      
      Good luck with your publishing journey!
    INSTRUCTIONS
    
    File.write(File.join(@book_dir, 'kdp_output', 'kdp_instructions.txt'), instructions)
    puts "📋 Publishing instructions saved!"
  end
end

# Run the generator
if __FILE__ == $0
  generator = SimpleKDPPDFGenerator.new
  generator.generate
end