package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadBook(t *testing.T) {
	// Create a temporary directory with test book content
	tempDir := t.TempDir()

	// Create book.md
	bookContent := `# Test Book
	
## Chapter 1: Introduction

This is a test book.`

	bookPath := filepath.Join(tempDir, "book.md")
	if err := os.WriteFile(bookPath, []byte(bookContent), 0644); err != nil {
		t.Fatalf("Failed to create test book.md: %v", err)
	}

	// Create chapters directory
	chaptersDir := filepath.Join(tempDir, "chapters")
	if err := os.MkdirAll(chaptersDir, 0755); err != nil {
		t.Fatalf("Failed to create chapters directory: %v", err)
	}

	// Create a test chapter
	chapter1Content := `# Chapter 1: The Beginning

This is the first chapter of our test book.

It has multiple paragraphs.

* * *

And even a scene break.`

	chapter1Path := filepath.Join(chaptersDir, "chapter-01.md")
	if err := os.WriteFile(chapter1Path, []byte(chapter1Content), 0644); err != nil {
		t.Fatalf("Failed to create test chapter: %v", err)
	}

	// Test loading the book
	book, err := LoadBook(tempDir)
	if err != nil {
		t.Fatalf("LoadBook() error = %v", err)
	}

	if book == nil {
		t.Fatal("LoadBook() returned nil book")
	}

	if book.Title != "Void Reavers" {
		t.Errorf("Expected title 'Void Reavers', got '%s'", book.Title)
	}

	if book.Author != "Captain J. Starwind" {
		t.Errorf("Expected author 'Captain J. Starwind', got '%s'", book.Author)
	}

	if len(book.Chapters) != 1 {
		t.Errorf("Expected 1 chapter, got %d", len(book.Chapters))
	}

	if len(book.Chapters) > 0 {
		chapter := book.Chapters[0]
		if chapter.Title != "Chapter 1: The Beginning" {
			t.Errorf("Expected chapter title 'Chapter 1: The Beginning', got '%s'", chapter.Title)
		}

		if chapter.Content == "" {
			t.Error("Chapter content is empty")
		}
	}
}

func TestLoadBookMissingDirectory(t *testing.T) {
	book, err := LoadBook("/nonexistent/directory")

	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}

	if book != nil {
		t.Error("Expected nil book for nonexistent directory")
	}
}

func TestLoadBookEmptyChapters(t *testing.T) {
	// Create a temporary directory with book.md but no chapters
	tempDir := t.TempDir()

	bookContent := `# Empty Book

A book with no chapters.`

	bookPath := filepath.Join(tempDir, "book.md")
	if err := os.WriteFile(bookPath, []byte(bookContent), 0644); err != nil {
		t.Fatalf("Failed to create test book.md: %v", err)
	}

	// Create empty chapters directory
	chaptersDir := filepath.Join(tempDir, "chapters")
	if err := os.MkdirAll(chaptersDir, 0755); err != nil {
		t.Fatalf("Failed to create chapters directory: %v", err)
	}

	book, err := LoadBook(tempDir)
	if err == nil {
		t.Fatal("LoadBook() should return error for empty chapters directory")
	}

	if book != nil {
		t.Fatal("LoadBook() should return nil book when no chapters found")
	}

	expectedErrMsg := "no chapter files found"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error containing '%s', got '%v'", expectedErrMsg, err)
	}
}

func TestParseMarkdownChapter(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedTitle   string
		expectedContent string
	}{
		{
			name:            "extracts chapter title and content",
			input:           "# Chapter 1: The Beginning\n\nThis is the first paragraph.\n\nThis is the second paragraph.",
			expectedTitle:   "Chapter 1: The Beginning",
			expectedContent: "This is the first paragraph.\n\nThis is the second paragraph.",
		},
		{
			name:            "handles content with italics",
			input:           "# The Ship\n\nThe *Void Reaver* sailed through space.",
			expectedTitle:   "The Ship",
			expectedContent: "The *Void Reaver* sailed through space.",
		},
		{
			name:            "handles content with bold",
			input:           "# Alert\n\n**Warning:** Pirates ahead!",
			expectedTitle:   "Alert",
			expectedContent: "**Warning:** Pirates ahead!",
		},
		{
			name:            "preserves scene breaks",
			input:           "# Scene Break Test\n\nBefore the break.\n\n* * *\n\nAfter the break.",
			expectedTitle:   "Scene Break Test",
			expectedContent: "Before the break.\n\n* * *\n\nAfter the break.",
		},
		{
			name:            "handles missing title",
			input:           "This is content without a title header.",
			expectedTitle:   "",
			expectedContent: "",
		},
		{
			name:            "ignores subsequent headers",
			input:           "# Main Title\n\nSome content.\n\n## Subsection\n\nMore content.",
			expectedTitle:   "Main Title",
			expectedContent: "Some content.\n\n## Subsection\n\nMore content.",
		},
		{
			name:            "trims whitespace properly",
			input:           "# Trimmed Title   \n\n\n\nContent with extra spaces.\n\n\n",
			expectedTitle:   "Trimmed Title",
			expectedContent: "Content with extra spaces.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMarkdownChapter(tt.input)
			if result.Title != tt.expectedTitle {
				t.Errorf("parseMarkdownChapter() title = %q, want %q", result.Title, tt.expectedTitle)
			}
			if result.Content != tt.expectedContent {
				t.Errorf("parseMarkdownChapter() content = %q, want %q", result.Content, tt.expectedContent)
			}
		})
	}
}
