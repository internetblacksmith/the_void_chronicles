#!/usr/bin/env ruby
# Convert Markdown source to KDP-ready PDF using Pandoc

require 'yaml'
require 'fileutils'
require 'tempfile'

class MarkdownToKDPConverter
  KDP_TRIM_SIZES = {
    'pocket' => { width: '5in', height: '8in', name: 'Pocket Book' },
    'standard' => { width: '6in', height: '9in', name: 'Standard' },
    'large' => { width: '6.14in', height: '9.21in', name: 'Large' }
  }
  
  def initialize(source_dir, trim_size = 'standard')
    @source_dir = source_dir
    @trim_size = KDP_TRIM_SIZES[trim_size] || KDP_TRIM_SIZES['standard']
    @metadata = load_metadata
  end
  
  def generate_pdf(output_file)
    puts "Generating KDP-ready PDF..."
    
    # Create temporary combined markdown file
    combined_md = create_combined_markdown
    
    # Create LaTeX template for KDP
    template = create_kdp_template
    
    # Run Pandoc with proper settings
    pandoc_command = build_pandoc_command(combined_md.path, template.path, output_file)
    
    puts "Running Pandoc..."
    success = system(pandoc_command)
    
    # Cleanup
    combined_md.unlink
    template.unlink
    
    if success
      puts "✓ PDF generated successfully: #{output_file}"
      
      # Check file size
      size_mb = File.size(output_file) / (1024.0 * 1024.0)
      puts "  File size: #{'%.2f' % size_mb} MB"
      
      if size_mb > 650
        puts "  ⚠️  WARNING: File exceeds Amazon's 650 MB limit!"
      end
    else
      puts "✗ PDF generation failed!"
      puts "  Make sure Pandoc is installed: brew install pandoc"
    end
    
    success
  end
  
  private
  
  def load_metadata
    metadata_file = File.join(@source_dir, 'metadata.yaml')
    if File.exist?(metadata_file)
      YAML.load_file(metadata_file)
    else
      {
        'title' => 'Untitled Book',
        'author' => 'Unknown Author',
        'language' => 'en-GB'
      }
    end
  end
  
  def create_combined_markdown
    tempfile = Tempfile.new(['book', '.md'])
    
    # Add front matter
    tempfile.puts "---"
    tempfile.puts "title: \"#{@metadata['title']}\""
    tempfile.puts "author: \"#{@metadata['author']}\""
    tempfile.puts "date: \"#{@metadata['publication_date'] || Date.today.year}\""
    tempfile.puts "lang: #{@metadata['language'] || 'en-GB'}"
    tempfile.puts "---"
    tempfile.puts
    
    # Add chapters
    if @metadata['chapters']
      @metadata['chapters'].each do |chapter|
        chapter_file = File.join(@source_dir, chapter['file'])
        if File.exist?(chapter_file)
          content = File.read(chapter_file, encoding: 'utf-8')
          tempfile.puts content
          tempfile.puts "\n\\newpage\n"
        end
      end
    else
      # Fallback: just concatenate all chapter files
      Dir.glob(File.join(@source_dir, 'chapters', 'chapter-*.md')).sort.each do |chapter_file|
        content = File.read(chapter_file, encoding: 'utf-8')
        tempfile.puts content
        tempfile.puts "\n\\newpage\n"
      end
    end
    
    tempfile.close
    tempfile
  end
  
  def create_kdp_template
    template = Tempfile.new(['kdp-template', '.tex'])
    
    template.puts <<~LATEX
      % KDP-optimized LaTeX template
      \\documentclass[11pt]{book}
      
      % Page geometry for KDP
      \\usepackage[
        paperwidth=#{@trim_size[:width]},
        paperheight=#{@trim_size[:height]},
        top=0.75in,
        bottom=0.75in,
        inner=0.875in,  % Larger inner margin for binding
        outer=0.75in
      ]{geometry}
      
      % Essential packages
      \\usepackage{microtype} % Better typography
      \\usepackage{setspace}
      \\usepackage{hyperref}
      \\usepackage[british]{babel}
      
      % No hyphenation for better e-reader compatibility
      \\usepackage[none]{hyphenat}
      
      % Headers and footers
      \\usepackage{fancyhdr}
      \\pagestyle{fancy}
      \\fancyhf{}
      \\fancyfoot[C]{\\thepage}
      \\renewcommand{\\headrulewidth}{0pt}
      
      % Chapter formatting
      \\usepackage{titlesec}
      \\titleformat{\\chapter}[display]
        {\\normalfont\\Large\\bfseries\\centering}
        {\\chaptertitlename\\ \\thechapter}
        {20pt}
        {\\Huge}
      \\titlespacing*{\\chapter}{0pt}{50pt}{40pt}
      
      % Typography
      \\setstretch{1.5} % 1.5 line spacing for readability
      \\setlength{\\parindent}{1.5em}
      \\setlength{\\parskip}{0pt}
      
      % PDF settings
      \\hypersetup{
        pdfauthor={$author$},
        pdftitle={$title$},
        pdfsubject={$subtitle$},
        pdfkeywords={fiction, science fiction, space opera, pirates},
        colorlinks=false,
        pdfborder={0 0 0}
      }
      
      % Start document
      \\begin{document}
      
      % Title page
      \\begin{titlepage}
      \\centering
      \\vspace*{2in}
      {\\Huge\\bfseries $title$\\par}
      \\vspace{0.5in}
      $if(subtitle)$
      {\\Large $subtitle$\\par}
      \\vspace{0.5in}
      $endif$
      {\\Large\\itshape $author$\\par}
      \\vfill
      \\end{titlepage}
      
      % Copyright page
      \\newpage
      \\thispagestyle{empty}
      \\vspace*{\\fill}
      \\begin{center}
      Copyright \\copyright\\ $date$ by $author$\\\\[0.5cm]
      All rights reserved.\\\\[1cm]
      $if(isbn)$
      ISBN: $isbn$\\\\[0.5cm]
      $endif$
      First Edition\\\\[1cm]
      This is a work of fiction. Names, characters, places, and incidents\\\\
      either are the product of the author's imagination or are used\\\\
      fictitiously, and any resemblance to actual persons, living or dead,\\\\
      business establishments, events, or locales is entirely coincidental.
      \\end{center}
      \\vspace*{\\fill}
      
      % Table of contents
      \\tableofcontents
      
      % Main content
      \\mainmatter
      
      $body$
      
      \\end{document}
    LATEX
    
    template.close
    template
  end
  
  def build_pandoc_command(input_file, template_file, output_file)
    cmd = [
      'pandoc',
      input_file,
      '-o', output_file,
      '--pdf-engine=xelatex',
      '--template', template_file,
      '--top-level-division=chapter',
      '-V documentclass=book',
      '-V fontsize=11pt',
      '-V mainfont="Georgia"',
      '-V sansfont="Arial"',
      '-V monofont="Courier New"',
      '-V linestretch=1.5',
      '-V indent=true',
      '--toc',
      '--toc-depth=1'
    ]
    
    cmd.join(' ')
  end
