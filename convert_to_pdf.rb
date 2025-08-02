#!/usr/bin/env ruby

require 'fileutils'
require 'open3'

class MarkdownToPdfConverter
  def initialize(base_dir = Dir.pwd)
    @base_dir = base_dir
    @markdown_dir = File.join(@base_dir, 'markdown')
    @pdf_dir = File.join(@base_dir, 'pdf')
  end

  def convert
    puts "Markdown to PDF Converter for Void Reavers"
    puts "=" * 60

    # Check if markdown directory exists
    unless Dir.exist?(@markdown_dir)
      puts "Error: markdown directory not found!"
      puts "Please run convert_to_md.rb first to generate markdown files."
      return
    end

    # Check for required tools
    unless check_dependencies
      return
    end

    # Create output directory
    FileUtils.mkdir_p(@pdf_dir)

    # Convert options
    puts "\nConversion options:"
    puts "1. Convert complete book to single PDF"
    puts "2. Convert individual chapters to separate PDFs"
    puts "3. Convert both"
    print "Choose option (1-3): "
    
    option = STDIN.gets.chomp
    
    case option
    when '1'
      convert_complete_book
    when '2'
      convert_individual_chapters
    when '3'
      convert_complete_book
      convert_individual_chapters
    else
      puts "Invalid option. Defaulting to complete book only."
      convert_complete_book
    end

    puts "\nConversion complete!"
    puts "PDF files saved in: #{@pdf_dir}"
  end

  private

  def check_dependencies
    # Check for pandoc
    pandoc_check = system('which pandoc > /dev/null 2>&1')
    unless pandoc_check
      puts "\nError: pandoc is not installed!"
      puts "Please install pandoc first:"
      puts "  Ubuntu/Debian: sudo apt-get install pandoc texlive-latex-base texlive-fonts-recommended"
      puts "  MacOS: brew install pandoc basictex"
      puts "  Or visit: https://pandoc.org/installing.html"
      puts "\nAlternatively, you can use markdown_to_html.rb to create an HTML version"
      puts "that can be printed to PDF from your web browser."
      return false
    end

    # Check for LaTeX (for better PDF output)
    latex_check = system('which pdflatex > /dev/null 2>&1')
    unless latex_check
      puts "\nWarning: LaTeX is not installed. PDF output quality may be reduced."
      puts "For better results, install LaTeX (texlive or basictex)."
      print "Continue anyway? (y/n): "
      response = STDIN.gets.chomp.downcase
      return response == 'y'
    end

    true
  end

  def convert_complete_book
    puts "\nConverting complete book to PDF..."
    
    md_file = File.join(@markdown_dir, 'void_reavers_complete.md')
    unless File.exist?(md_file)
      puts "Error: void_reavers_complete.md not found!"
      return
    end

    pdf_file = File.join(@pdf_dir, 'void_reavers_complete.pdf')
    
    # Pandoc options for book-like output
    pandoc_cmd = [
      'pandoc',
      md_file,
      '-o', pdf_file,
      '--pdf-engine=pdflatex',
      '-V', 'documentclass=book',
      '-V', 'geometry:margin=1in',
      '-V', 'fontsize=11pt',
      '-V', 'linkcolor=blue',
      '--toc',
      '--toc-depth=1',
      '-V', 'title=Void Reavers',
      '-V', 'author=Captain J. Starwind',
      '-V', 'subtitle=A Tale of Space Pirates and Cosmic Plunder',
      '--metadata', 'title=Void Reavers',
      '--standalone'
    ]

    # Try with LaTeX first, fall back to basic if it fails
    success = run_pandoc_command(pandoc_cmd)
    
    if !success && system('which pdflatex > /dev/null 2>&1')
      puts "LaTeX conversion failed, trying simpler approach..."
      pandoc_cmd = [
        'pandoc',
        md_file,
        '-o', pdf_file,
        '--pdf-engine=wkhtmltopdf',
        '--toc',
        '--metadata', 'title=Void Reavers'
      ]
      success = run_pandoc_command(pandoc_cmd)
    end

    if success
      puts "  ✓ Complete book converted to PDF"
      puts "  → #{pdf_file}"
    else
      puts "  ✗ Failed to convert complete book"
    end
  end

  def convert_individual_chapters
    puts "\nConverting individual chapters to PDF..."
    
    chapter_files = Dir.glob(File.join(@markdown_dir, 'chapter*.md')).sort
    
    if chapter_files.empty?
      puts "No chapter files found!"
      return
    end

    chapter_files.each do |md_file|
      basename = File.basename(md_file, '.md')
      pdf_file = File.join(@pdf_dir, "#{basename}.pdf")
      
      # Extract chapter title from the file
      chapter_title = extract_chapter_title(md_file)
      
      pandoc_cmd = [
        'pandoc',
        md_file,
        '-o', pdf_file,
        '--pdf-engine=pdflatex',
        '-V', 'geometry:margin=1in',
        '-V', 'fontsize=11pt',
        '-V', 'linkcolor=blue',
        '-V', "title=#{chapter_title}",
        '-V', 'author=Captain J. Starwind',
        '--standalone'
      ]

      success = run_pandoc_command(pandoc_cmd, quiet: true)
      
      if success
        puts "  ✓ #{basename}.md → #{basename}.pdf"
      else
        puts "  ✗ Failed to convert #{basename}.md"
      end
    end
  end

  def run_pandoc_command(cmd, quiet: false)
    if quiet
      system(*cmd, out: File::NULL, err: File::NULL)
    else
      system(*cmd)
    end
  end

  def extract_chapter_title(md_file)
    # Try to extract the first # heading as the chapter title
    File.foreach(md_file) do |line|
      if line.start_with?('# ')
        return line.sub('# ', '').strip
      end
    end
    'Chapter'
  end
end

# Alternative simple conversion function using markdown-pdf if available
def simple_markdown_to_pdf
  puts "\nAttempting simple conversion using markdown-pdf..."
  
  # Check if markdown-pdf is installed
  unless system('which markdown-pdf > /dev/null 2>&1')
    puts "markdown-pdf is not installed."
    puts "Install with: npm install -g markdown-pdf"
    return false
  end

  FileUtils.mkdir_p('pdf')
  
  # Convert complete book
  system('markdown-pdf markdown/void_reavers_complete.md -o pdf/void_reavers_complete.pdf')
  
  puts "Simple conversion complete!"
  true
end

# Create a wrapper script for easy PDF generation
def create_pdf_wrapper_script
  wrapper_content = <<~SCRIPT
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
        pandoc markdown/void_reavers_complete.md \\
            -o pdf/void_reavers_complete.pdf \\
            --toc \\
            --standalone \\
            -V geometry:margin=1in \\
            -V fontsize=11pt \\
            -V title="Void Reavers" \\
            -V author="Captain J. Starwind" \\
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
  SCRIPT

  File.write('convert_to_pdf.sh', wrapper_content)
  File.chmod(0755, 'convert_to_pdf.sh')
  puts "Created convert_to_pdf.sh wrapper script"
end

# Run the converter if this script is executed directly
if __FILE__ == $0
  # Create the wrapper script
  create_pdf_wrapper_script
  
  # Run the Ruby converter
  converter = MarkdownToPdfConverter.new
  converter.convert
end