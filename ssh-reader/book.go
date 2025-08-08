package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/muesli/reflow/wordwrap"
)

type Chapter struct {
	Title   string
	Content string
}

type Book struct {
	Title    string
	Author   string
	Chapters []Chapter
}

func LoadBook(bookDir string) (*Book, error) {
	// Try to load from markdown first, then LaTeX
	book, err := loadFromMarkdown(bookDir)
	if err != nil {
		book, err = loadFromLaTeX(bookDir)
		if err != nil {
			return nil, fmt.Errorf("failed to load book from both markdown and LaTeX: %v", err)
		}
	}

	return book, nil
}

func loadFromMarkdown(bookDir string) (*Book, error) {
	markdownDir := filepath.Join(bookDir, "markdown")
	
	// Try to find chapter files
	files, err := ioutil.ReadDir(markdownDir)
	if err != nil {
		return nil, fmt.Errorf("could not read markdown directory: %v", err)
	}

	var chapterFiles []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "chapter") && strings.HasSuffix(file.Name(), ".md") {
			chapterFiles = append(chapterFiles, file.Name())
		}
	}

	if len(chapterFiles) == 0 {
		return nil, fmt.Errorf("no chapter files found in %s", markdownDir)
	}

	sort.Strings(chapterFiles)

	var chapters []Chapter
	for _, filename := range chapterFiles {
		content, err := ioutil.ReadFile(filepath.Join(markdownDir, filename))
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %v", filename, err)
		}

		chapter := parseMarkdownChapter(string(content))
		if chapter.Title == "" {
			// Generate title from filename if not found
			chapter.Title = fmt.Sprintf("Chapter %d", len(chapters)+1)
		}
		chapters = append(chapters, chapter)
	}

	return &Book{
		Title:    "Void Reavers",
		Author:   "Captain J. Starwind",
		Chapters: chapters,
	}, nil
}

func loadFromLaTeX(bookDir string) (*Book, error) {
	// Try to find chapter files
	files, err := ioutil.ReadDir(bookDir)
	if err != nil {
		return nil, fmt.Errorf("could not read book directory: %v", err)
	}

	var chapterFiles []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "chapter") && strings.HasSuffix(file.Name(), ".tex") {
			chapterFiles = append(chapterFiles, file.Name())
		}
	}

	if len(chapterFiles) == 0 {
		return nil, fmt.Errorf("no chapter files found in %s", bookDir)
	}

	sort.Strings(chapterFiles)

	var chapters []Chapter
	for _, filename := range chapterFiles {
		content, err := ioutil.ReadFile(filepath.Join(bookDir, filename))
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %v", filename, err)
		}

		chapter := parseLaTeXChapter(string(content))
		if chapter.Title == "" {
			// Generate title from filename if not found
			chapter.Title = fmt.Sprintf("Chapter %d", len(chapters)+1)
		}
		chapters = append(chapters, chapter)
	}

	return &Book{
		Title:    "Void Reavers",
		Author:   "Captain J. Starwind",
		Chapters: chapters,
	}, nil
}

func parseMarkdownChapter(content string) Chapter {
	lines := strings.Split(content, "\n")
	var title string
	var contentLines []string
	
	titleFound := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") && !titleFound {
			title = strings.TrimSpace(trimmed[2:])
			titleFound = true
		} else if titleFound {
			contentLines = append(contentLines, line)
		}
	}

	// Clean up content
	content = strings.Join(contentLines, "\n")
	content = strings.TrimSpace(content)

	return Chapter{
		Title:   title,
		Content: content,
	}
}

func parseLaTeXChapter(content string) Chapter {
	// Extract chapter title
	chapterRegex := regexp.MustCompile(`\\chapter\{([^}]+)\}`)
	matches := chapterRegex.FindStringSubmatch(content)
	
	var title string
	if len(matches) > 1 {
		title = matches[1]
	}

	// Convert LaTeX to plain text
	content = convertLaTeXToPlainText(content)

	return Chapter{
		Title:   title,
		Content: content,
	}
}

func convertLaTeXToPlainText(content string) string {
	// Remove chapter command
	content = regexp.MustCompile(`\\chapter\{[^}]+\}`).ReplaceAllString(content, "")
	
	// Convert italics
	content = regexp.MustCompile(`\\textit\{([^}]+)\}`).ReplaceAllString(content, "*$1*")
	
	// Convert bold
	content = regexp.MustCompile(`\\textbf\{([^}]+)\}`).ReplaceAllString(content, "**$1**")
	
	// Convert double quotes
	content = regexp.MustCompile("``([^']+)''").ReplaceAllString(content, "\"$1\"")
	
	// Convert single quotes
	content = regexp.MustCompile("`([^']+)'").ReplaceAllString(content, "'$1'")
	
	// Convert em-dashes
	content = strings.ReplaceAll(content, "---", "—")
	
	// Remove LaTeX escapes
	replacements := map[string]string{
		"\\%": "%",
		"\\$": "$",
		"\\&": "&",
		"\\_": "_",
		"\\#": "#",
	}
	
	for latex, plain := range replacements {
		content = strings.ReplaceAll(content, latex, plain)
	}
	
	// Clean up excessive newlines
	content = regexp.MustCompile(`\n{3,}`).ReplaceAllString(content, "\n\n")
	content = strings.TrimSpace(content)
	
	return content
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		width = 80
	}
	
	paragraphs := strings.Split(text, "\n\n")
	var wrappedLines []string
	
	for i, paragraph := range paragraphs {
		if strings.TrimSpace(paragraph) == "" {
			continue
		}
		
		// Handle special formatting
		if strings.HasPrefix(paragraph, "*") && strings.HasSuffix(paragraph, "*") {
			// Italic text
			wrapped := wordwrap.String(paragraph, width)
			wrappedLines = append(wrappedLines, strings.Split(wrapped, "\n")...)
		} else if strings.HasPrefix(paragraph, "**") && strings.HasSuffix(paragraph, "**") {
			// Bold text
			wrapped := wordwrap.String(paragraph, width)
			wrappedLines = append(wrappedLines, strings.Split(wrapped, "\n")...)
		} else {
			// Regular paragraph
			wrapped := wordwrap.String(paragraph, width)
			wrappedLines = append(wrappedLines, strings.Split(wrapped, "\n")...)
		}
		
		// Add spacing between paragraphs (except for the last one)
		if i < len(paragraphs)-1 {
			wrappedLines = append(wrappedLines, "")
		}
	}
	
	return wrappedLines
}