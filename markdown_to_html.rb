#!/usr/bin/env ruby

require 'fileutils'

class MarkdownToHtmlConverter
  def initialize(base_dir = Dir.pwd)
    @base_dir = base_dir
    @markdown_dir = File.join(@base_dir, 'markdown')
    @html_dir = File.join(@base_dir, 'html')
  end

  def convert
    puts "Markdown to HTML Converter for Void Reavers"
    puts "=" * 60

    # Check if markdown directory exists
    unless Dir.exist?(@markdown_dir)
      puts "Error: markdown directory not found!"
      puts "Please run convert_to_md.rb first to generate markdown files."
      return
    end

    # Create output directory
    FileUtils.mkdir_p(@html_dir)

    # Convert complete book
    convert_complete_book_to_html
    
    # Create CSS file
    create_css_file

    puts "\nConversion complete!"
    puts "HTML files saved in: #{@html_dir}"
    puts "\nTo view the book, open: #{File.join(@html_dir, 'void_reavers_complete.html')}"
    puts "\nNote: You can print this HTML file to PDF using your web browser's print function."
  end

  private

  def convert_complete_book_to_html
    puts "\nConverting complete book to HTML..."
    
    md_file = File.join(@markdown_dir, 'void_reavers_complete.md')
    unless File.exist?(md_file)
      puts "Error: void_reavers_complete.md not found!"
      return
    end

    content = File.read(md_file, encoding: 'utf-8')
    html_content = markdown_to_html(content)
    
    # Create full HTML document
    full_html = create_html_document(html_content)
    
    # Save HTML file
    html_file = File.join(@html_dir, 'void_reavers_complete.html')
    File.write(html_file, full_html, encoding: 'utf-8')
    
    puts "  ✓ Complete book converted to HTML"
  end

  def markdown_to_html(markdown)
    html = markdown
    
    # Convert headers
    html = html.gsub(/^### (.+)$/, '<h3>\1</h3>')
    html = html.gsub(/^## (.+)$/, '<h2>\1</h2>')
    html = html.gsub(/^# (.+)$/, '<h1>\1</h1>')
    
    # Convert bold and italic
    html = html.gsub(/\*\*([^*]+)\*\*/, '<strong>\1</strong>')
    html = html.gsub(/\*([^*]+)\*/, '<em>\1</em>')
    
    # Convert horizontal rules
    html = html.gsub(/^---$/, '<hr>')
    
    # Convert paragraphs
    paragraphs = html.split(/\n\n+/)
    html = paragraphs.map do |para|
      para = para.strip
      if para.start_with?('<h', '<hr')
        para
      elsif para.empty?
        ''
      else
        "<p>#{para}</p>"
      end
    end.join("\n\n")
    
    # Fix quotes
    html = html.gsub('"', '&ldquo;').gsub('"', '&rdquo;')
    html = html.gsub("'", '&lsquo;').gsub("'", '&rsquo;')
    
    html
  end

  def create_html_document(content)
    <<~HTML
      <!DOCTYPE html>
      <html lang="en">
      <head>
          <meta charset="UTF-8">
          <meta name="viewport" content="width=device-width, initial-scale=1.0">
          <title>Void Reavers - A Tale of Space Pirates and Cosmic Plunder</title>
          <link rel="stylesheet" href="style.css">
      </head>
      <body>
          <div class="container">
              #{content}
          </div>
          
          <script>
              // Add print-friendly features
              window.onload = function() {
                  // Add print button
                  const printBtn = document.createElement('button');
                  printBtn.textContent = 'Print / Save as PDF';
                  printBtn.className = 'print-button';
                  printBtn.onclick = function() { window.print(); };
                  document.body.insertBefore(printBtn, document.body.firstChild);
                  
                  // Add table of contents
                  const toc = document.createElement('div');
                  toc.className = 'toc';
                  toc.innerHTML = '<h2>Table of Contents</h2>';
                  const tocList = document.createElement('ol');
                  
                  document.querySelectorAll('h1').forEach((h1, index) => {
                      if (index > 0) { // Skip the main title
                          const li = document.createElement('li');
                          const a = document.createElement('a');
                          a.href = '#chapter' + index;
                          a.textContent = h1.textContent;
                          h1.id = 'chapter' + index;
                          li.appendChild(a);
                          tocList.appendChild(li);
                      }
                  });
                  
                  toc.appendChild(tocList);
                  const firstH1 = document.querySelector('h1');
                  firstH1.parentNode.insertBefore(toc, firstH1.nextSibling.nextSibling);
              };
          </script>
      </body>
      </html>
    HTML
  end

  def create_css_file
    css_content = <<~CSS
      /* Void Reavers Book Styles */
      
      body {
          font-family: Georgia, 'Times New Roman', serif;
          line-height: 1.8;
          color: #333;
          background-color: #f9f9f9;
          margin: 0;
          padding: 20px;
      }
      
      .container {
          max-width: 800px;
          margin: 0 auto;
          background-color: white;
          padding: 40px;
          box-shadow: 0 0 20px rgba(0,0,0,0.1);
      }
      
      h1 {
          font-size: 2.5em;
          color: #2c3e50;
          margin-bottom: 0.5em;
          page-break-before: always;
      }
      
      h1:first-of-type {
          page-break-before: avoid;
          text-align: center;
          font-size: 3em;
      }
      
      h2 {
          font-size: 2em;
          color: #34495e;
          margin-top: 1.5em;
          margin-bottom: 0.5em;
          text-align: center;
      }
      
      h3 {
          font-size: 1.5em;
          color: #7f8c8d;
          margin-top: 1em;
      }
      
      p {
          text-align: justify;
          margin-bottom: 1em;
          text-indent: 1.5em;
      }
      
      p:first-of-type {
          text-indent: 0;
      }
      
      hr {
          border: none;
          text-align: center;
          margin: 2em 0;
      }
      
      hr::before {
          content: "* * *";
          color: #bdc3c7;
          font-size: 1.2em;
      }
      
      em {
          font-style: italic;
      }
      
      strong {
          font-weight: bold;
      }
      
      .toc {
          background-color: #ecf0f1;
          padding: 20px;
          margin: 2em 0;
          border-radius: 5px;
      }
      
      .toc h2 {
          color: #2c3e50;
          margin-top: 0;
      }
      
      .toc ol {
          counter-reset: chapter;
          list-style: none;
          padding-left: 0;
      }
      
      .toc li {
          counter-increment: chapter;
          margin-bottom: 0.5em;
      }
      
      .toc li::before {
          content: counter(chapter) ". ";
          font-weight: bold;
          margin-right: 0.5em;
      }
      
      .toc a {
          color: #3498db;
          text-decoration: none;
      }
      
      .toc a:hover {
          text-decoration: underline;
      }
      
      .print-button {
          position: fixed;
          top: 20px;
          right: 20px;
          padding: 10px 20px;
          background-color: #3498db;
          color: white;
          border: none;
          border-radius: 5px;
          cursor: pointer;
          font-size: 16px;
          box-shadow: 0 2px 5px rgba(0,0,0,0.2);
          z-index: 1000;
      }
      
      .print-button:hover {
          background-color: #2980b9;
      }
      
      /* Print styles */
      @media print {
          body {
              background-color: white;
              padding: 0;
          }
          
          .container {
              box-shadow: none;
              padding: 0;
              max-width: 100%;
          }
          
          .print-button {
              display: none;
          }
          
          h1 {
              page-break-before: always;
              page-break-after: avoid;
          }
          
          h1:first-of-type {
              page-break-before: avoid;
          }
          
          p {
              orphans: 3;
              widows: 3;
          }
          
          .toc {
              page-break-after: always;
          }
      }
      
      @page {
          margin: 1in;
      }
    CSS

    css_file = File.join(@html_dir, 'style.css')
    File.write(css_file, css_content)
    puts "  ✓ CSS file created"
  end
end

# Run the converter if this script is executed directly
if __FILE__ == $0
  converter = MarkdownToHtmlConverter.new
  converter.convert
end