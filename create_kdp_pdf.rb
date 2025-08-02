#!/usr/bin/env ruby

require 'fileutils'
require 'erb'

class KDPPDFGenerator
  # Amazon KDP standard book sizes
  TRIM_SIZES = {
    '5x8' => { width: '5in', height: '8in', name: '5" × 8" (Small Novel)' },
    '5.25x8' => { width: '5.25in', height: '8in', name: '5.25" × 8" (Standard Novel)' },
    '5.5x8.5' => { width: '5.5in', height: '8.5in', name: '5.5" × 8.5" (Trade)' },
    '6x9' => { width: '6in', height: '9in', name: '6" × 9" (US Trade - Most Common)' },
    '6.14x9.21' => { width: '6.14in', height: '9.21in', name: '6.14" × 9.21" (Royal)' },
    '7x10' => { width: '7in', height: '10in', name: '7" × 10" (Textbook)' },
    '8.5x11' => { width: '8.5in', height: '11in', name: '8.5" × 11" (Large)' }
  }.freeze

  def initialize(base_dir = Dir.pwd)
    @base_dir = base_dir
    @book_dir = find_book_directory
  end

  def generate
    unless @book_dir
      puts "❌ No book directory found!"
      return
    end

    puts "\n📚 Amazon KDP PDF Generator"
    puts "==========================="
    puts "\nThis will create a print-ready PDF for Amazon KDP publishing."
    
    # Select trim size
    trim_size = select_trim_size
    
    # Get page count estimate
    page_count = estimate_page_count
    
    # Ask about bleed
    use_bleed = ask_about_bleed
    
    puts "\n🔧 Generating KDP-compliant PDF..."
    
    # Create output directory
    output_dir = File.join(@book_dir, 'kdp_output')
    FileUtils.mkdir_p(output_dir)
    
    # Generate LaTeX file with KDP settings
    latex_file = generate_kdp_latex(trim_size, use_bleed, page_count)
    
    # Compile to PDF
    compile_to_pdf(latex_file, output_dir)
    
    # Generate cover template info
    generate_cover_info(trim_size, page_count)
    
    puts "\n✅ KDP PDF generation complete!"
    puts "\n📋 Next Steps:"
    puts "1. Review the PDF in #{output_dir}"
    puts "2. Create your cover using the specifications in cover_specs.txt"
    puts "3. Upload both files to Amazon KDP"
    puts "4. Order a proof copy before publishing"
  end

  private

  def find_book_directory
    Dir.glob(File.join(@base_dir, 'book*')).select { |d| File.directory?(d) }.first
  end

  def select_trim_size
    puts "\n📏 Select your book trim size:"
    puts "(This is the final size of your printed book)"
    puts ""
    
    TRIM_SIZES.each_with_index do |(key, info), index|
      puts "#{index + 1}. #{info[:name]}"
    end
    
    print "\nSelect size (1-#{TRIM_SIZES.length}) [default: 4 for 6×9]: "
    choice = gets.chomp
    choice = '4' if choice.empty?
    
    index = choice.to_i - 1
    if index >= 0 && index < TRIM_SIZES.length
      TRIM_SIZES.keys[index]
    else
      '6x9' # Default to most common
    end
  end

  def estimate_page_count
    # Count words in all chapters
    word_count = 0
    Dir.glob(File.join(@book_dir, 'chapter*.tex')).each do |file|
      content = File.read(file)
      # Remove LaTeX commands for word count
      text = content.gsub(/\\[a-zA-Z]+\{[^}]*\}/, ' ')
      text = text.gsub(/\\[a-zA-Z]+/, ' ')
      word_count += text.split.length
    end
    
    # Estimate pages (average 250 words per page for 6x9)
    estimated_pages = (word_count / 250.0).ceil
    
    puts "\n📊 Book Statistics:"
    puts "   Word count: ~#{word_count}"
    puts "   Estimated pages: ~#{estimated_pages}"
    
    print "\nEnter actual/expected page count [#{estimated_pages}]: "
    input = gets.chomp
    input.empty? ? estimated_pages : input.to_i
  end

  def ask_about_bleed
    puts "\n🎨 Do you have images that should extend to the page edges?"
    print "Use bleed settings? (y/N): "
    gets.chomp.downcase == 'y'
  end

  def generate_kdp_latex(trim_size, use_bleed, page_count)
    size_info = TRIM_SIZES[trim_size]
    
    # Calculate dimensions with bleed if needed
    if use_bleed
      # Add 0.125" bleed on three sides (0.25" total height, 0.125" width)
      width = size_info[:width].sub(/(\d+\.?\d*)in/, "\\1")
      height = size_info[:height].sub(/(\d+\.?\d*)in/, "\\1")
      
      width_with_bleed = (width.to_f + 0.125).to_s + "in"
      height_with_bleed = (height.to_f + 0.25).to_s + "in"
      
      paperwidth = width_with_bleed
      paperheight = height_with_bleed
    else
      paperwidth = size_info[:width]
      paperheight = size_info[:height]
    end
    
    latex_content = <<~LATEX
      % Amazon KDP-compliant book format
      % Generated for #{size_info[:name]} trim size
      % #{use_bleed ? 'With' : 'Without'} bleed
      
      \\documentclass[12pt,oneside]{book}
      
      % Page setup for KDP
      \\usepackage[
        paperwidth=#{paperwidth},
        paperheight=#{paperheight},
        % Inner margin (binding side) - slightly larger
        inner=0.75in,
        % Outer margin
        outer=0.5in,
        % Top and bottom margins
        top=0.75in,
        bottom=0.75in,
        % No headers/footers in margins
        includehead=false,
        includefoot=false
      ]{geometry}
      
      % Font settings
      \\usepackage{times} % or {bookman} or {palatino}
      \\usepackage[T1]{fontenc}
      \\usepackage[utf8]{inputenc}
      
      % Better typography
      \\usepackage{microtype}
      
      % Paragraph settings
      \\setlength{\\parindent}{0.25in}
      \\setlength{\\parskip}{0pt}
      
      % Line spacing (1.2 is good for novels)
      \\linespread{1.2}
      
      % Chapter title formatting
      \\usepackage{titlesec}
      \\titleformat{\\chapter}[display]
        {\\normalfont\\Large\\bfseries\\centering}
        {\\chaptertitlename\\ \\thechapter}
        {20pt}
        {\\Huge}
      \\titlespacing*{\\chapter}{0pt}{50pt}{40pt}
      
      % Headers and footers
      \\usepackage{fancyhdr}
      \\pagestyle{fancy}
      \\fancyhf{} % Clear all headers and footers
      \\fancyfoot[C]{\\thepage} % Page number at bottom center
      \\renewcommand{\\headrulewidth}{0pt} % No header line
      
      % Ensure 300 DPI for images
      \\pdfcompresslevel=9
      \\pdfimageresolution=300
      
      % PDF metadata
      \\usepackage[
        pdftitle={Void Reavers: A Tale of Space Pirates and Cosmic Plunder},
        pdfauthor={Captain J. Starwind},
        pdfsubject={Science Fiction Novel},
        pdfkeywords={space pirates, science fiction, space opera},
        pdfproducer={LaTeX},
        pdfcreator={pdflatex}
      ]{hyperref}
      
      \\begin{document}
      
      % Front matter
      \\frontmatter
      
      % Title page
      \\begin{titlepage}
      \\vspace*{\\fill}
      \\begin{center}
      {\\Huge\\bfseries Void Reavers}\\\\[0.5cm]
      {\\Large A Tale of Space Pirates and Cosmic Plunder}\\\\[2cm]
      {\\Large\\itshape Captain J. Starwind}\\\\[1cm]
      \\vfill
      {\\large Book One of The Void Chronicles}
      \\end{center}
      \\vspace*{\\fill}
      \\end{titlepage}
      
      % Copyright page
      \\clearpage
      \\thispagestyle{empty}
      \\vspace*{\\fill}
      \\begin{center}
      Copyright \\copyright\\ 2024 by Captain J. Starwind\\\\[0.5cm]
      All rights reserved.\\\\[1cm]
      ISBN: [Your ISBN here]\\\\[0.5cm]
      First Edition\\\\[1cm]
      This is a work of fiction. Names, characters, places, and incidents\\\\
      either are the product of the author's imagination or are used\\\\
      fictitiously, and any resemblance to actual persons, living or dead,\\\\
      business establishments, events, or locales is entirely coincidental.\\\\[1cm]
      Published in the United States of America
      \\end{center}
      \\vspace*{\\fill}
      
      % Table of contents
      \\tableofcontents
      
      % Main matter
      \\mainmatter
      
      % Include all chapters
      \\input{chapter01}
      \\input{chapter02}
      \\input{chapter03}
      \\input{chapter04}
      \\input{chapter05}
      \\input{chapter06}
      \\input{chapter07}
      \\input{chapter08}
      \\input{chapter09}
      \\input{chapter10}
      \\input{chapter11}
      \\input{chapter12}
      \\input{chapter13}
      \\input{chapter14}
      \\input{chapter15}
      \\input{chapter16}
      \\input{chapter17}
      \\input{chapter18}
      \\input{chapter19}
      \\input{chapter20}
      
      % Back matter
      \\backmatter
      
      % About the Author
      \\chapter*{About the Author}
      \\addcontentsline{toc}{chapter}{About the Author}
      
      Captain J. Starwind is a pseudonymous author who has spent years exploring
      the cosmic seas of imagination. When not charting new galaxies or negotiating
      with alien civilizations, the Captain enjoys studying the real physics of
      space travel and the sociology of frontier societies.
      
      \\textit{Void Reavers} is the first book in The Void Chronicles series,
      with nine more adventures planned in this epic space opera saga.
      
      % End of book page
      \\clearpage
      \\thispagestyle{empty}
      \\vspace*{\\fill}
      \\begin{center}
      {\\Large Thank you for reading!}\\\\[1cm]
      If you enjoyed this book, please consider leaving a review.\\\\[1cm]
      {\\large Coming Next in The Void Chronicles:}\\\\[0.5cm]
      {\\itshape Book Two: The Entropy Dancers}
      \\end{center}
      \\vspace*{\\fill}
      
      \\end{document}
    LATEX
    
    # Save LaTeX file
    latex_file = File.join(@book_dir, 'void_reavers_kdp.tex')
    File.write(latex_file, latex_content)
    
    latex_file
  end

  def compile_to_pdf(latex_file, output_dir)
    puts "\n🔨 Compiling PDF (this may take a moment)..."
    
    Dir.chdir(File.dirname(latex_file)) do
      # First pass
      system("pdflatex -interaction=nonstopmode -output-directory=#{output_dir} #{latex_file} > /dev/null 2>&1")
      
      # Second pass for TOC and references
      system("pdflatex -interaction=nonstopmode -output-directory=#{output_dir} #{latex_file} > /dev/null 2>&1")
      
      # Third pass to ensure everything is correct
      system("pdflatex -interaction=nonstopmode -output-directory=#{output_dir} #{latex_file} > /dev/null 2>&1")
    end
    
    pdf_file = File.join(output_dir, 'void_reavers_kdp.pdf')
    if File.exist?(pdf_file)
      puts "✅ PDF created: #{pdf_file}"
      
      # Check file size
      size_mb = File.size(pdf_file) / 1024.0 / 1024.0
      puts "📏 File size: #{'%.2f' % size_mb} MB (max allowed: 650 MB)"
      
      if size_mb > 650
        puts "⚠️  WARNING: File size exceeds Amazon's 650 MB limit!"
        puts "   Consider reducing image quality or removing some images."
      end
    else
      puts "❌ PDF compilation failed. Please check if pdflatex is installed."
      puts "   Ubuntu/Debian: sudo apt install texlive-full"
      puts "   macOS: brew install --cask mactex"
    end
  end

  def generate_cover_info(trim_size, page_count)
    size_info = TRIM_SIZES[trim_size]
    
    # Calculate spine width (KDP formula)
    # White paper: page count * 0.0025"
    # Cream paper: page count * 0.002739"
    spine_width_white = page_count * 0.0025
    spine_width_cream = page_count * 0.002739
    
    # Parse dimensions
    width = size_info[:width].sub(/(\d+\.?\d*)in/, "\\1").to_f
    height = size_info[:height].sub(/(\d+\.?\d*)in/, "\\1").to_f
    
    # Calculate full cover dimensions (front + spine + back + bleed)
    # Bleed: 0.125" on all sides
    bleed = 0.125
    
    cover_width_white = (width * 2) + spine_width_white + (bleed * 2)
    cover_height = height + (bleed * 2)
    
    cover_width_cream = (width * 2) + spine_width_cream + (bleed * 2)
    
    cover_specs = <<~SPECS
      Amazon KDP Cover Specifications
      ================================
      
      Book Details:
      - Trim Size: #{size_info[:name]}
      - Page Count: #{page_count}
      
      Cover Dimensions (White Paper):
      - Width: #{'%.3f' % cover_width_white}" (#{'%.1f' % (cover_width_white * 25.4)} mm)
      - Height: #{'%.3f' % cover_height}" (#{'%.1f' % (cover_height * 25.4)} mm)
      - Spine Width: #{'%.3f' % spine_width_white}" (#{'%.1f' % (spine_width_white * 25.4)} mm)
      
      Cover Dimensions (Cream Paper):
      - Width: #{'%.3f' % cover_width_cream}" (#{'%.1f' % (cover_width_cream * 25.4)} mm)
      - Height: #{'%.3f' % cover_height}" (#{'%.1f' % (cover_height * 25.4)} mm)
      - Spine Width: #{'%.3f' % spine_width_cream}" (#{'%.1f' % (spine_width_cream * 25.4)} mm)
      
      Requirements:
      - Resolution: 300 DPI minimum
      - Color mode: CMYK for best print results (RGB accepted)
      - Bleed: 0.125" (3.2 mm) on all sides
      - Safe zone: Keep important content 0.25" (6.4 mm) from edges
      - Spine text: Only for books with 79+ pages
      - Spine text margin: 0.0625" (1.6 mm) from spine edges
      - File format: PDF
      - File size: Maximum 650 MB
      
      Cover Layout (left to right):
      1. Back Cover: #{width}" × #{height}"
      2. Spine: #{'%.3f' % spine_width_white}" × #{height}" (white paper)
      3. Front Cover: #{width}" × #{height}"
      
      Plus 0.125" bleed on all outer edges.
      
      Important Notes:
      - Download KDP's cover template for exact measurements
      - Ensure all text is at least 0.125" from spine edges
      - Barcode space (if needed): 2" × 1.2" on back cover
      - ISBN barcode will be added automatically by Amazon
      
      For cover template generator, visit:
      https://kdp.amazon.com/en_US/cover-templates
    SPECS
    
    specs_file = File.join(@book_dir, 'kdp_output', 'cover_specs.txt')
    File.write(specs_file, cover_specs)
    
    puts "\n📐 Cover specifications saved to: cover_specs.txt"
    puts "   Spine width (white paper): #{'%.3f' % spine_width_white}\" (#{page_count} pages)"
    puts "   Total cover width: #{'%.3f' % cover_width_white}\" × #{'%.3f' % cover_height}\""
  end
end

# Check for pdflatex
unless system('which pdflatex > /dev/null 2>&1')
  puts "❌ pdflatex not found. Please install LaTeX:"
  puts "   Ubuntu/Debian: sudo apt install texlive-full"
  puts "   macOS: brew install --cask mactex"
  puts "   Windows: Install MiKTeX from miktex.org"
  exit 1
end

# Run the generator
if __FILE__ == $0
  generator = KDPPDFGenerator.new
  generator.generate
end