// Copyright (C) 2024 Paolo Fabbri
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	// Load from markdown chapters
	book, err := loadFromMarkdown(bookDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load book: %v", err)
	}

	return book, nil
}

func loadFromMarkdown(bookDir string) (*Book, error) {
	// Read chapters from the chapters subdirectory
	chaptersDir := filepath.Join(bookDir, "chapters")
	files, err := ioutil.ReadDir(chaptersDir)
	if err != nil {
		return nil, fmt.Errorf("could not read chapters directory: %v", err)
	}

	var chapterFiles []string
	for _, file := range files {
		// Support both "chapter" and "chapter-" prefixes
		if (strings.HasPrefix(file.Name(), "chapter") || strings.HasPrefix(file.Name(), "chapter-")) && strings.HasSuffix(file.Name(), ".md") {
			chapterFiles = append(chapterFiles, file.Name())
		}
	}

	if len(chapterFiles) == 0 {
		return nil, fmt.Errorf("no chapter files found in %s", chaptersDir)
	}

	sort.Strings(chapterFiles)

	var chapters []Chapter
	for _, filename := range chapterFiles {
		content, err := ioutil.ReadFile(filepath.Join(chaptersDir, filename))
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