end

# Command-line interface
puts "📚 Markdown to KDP PDF Converter"
puts "================================"
puts

# Find book directories
book_dirs = Dir.glob('book*_source').select { |d| File.directory?(d) }

if book_dirs.empty?
  puts "❌ No book source directories found!"
  puts "   Expected format: book*_source"
  exit 1
end

# Select book if multiple found, or use the only one
if book_dirs.size == 1
  source_dir = book_dirs.first
  puts "Found book: #{source_dir}"
else
  puts "Multiple books found:"
  book_dirs.each_with_index do |dir, idx|
    puts "  #{idx + 1}. #{dir}"
  end
  print "\nSelect book (1-#{book_dirs.size}): "
  choice = gets.chomp.to_i
  
  if choice < 1 || choice > book_dirs.size
    puts "Invalid selection!"
    exit 1
  end
  
  source_dir = book_dirs[choice - 1]
end

# Get book name for output file
book_name = source_dir.gsub(/_source$/, '')
default_output = "#{book_name}_kdp.pdf"

# Get output filename
if ARGV[0]
  output_file = ARGV[0]
else
  print "\nOutput filename [#{default_output}]: "
  input = gets.chomp
  output_file = input.empty? ? default_output : input
end

# Select trim size
trim_sizes = {
  '1' => 'pocket',
  '2' => 'standard',
  '3' => 'large'
}

puts "\nSelect trim size:"
puts "  1. Pocket (5\" x 8\")"
puts "  2. Standard (6\" x 9\") [default]"
puts "  3. Large (6.14\" x 9.21\")"
print "Choice (1-3) [2]: "
size_choice = gets.chomp
size_choice = '2' if size_choice.empty?

trim_size = trim_sizes[size_choice] || 'standard'

# Check for Pandoc
unless system('which pandoc > /dev/null 2>&1')
  puts "Error: Pandoc is required but not installed"
  puts "Install with:"
  puts "  macOS: brew install pandoc"
  puts "  Ubuntu: sudo apt install pandoc"
  puts "  Arch: sudo pacman -S pandoc"
  exit 1
end

converter = MarkdownToKDPConverter.new(source_dir, trim_size)
success = converter.generate_pdf(output_file)

exit(success ? 0 : 1)