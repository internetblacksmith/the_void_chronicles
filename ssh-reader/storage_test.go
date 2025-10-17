package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStorageCleanup(t *testing.T) {
	t.Run("removes old files", func(t *testing.T) {
		tempDir := t.TempDir()
		pm := &ProgressManager{dataDir: tempDir}

		// Create some test files with different ages
		oldFile := filepath.Join(tempDir, "olduser.json")
		newFile := filepath.Join(tempDir, "newuser.json")

		// Create old file
		os.WriteFile(oldFile, []byte("{}"), 0644)
		// Modify its time to be 100 days ago
		oldTime := time.Now().Add(-100 * 24 * time.Hour)
		os.Chtimes(oldFile, oldTime, oldTime)

		// Create new file
		os.WriteFile(newFile, []byte("{}"), 0644)

		// Run cleanup
		cleaned, err := pm.CleanupOldFiles()
		if err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		// Should have cleaned 1 file
		if cleaned != 1 {
			t.Errorf("Expected 1 file cleaned, got %d", cleaned)
		}

		// Old file should be gone
		if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
			t.Error("Old file should have been deleted")
		}

		// New file should still exist
		if _, err := os.Stat(newFile); os.IsNotExist(err) {
			t.Error("New file should not have been deleted")
		}
	})

	t.Run("skips directories", func(t *testing.T) {
		tempDir := t.TempDir()
		pm := &ProgressManager{dataDir: tempDir}

		// Create a subdirectory
		subDir := filepath.Join(tempDir, "subdir")
		os.Mkdir(subDir, 0755)

		// Make it old
		oldTime := time.Now().Add(-100 * 24 * time.Hour)
		os.Chtimes(subDir, oldTime, oldTime)

		// Run cleanup
		cleaned, err := pm.CleanupOldFiles()
		if err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		// Should not have cleaned the directory
		if cleaned != 0 {
			t.Errorf("Expected 0 files cleaned (directory should be skipped), got %d", cleaned)
		}

		// Directory should still exist
		if _, err := os.Stat(subDir); os.IsNotExist(err) {
			t.Error("Directory should not have been deleted")
		}
	})
}

func TestProgressValidation(t *testing.T) {
	pm := &ProgressManager{dataDir: t.TempDir()}

	t.Run("trims excess bookmarks", func(t *testing.T) {
		progress := &UserProgress{
			Username: "testuser",
			Bookmarks: make([]Bookmark, 25), // More than MaxBookmarksPerUser
		}

		// Fill bookmarks with test data
		for i := range progress.Bookmarks {
			progress.Bookmarks[i] = Bookmark{
				Chapter: i,
				Note:    "Test bookmark",
			}
		}

		err := pm.ValidateAndTrimProgress(progress)
		if err != nil {
			t.Fatalf("Validation failed: %v", err)
		}

		if len(progress.Bookmarks) != MaxBookmarksPerUser {
			t.Errorf("Expected %d bookmarks after trim, got %d", MaxBookmarksPerUser, len(progress.Bookmarks))
		}

		// Should keep the most recent bookmarks (last 20)
		if progress.Bookmarks[0].Chapter != 5 {
			t.Errorf("Expected first bookmark to be from position 5, got %d", progress.Bookmarks[0].Chapter)
		}
	})

	t.Run("trims long bookmark notes", func(t *testing.T) {
		longNote := strings.Repeat("a", MaxNoteLength+50)
		progress := &UserProgress{
			Username: "testuser",
			Bookmarks: []Bookmark{
				{
					Chapter: 1,
					Note:    longNote,
				},
			},
		}

		err := pm.ValidateAndTrimProgress(progress)
		if err != nil {
			t.Fatalf("Validation failed: %v", err)
		}

		if len(progress.Bookmarks[0].Note) != MaxNoteLength {
			t.Errorf("Expected note length %d, got %d", MaxNoteLength, len(progress.Bookmarks[0].Note))
		}
	})

	t.Run("enforces file size limit", func(t *testing.T) {
		progress := &UserProgress{
			Username: "testuser",
			Bookmarks: make([]Bookmark, 100), // Way too many
		}

		// Fill with large notes to exceed size limit
		for i := range progress.Bookmarks {
			progress.Bookmarks[i] = Bookmark{
				Chapter: i,
				Note:    strings.Repeat("x", 200),
			}
		}

		err := pm.ValidateAndTrimProgress(progress)
		if err != nil {
			t.Fatalf("Validation failed: %v", err)
		}

		// Should have removed bookmarks to fit under MaxFileSize
		data, _ := json.Marshal(progress)
		if len(data) > MaxFileSize {
			t.Errorf("File size %d exceeds max %d after validation", len(data), MaxFileSize)
		}
	})
}

func TestStorageStats(t *testing.T) {
	t.Run("calculates storage statistics", func(t *testing.T) {
		tempDir := t.TempDir()
		pm := &ProgressManager{dataDir: tempDir}

		// Create test files
		files := []struct {
			name string
			size int
		}{
			{"user1.json", 500},
			{"user2.json", 1000},
			{"user3.json", 1500},
		}

		for _, f := range files {
			content := make([]byte, f.size)
			os.WriteFile(filepath.Join(tempDir, f.name), content, 0644)
		}

		stats, err := pm.GetStorageStats()
		if err != nil {
			t.Fatalf("Failed to get stats: %v", err)
		}

		if stats.FileCount != 3 {
			t.Errorf("Expected 3 files, got %d", stats.FileCount)
		}

		expectedTotal := int64(3000)
		if stats.TotalSizeBytes != expectedTotal {
			t.Errorf("Expected total size %d, got %d", expectedTotal, stats.TotalSizeBytes)
		}

		expectedAvg := int64(1000)
		if stats.AvgFileSize != expectedAvg {
			t.Errorf("Expected average size %d, got %d", expectedAvg, stats.AvgFileSize)
		}

		// Check capacity calculation (1GB / 1000 bytes avg)
		expectedCapacity := int64(1073741824 / 1000)
		if stats.EstimatedCapacity != expectedCapacity {
			t.Errorf("Expected capacity %d, got %d", expectedCapacity, stats.EstimatedCapacity)
		}
	})

	t.Run("handles empty directory", func(t *testing.T) {
		tempDir := t.TempDir()
		pm := &ProgressManager{dataDir: tempDir}

		stats, err := pm.GetStorageStats()
		if err != nil {
			t.Fatalf("Failed to get stats: %v", err)
		}

		if stats.FileCount != 0 {
			t.Errorf("Expected 0 files, got %d", stats.FileCount)
		}

		if stats.TotalSizeBytes != 0 {
			t.Errorf("Expected 0 total size, got %d", stats.TotalSizeBytes)
		}
	})
}