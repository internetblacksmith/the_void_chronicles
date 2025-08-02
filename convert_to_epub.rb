#!/usr/bin/env ruby

require 'fileutils'
require 'erb'
require 'cgi'

class LaTeXToEPUBConverter
  def initialize(base_dir = Dir.pwd)
    @base_dir = base_dir
    @book_dirs = find_book_directories
  end

  def convert
    if @book_dirs.empty?
      puts "No book directories found!"
      return
    end

    book_dir = select_book_directory
    return unless book_dir

    puts "\n📚 Converting #{File.basename(book_dir)} to EPUB..."
    
    # Create temporary working directory
    temp_dir = File.join(book_dir, 'epub_temp')
    FileUtils.mkdir_p(temp_dir)
    
    begin
      # Load book metadata and chapters
      book_info = load_book_info(book_dir)
      chapters = load_chapters(book_dir)
      
      if chapters.empty?
        puts "❌ No chapters found!"
        return
      end
      
      # Create EPUB structure
      create_epub_structure(temp_dir)
      
      # Generate EPUB files
      generate_mimetype(temp_dir)
      generate_container_xml(temp_dir)
      generate_content_opf(temp_dir, book_info, chapters)
      generate_toc_ncx(temp_dir, book_info, chapters)
      generate_toc_xhtml(temp_dir, book_info, chapters)
      generate_cover_xhtml(temp_dir, book_info)
      generate_chapters_xhtml(temp_dir, chapters)
      generate_css(temp_dir)
      
      # Create EPUB file
      epub_file = create_epub_file(book_dir, temp_dir, book_info)
      
      puts "\n✅ EPUB created successfully: #{epub_file}"
      puts "\n📱 Transfer this file to your Kindle via:"
      puts "   - Email to your Kindle address"
      puts "   - USB cable transfer"
      puts "   - Kindle app's Send to Kindle feature"
      
    ensure
      # Clean up temporary directory
      FileUtils.rm_rf(temp_dir)
    end
  end

  private

  def find_book_directories
    Dir.glob(File.join(@base_dir, 'book*')).select { |d| File.directory?(d) }
  end

  def select_book_directory
    if @book_dirs.length == 1
      return @book_dirs.first
    end

    puts "\nFound multiple book directories:"
    @book_dirs.each_with_index do |dir, index|
      puts "#{index + 1}. #{File.basename(dir)}"
    end
    puts "#{@book_dirs.length + 1}. All books"

    print "\nSelect a book to convert (1-#{@book_dirs.length + 1}): "
    choice = gets.chomp.to_i

    if choice == @book_dirs.length + 1
      # Convert all books
      @book_dirs.each { |dir| convert_single_book(dir) }
      nil
    elsif choice > 0 && choice <= @book_dirs.length
      @book_dirs[choice - 1]
    else
      puts "Invalid selection!"
      nil
    end
  end

  def load_book_info(book_dir)
    {
      title: "Void Reavers",
      subtitle: "A Tale of Space Pirates and Cosmic Plunder",
      author: "Captain J. Starwind",
      language: "en",
      publisher: "Void Chronicles Publishing",
      year: Time.now.year.to_s,
      description: "Follow Captain Zara 'Bloodhawk' Vega's fifty-year journey from pirate to diplomat in a universe where humanity must prove itself worthy of cosmic citizenship.",
      isbn: generate_uuid
    }
  end

  def load_chapters(book_dir)
    chapters = []
    
    # Try to load chapters from LaTeX files
    (1..20).each do |num|
      filename = sprintf("chapter%02d.tex", num)
      filepath = File.join(book_dir, filename)
      
      if File.exist?(filepath)
        content = File.read(filepath, encoding: 'UTF-8')
        title, text = parse_latex_chapter(content)
        
        chapters << {
          number: num,
          title: title || "Chapter #{num}",
          content: text,
          filename: sprintf("chapter%02d.xhtml", num)
        }
      end
    end
    
    chapters
  end

  def parse_latex_chapter(content)
    # Extract chapter title
    title_match = content.match(/\\chapter\{([^}]+)\}/)
    title = title_match ? title_match[1] : nil
    
    # Convert LaTeX to HTML-friendly format
    text = content.dup
    
    # Remove chapter command
    text.gsub!(/\\chapter\{[^}]+\}/, '')
    
    # Convert LaTeX formatting to HTML
    text.gsub!(/\\textit\{([^}]+)\}/, '<em>\1</em>')
    text.gsub!(/\\textbf\{([^}]+)\}/, '<strong>\1</strong>')
    text.gsub!(/``([^']+)''/, '"\1"')
    text.gsub!(/`([^']+)'/, "'\1'")
    text.gsub!(/---/, '—')
    
    # Handle special characters
    text.gsub!(/\\%/, '%')
    text.gsub!(/\\\$/, '$')
    text.gsub!(/\\&/, '&')
    text.gsub!(/\\_/, '_')
    text.gsub!(/\\#/, '#')
    
    # Convert paragraphs
    paragraphs = text.split(/\n\n+/).map(&:strip).reject(&:empty?)
    formatted_text = paragraphs.map { |p| "<p>#{CGI.escapeHTML(p).gsub(/&lt;(\/?(?:em|strong))&gt;/, '<\1>')}</p>" }.join("\n\n")
    
    [title, formatted_text]
  end

  def create_epub_structure(temp_dir)
    FileUtils.mkdir_p(File.join(temp_dir, 'META-INF'))
    FileUtils.mkdir_p(File.join(temp_dir, 'OEBPS'))
    FileUtils.mkdir_p(File.join(temp_dir, 'OEBPS', 'css'))
    FileUtils.mkdir_p(File.join(temp_dir, 'OEBPS', 'images'))
  end

  def generate_mimetype(temp_dir)
    File.write(File.join(temp_dir, 'mimetype'), 'application/epub+zip', mode: 'wb')
  end

  def generate_container_xml(temp_dir)
    container = <<~XML
      <?xml version="1.0" encoding="UTF-8"?>
      <container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
        <rootfiles>
          <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
        </rootfiles>
      </container>
    XML
    
    File.write(File.join(temp_dir, 'META-INF', 'container.xml'), container)
  end

  def generate_content_opf(temp_dir, book_info, chapters)
    content_opf = <<~XML
      <?xml version="1.0" encoding="UTF-8"?>
      <package xmlns="http://www.idpf.org/2007/opf" unique-identifier="uid" version="3.0">
        <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
          <dc:identifier id="uid">#{book_info[:isbn]}</dc:identifier>
          <dc:title>#{CGI.escapeHTML(book_info[:title])}</dc:title>
          <dc:creator>#{CGI.escapeHTML(book_info[:author])}</dc:creator>
          <dc:language>#{book_info[:language]}</dc:language>
          <dc:publisher>#{CGI.escapeHTML(book_info[:publisher])}</dc:publisher>
          <dc:date>#{book_info[:year]}</dc:date>
          <dc:description>#{CGI.escapeHTML(book_info[:description])}</dc:description>
          <meta property="dcterms:modified">#{Time.now.strftime('%Y-%m-%dT%H:%M:%SZ')}</meta>
        </metadata>
        
        <manifest>
          <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"/>
          <item id="toc" href="toc.xhtml" media-type="application/xhtml+xml" properties="nav"/>
          <item id="css" href="css/style.css" media-type="text/css"/>
          <item id="cover" href="cover.xhtml" media-type="application/xhtml+xml"/>
          #{chapters.map { |ch| %Q[<item id="chapter#{ch[:number]}" href="#{ch[:filename]}" media-type="application/xhtml+xml"/>] }.join("\n          ")}
        </manifest>
        
        <spine toc="ncx">
          <itemref idref="cover"/>
          <itemref idref="toc"/>
          #{chapters.map { |ch| %Q[<itemref idref="chapter#{ch[:number]}"/>] }.join("\n          ")}
        </spine>
        
        <guide>
          <reference type="toc" title="Table of Contents" href="toc.xhtml"/>
          <reference type="text" title="Beginning" href="chapter01.xhtml"/>
        </guide>
      </package>
    XML
    
    File.write(File.join(temp_dir, 'OEBPS', 'content.opf'), content_opf)
  end

  def generate_toc_ncx(temp_dir, book_info, chapters)
    toc_ncx = <<~XML
      <?xml version="1.0" encoding="UTF-8"?>
      <ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
        <head>
          <meta name="dtb:uid" content="#{book_info[:isbn]}"/>
          <meta name="dtb:depth" content="1"/>
          <meta name="dtb:totalPageCount" content="0"/>
          <meta name="dtb:maxPageNumber" content="0"/>
        </head>
        
        <docTitle>
          <text>#{CGI.escapeHTML(book_info[:title])}</text>
        </docTitle>
        
        <navMap>
          <navPoint id="cover" playOrder="1">
            <navLabel>
              <text>Cover</text>
            </navLabel>
            <content src="cover.xhtml"/>
          </navPoint>
          
          <navPoint id="toc" playOrder="2">
            <navLabel>
              <text>Table of Contents</text>
            </navLabel>
            <content src="toc.xhtml"/>
          </navPoint>
          
          #{chapters.map.with_index { |ch, idx| 
            %Q[<navPoint id="navPoint-#{idx + 3}" playOrder="#{idx + 3}">
            <navLabel>
              <text>#{CGI.escapeHTML(ch[:title])}</text>
            </navLabel>
            <content src="#{ch[:filename]}"/>
          </navPoint>]
          }.join("\n          ")}
        </navMap>
      </ncx>
    XML
    
    File.write(File.join(temp_dir, 'OEBPS', 'toc.ncx'), toc_ncx)
  end

  def generate_toc_xhtml(temp_dir, book_info, chapters)
    toc_xhtml = <<~HTML
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE html>
      <html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
      <head>
        <title>Table of Contents</title>
        <link rel="stylesheet" type="text/css" href="css/style.css"/>
      </head>
      <body>
        <nav epub:type="toc">
          <h1>Table of Contents</h1>
          <ol>
            <li><a href="cover.xhtml">Cover</a></li>
            #{chapters.map { |ch| %Q[<li><a href="#{ch[:filename]}">#{CGI.escapeHTML(ch[:title])}</a></li>] }.join("\n            ")}
          </ol>
        </nav>
      </body>
      </html>
    HTML
    
    File.write(File.join(temp_dir, 'OEBPS', 'toc.xhtml'), toc_xhtml)
  end

  def generate_cover_xhtml(temp_dir, book_info)
    cover_xhtml = <<~HTML
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE html>
      <html xmlns="http://www.w3.org/1999/xhtml">
      <head>
        <title>#{CGI.escapeHTML(book_info[:title])}</title>
        <link rel="stylesheet" type="text/css" href="css/style.css"/>
      </head>
      <body class="cover">
        <div class="cover-content">
          <h1 class="title">#{CGI.escapeHTML(book_info[:title])}</h1>
          <h2 class="subtitle">#{CGI.escapeHTML(book_info[:subtitle])}</h2>
          <div class="spacer"></div>
          <p class="author">by #{CGI.escapeHTML(book_info[:author])}</p>
          <div class="spacer"></div>
          <p class="series">Book One of The Void Chronicles</p>
        </div>
      </body>
      </html>
    HTML
    
    File.write(File.join(temp_dir, 'OEBPS', 'cover.xhtml'), cover_xhtml)
  end

  def generate_chapters_xhtml(temp_dir, chapters)
    chapters.each do |chapter|
      chapter_xhtml = <<~HTML
        <?xml version="1.0" encoding="UTF-8"?>
        <!DOCTYPE html>
        <html xmlns="http://www.w3.org/1999/xhtml">
        <head>
          <title>#{CGI.escapeHTML(chapter[:title])}</title>
          <link rel="stylesheet" type="text/css" href="css/style.css"/>
        </head>
        <body>
          <h1>#{CGI.escapeHTML(chapter[:title])}</h1>
          #{chapter[:content]}
        </body>
        </html>
      HTML
      
      File.write(File.join(temp_dir, 'OEBPS', chapter[:filename]), chapter_xhtml)
    end
  end

  def generate_css(temp_dir)
    css = <<~CSS
      /* Basic styles for Kindle compatibility */
      body {
        font-family: Georgia, serif;
        font-size: 1em;
        line-height: 1.6;
        margin: 0;
        padding: 0;
        text-align: justify;
      }
      
      h1 {
        font-size: 1.5em;
        font-weight: bold;
        text-align: center;
        margin: 1em 0;
        page-break-before: always;
      }
      
      h2 {
        font-size: 1.3em;
        font-weight: bold;
        text-align: center;
        margin: 1em 0;
      }
      
      p {
        text-indent: 1.5em;
        margin: 0 0 0.5em 0;
      }
      
      p:first-of-type {
        text-indent: 0;
      }
      
      em {
        font-style: italic;
      }
      
      strong {
        font-weight: bold;
      }
      
      /* Cover page styles */
      .cover {
        text-align: center;
        height: 100vh;
        display: flex;
        flex-direction: column;
        justify-content: center;
      }
      
      .cover-content {
        margin: auto;
      }
      
      .title {
        font-size: 2.5em;
        font-weight: bold;
        margin: 0.5em 0;
        text-indent: 0;
        page-break-before: avoid;
      }
      
      .subtitle {
        font-size: 1.5em;
        font-weight: normal;
        margin: 0.5em 0;
      }
      
      .author {
        font-size: 1.3em;
        margin: 1em 0;
        font-style: italic;
      }
      
      .series {
        font-size: 1.1em;
        margin: 1em 0;
      }
      
      .spacer {
        height: 2em;
      }
      
      /* Table of contents */
      nav ol {
        list-style-type: none;
        padding-left: 0;
      }
      
      nav li {
        margin: 0.5em 0;
      }
      
      nav a {
        text-decoration: none;
        color: inherit;
      }
    CSS
    
    File.write(File.join(temp_dir, 'OEBPS', 'css', 'style.css'), css)
  end

  def create_epub_file(book_dir, temp_dir, book_info)
    epub_filename = File.join(book_dir, "#{book_info[:title].downcase.gsub(/\s+/, '_')}.epub")
    
    # Remove existing EPUB if present
    FileUtils.rm_f(epub_filename)
    
    # Create EPUB using zip command
    Dir.chdir(temp_dir) do
      # First add mimetype without compression
      system("zip -X0 '#{epub_filename}' mimetype")
      # Then add everything else with compression
      system("zip -Xr9D '#{epub_filename}' META-INF/ OEBPS/")
    end
    
    epub_filename
  end

  def generate_uuid
    "urn:uuid:#{SecureRandom.uuid}" rescue "urn:uuid:#{Time.now.to_i}-#{rand(1000000)}"
  end
end

# Run the converter
if __FILE__ == $0
  converter = LaTeXToEPUBConverter.new
  converter.convert
end