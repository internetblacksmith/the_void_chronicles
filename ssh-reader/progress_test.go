package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProgressManager(t *testing.T) {
	// Create a temporary directory for test data
	tempDir := t.TempDir()

	t.Run("creates and loads user progress", func(t *testing.T) {
		pm := &ProgressManager{dataDir: tempDir}

		// Create test progress
		progress := &UserProgress{
			Username:       "testuser",
			CurrentChapter: 5,
			ScrollOffset:   100,
			LastRead:       time.Now(),
			ChapterProgress: map[int]bool{
				1: true,
				2: true,
				3: true,
				4: true,
				5: false,
			},
			Bookmarks: []Bookmark{
				{
					Chapter:      3,
					ScrollOffset: 50,
					Note:         "Important scene",
					Created:      time.Now(),
				},
			},
			ReadingTime: time.Hour,
		}

		// Save progress
		err := pm.SaveProgress(progress)
		if err != nil {
			t.Fatalf("Failed to save progress: %v", err)
		}

		// Verify file was created
		filename := filepath.Join(tempDir, "testuser.json")
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Error("Progress file was not created")
		}

		// Load progress back
		loaded, err := pm.LoadProgress("testuser")
		if err != nil {
			t.Fatalf("Failed to load progress: %v", err)
		}

		// Verify loaded data
		if loaded.CurrentChapter != 5 {
			t.Errorf("Expected current chapter 5, got %d", loaded.CurrentChapter)
		}

		if loaded.ScrollOffset != 100 {
			t.Errorf("Expected scroll offset 100, got %d", loaded.ScrollOffset)
		}

		if len(loaded.Bookmarks) != 1 {
			t.Errorf("Expected 1 bookmark, got %d", len(loaded.Bookmarks))
		}

		if loaded.Bookmarks[0].Note != "Important scene" {
			t.Errorf("Expected bookmark note 'Important scene', got '%s'", loaded.Bookmarks[0].Note)
		}
	})

	t.Run("handles missing user gracefully", func(t *testing.T) {
		pm := &ProgressManager{dataDir: tempDir}

		progress, err := pm.LoadProgress("nonexistent")
		
		// Should return empty progress, not error
		if err != nil {
			t.Errorf("Expected no error for missing user, got: %v", err)
		}

		if progress == nil {
			t.Fatal("Expected empty progress object, got nil")
		}

		if progress.CurrentChapter != 0 {
			t.Errorf("Expected chapter 0 for new user, got %d", progress.CurrentChapter)
		}
	})

	t.Run("adds and removes bookmarks", func(t *testing.T) {
		progress := &UserProgress{
			Username:        "bookmarkuser",
			ChapterProgress: make(map[int]bool),
			Bookmarks:       []Bookmark{},
		}

		// Add bookmark
		progress.AddBookmark(3, 100, "Test note")

		if len(progress.Bookmarks) != 1 {
			t.Errorf("Expected 1 bookmark after adding, got %d", len(progress.Bookmarks))
		}

		// Remove bookmark
		progress.RemoveBookmark(0)

		if len(progress.Bookmarks) != 0 {
			t.Errorf("Expected 0 bookmarks after removing, got %d", len(progress.Bookmarks))
		}
	})

	t.Run("tracks chapter completion", func(t *testing.T) {
		progress := &UserProgress{
			Username:        "completionuser",
			ChapterProgress: make(map[int]bool),
		}

		// Mark chapters as complete
		progress.MarkChapterComplete(1)
		progress.MarkChapterComplete(2)
		progress.MarkChapterComplete(3)

		if !progress.IsChapterComplete(1) {
			t.Error("Expected chapter 1 to be complete")
		}

		if !progress.IsChapterComplete(2) {
			t.Error("Expected chapter 2 to be complete")
		}

		if progress.IsChapterComplete(4) {
			t.Error("Expected chapter 4 to not be complete")
		}
	})

	t.Run("calculates completion percentage", func(t *testing.T) {
		progress := &UserProgress{
			ChapterProgress: map[int]bool{
				0: true,
				1: true,
				2: false,
				3: false,
			},
		}

		percentage := progress.GetCompletionPercentage(4)
		
		// 2 out of 4 chapters = 50%
		if percentage != 50.0 {
			t.Errorf("Expected 50%% completion, got %.1f%%", percentage)
		}
	})

	t.Run("generates reading stats", func(t *testing.T) {
		progress := &UserProgress{
			CurrentChapter: 5,
			ChapterProgress: map[int]bool{
				1: true,
				2: true,
			},
			Bookmarks:   make([]Bookmark, 3),
			ReadingTime: 2 * time.Hour,
		}

		stats := progress.GetReadingStats()

		if stats["chapters_completed"] != 2 {
			t.Errorf("Expected 2 chapters completed, got %v", stats["chapters_completed"])
		}

		if stats["current_chapter"] != 6 { // 1-based display
			t.Errorf("Expected current chapter 6 (1-based), got %v", stats["current_chapter"])
		}

		if stats["bookmarks_count"] != 3 {
			t.Errorf("Expected 3 bookmarks, got %v", stats["bookmarks_count"])
		}
	})
}