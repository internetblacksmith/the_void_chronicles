#!/usr/bin/env ruby
# Publish Markdown source to PDF and EPUB without build metadata

require 'yaml'
require 'fileutils'
require 'tempfile'

class MarkdownToPublishConverter
  def initialize(source_dir)
    @source_dir = source_dir
    @metadata = load_metadata
  end

  def publish(output_file, format)
    puts "Publishing #{format.upcase} from Markdown source..."

    # Create temporary combined markdown file
    combined_md = create_combined_markdown

    # Build Pandoc command
    cmd = build_pandoc_command(combined_md.path, output_file, format)

    puts "Running Pandoc..."
    success = system(cmd)

    # Cleanup
    combined_md.unlink

    if success
      puts "\u2713 #{format.upcase} published successfully: #{output_file}"
      size_mb = File.size(output_file) / (1024.0 * 1024.0)
      puts "  File size: #{'%.2f' % size_mb} MB"
    else
      puts "\u2717 #{format.upcase} publishing failed!"
      puts "  Make sure Pandoc is installed"
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

    # Add metadata header
    tempfile.puts "---"
    tempfile.puts "title: \"#{@metadata['title']}\""
    tempfile.puts "author: \"#{@metadata['author']}\""
    tempfile.puts "date: \"#{@metadata['publication_date'] || Date.today.year}\""
    tempfile.puts "lang: #{@metadata['language'] || 'en-GB'}"
    tempfile.puts "---"
    tempfile.puts

    # Add all chapters
    if @metadata['chapters']
      @metadata['chapters'].each do |chapter|
        chapter_file = File.join(@source_dir, chapter['file'])
        if File.exist?(chapter_file)
          content = File.read(chapter_file, encoding: 'utf-8')
          tempfile.puts content
          tempfile.puts "\n"
        end
      end
    else
      Dir.glob(File.join(@source_dir, 'chapters', 'chapter-*.md')).sort.each do |chapter_file|
        content = File.read(chapter_file, encoding: 'utf-8')
        tempfile.puts content
        tempfile.puts "\n"
      end
    end

    tempfile.close
    tempfile
  end

  def build_pandoc_command(input_file, output_file, format)
    cmd = [
      'pandoc',
      input_file,
      '-o', output_file,
      '--toc',
      '--toc-depth=1'
    ]

    if format == 'epub'
      cmd << '--epub-chapter-level=1'
    elsif format == 'pdf'
      cmd << '--pdf-engine=xelatex'
      cmd << '-V documentclass=book'
      cmd << '-V fontsize=11pt'
      cmd << '-V mainfont="Georgia"'
      cmd << '-V sansfont="Arial"'
      cmd << '-V monofont="Courier New"'
      cmd << '-V linestretch=1.5'
      cmd << '-V indent=true'
    end

    cmd.join(' ')
  end
end

# Main script
puts "\u2728 Markdown to Publish Converter"
puts "============================"
puts

# Find book directories
book_dirs = Dir.glob('book*_source').select { |d| File.directory?(d) }

if book_dirs.empty?
  puts "\u274c No book source directories found!"
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
default_output_pdf = "#{book_name}.pdf"
default_output_epub = "#{book_name}.epub"

# Get output filename
print "\nOutput PDF filename [#{default_output_pdf}]: "
pdf_input = gets.chomp
output_pdf = pdf_input.empty? ? default_output_pdf : pdf_input

print "\nOutput EPUB filename [#{default_output_epub}]: "
epub_input = gets.chomp
output_epub = epub_input.empty? ? default_output_epub : epub_input

unless system('which pandoc > /dev/null 2>&1')
  puts "Error: Pandoc is required but not installed"
  exit 1
end

converter = MarkdownToPublishConverter.new(source_dir)
converter.publish(output_pdf, 'pdf')
converter.publish(output_epub, 'epub')