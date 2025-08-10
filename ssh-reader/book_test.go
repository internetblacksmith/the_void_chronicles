package main

import (
	"os"
	"path/filepath"
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
	
	if book.Title != "Test Book" {
		t.Errorf("Expected title 'Test Book', got '%s'", book.Title)
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
	if err != nil {
		t.Fatalf("LoadBook() error = %v", err)
	}
	
	if book == nil {
		t.Fatal("LoadBook() returned nil book")
	}
	
	if len(book.Chapters) != 0 {
		t.Errorf("Expected 0 chapters for empty directory, got %d", len(book.Chapters))
	}
}

func TestProcessMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "preserves italics",
			input:    "This is *italic* text",
			expected: "This is *italic* text",
		},
		{
			name:     "preserves bold",
			input:    "This is **bold** text",
			expected: "This is **bold** text",
		},
		{
			name:     "preserves scene breaks",
			input:    "Before\n\n* * *\n\nAfter",
			expected: "Before\n\n* * *\n\nAfter",
		},
		{
			name:     "handles multiple paragraphs",
			input:    "First paragraph.\n\nSecond paragraph.\n\nThird paragraph.",
			expected: "First paragraph.\n\nSecond paragraph.\n\nThird paragraph.",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("processMarkdown() = %v, want %v", result, tt.expected)
			}
		})
	}
}