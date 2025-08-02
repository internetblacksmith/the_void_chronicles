#!/usr/bin/env ruby

require 'fileutils'

class LatexToMarkdownConverter
  def initialize(base_dir = Dir.pwd)
    @base_dir = base_dir
    # Check if we're in the main directory or a book directory
    if File.exist?(File.join(@base_dir, 'book.tex'))
      @source_dir = @base_dir
      @output_dir = File.join(@base_dir, 'markdown')
    else
      # Look for book directories
      book_dirs = Dir.glob(File.join(@base_dir, 'book*')).select { |d| File.directory?(d) }
      if book_dirs.any?
        puts "Found book directories: #{book_dirs.map { |d| File.basename(d) }.join(', ')}"
        puts "Which book would you like to convert? (or 'all' for all books)"
        choice = STDIN.gets.chomp
        if choice.downcase == 'all'
          convert_all_books(book_dirs)
          return
        else
          @source_dir = book_dirs.find { |d| File.basename(d).include?(choice) } || book_dirs.first
        end
      else
        puts "No book.tex found and no book directories detected."
        @source_dir = @base_dir
      end
      @output_dir = File.join(@source_dir, 'markdown')
    end
  end

  def convert_all_books(book_dirs)
    book_dirs.each do |book_dir|
      puts "\n" + "=" * 60
      puts "Converting #{File.basename(book_dir)}..."
      puts "=" * 60
      
      converter = LatexToMarkdownConverter.new
      converter.instance_variable_set(:@source_dir, book_dir)
      converter.instance_variable_set(:@output_dir, File.join(book_dir, 'markdown'))
      converter.convert_single_book
    end
  end

  def convert
    return if @source_dir.nil? # Already handled in convert_all_books
    convert_single_book
  end

  def convert_single_book
    puts "LaTeX to Markdown Converter for Void Reavers (Ruby Version)"
    puts "=" * 60

    # Create output directory
    FileUtils.mkdir_p(@output_dir)

    # Get list of chapter files
    chapters = find_chapters

    if chapters.empty?
      puts "No chapter files found in #{@source_dir}!"
      return
    end

    # Convert each chapter
    all_content = []
    
    chapters.each do |chapter_file|
      content = convert_chapter(chapter_file)
      all_content << content if content
    end

    # Create complete book
    create_complete_book(all_content)
    
    # Create README
    create_readme

    puts "\nConversion complete!"
    puts "Markdown files saved in: #{@output_dir}"
  end

  private

  def find_chapters
    # First try to read from book.tex
    book_file = File.join(@source_dir, 'book.tex')
    chapters = []

    if File.exist?(book_file)
      content = File.read(book_file, encoding: 'utf-8')
      # Extract chapter includes
      includes = content.scan(/\\include\{(chapter\d+)\}/).flatten
      chapters = includes.map { |ch| "#{ch}.tex" }
    end

    # Fallback to finding chapter files directly
    if chapters.empty?
      chapters = Dir.glob(File.join(@source_dir, 'chapter*.tex')).map { |f| File.basename(f) }.sort
    end

    chapters
  end

  def convert_chapter(chapter_file)
    tex_path = File.join(@source_dir, chapter_file)
    
    unless File.exist?(tex_path)
      puts "Warning: #{chapter_file} not found!"
      return nil
    end

    puts "Converting #{chapter_file}..."

    # Read LaTeX content
    latex_content = File.read(tex_path, encoding: 'utf-8')
    
    # Convert to Markdown
    md_content = convert_latex_to_markdown(latex_content)
    
    # Save individual chapter
    md_filename = chapter_file.sub('.tex', '.md')
    md_path = File.join(@output_dir, md_filename)
    
    File.write(md_path, md_content, encoding: 'utf-8')
    puts "  → Saved to #{md_filename}"

    md_content
  end

  def convert_latex_to_markdown(content)
    # Convert chapter headings
    content = content.gsub(/\\chapter\{([^}]+)\}/, '# \1')
    
    # Convert double quotes
    content = content.gsub(/``([^']+)''/, '"\1"')
    
    # Convert italics
    content = content.gsub(/\\textit\{([^}]+)\}/, '*\1*')
    
    # Convert bold
    content = content.gsub(/\\textbf\{([^}]+)\}/, '**\1**')
    
    # Convert single quotes with backticks
    content = content.gsub(/`([^']+)'/, "'\1'")
    
    # Convert em-dashes
    content = content.gsub('---', '—')
    
    # Remove LaTeX escapes for special characters
    replacements = {
      '\%' => '%',
      '\$' => '$',
      '\&' => '&',
      '\_' => '_',
      '\#' => '#'
    }
    
    replacements.each do |latex, plain|
      content = content.gsub(latex, plain)
    end
    
    # Clean up excessive newlines
    content = content.gsub(/\n{3,}/, "\n\n")
    
    content
  end

  def create_complete_book(chapters)
    return if chapters.empty?

    puts "\nCreating complete book..."

    book_content = []
    
    # Add header
    book_content << "# Void Reavers\n"
    book_content << "## A Tale of Space Pirates and Cosmic Plunder\n"
    book_content << "By Captain J. Starwind\n"
    book_content << "---\n"
    
    # Add all chapters with separators
    chapters.each_with_index do |chapter, index|
      book_content << chapter
      book_content << "\n---\n" unless index == chapters.length - 1
    end
    
    # Write complete book
    book_path = File.join(@output_dir, 'void_reavers_complete.md')
    File.write(book_path, book_content.join("\n"), encoding: 'utf-8')
    
    puts "Complete book saved to: void_reavers_complete.md"
  end

  def create_readme
    readme_content = <<~README
      # Void Reavers - Markdown Version

      This directory contains the markdown version of the book "Void Reavers".

      ## Files

      - `void_reavers_complete.md` - The complete book in a single file
      - `chapter01.md` through `chapter20.md` - Individual chapter files

      ## About

      This book tells the story of Captain Zara "Bloodhawk" Vega and the transformation of space piracy in a universe where humanity must prove itself worthy of the stars.

      ## Synopsis

      From the lawless void between stars to the halls of galactic diplomacy, follow Captain Zara Vega's fifty-year journey as she transforms from a young pirate forced into Rex Morrison's brutal crew to humanity's ambassador to alien civilizations. 

      In a universe where quantum physics can tear reality apart and ancient alien watchers judge humanity's every move, pirates must evolve from raiders to protectors, proving that even thieves can have honor when the survival of the species is at stake.

      ## Reading Order

      The chapters should be read in numerical order from 01 to 20. The complete book file contains all chapters in the correct order.

      ## Conversion Details

      This markdown version was converted from the original LaTeX source using a Ruby script that preserves formatting while making the text more accessible for general readers.
    README

    readme_path = File.join(@output_dir, 'README.md')
    File.write(readme_path, readme_content, encoding: 'utf-8')
    
    puts "README created"
  end
end

# Run the converter if this script is executed directly
if __FILE__ == $0
  converter = LatexToMarkdownConverter.new
  converter.convert
end