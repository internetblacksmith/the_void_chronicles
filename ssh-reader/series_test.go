package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestGetSeriesBooks(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	t.Run("loads series from JSON file", func(t *testing.T) {
		// Create test series.json
		testSeries := SeriesInfo{
			Series: "Test Series",
			Books: []BookInfo{
				{
					Number:    1,
					Title:     "Test Book 1",
					Subtitle:  "Test Subtitle",
					Status:    "Available",
					Available: true,
					Summary:   "Test summary",
				},
				{
					Number:    2,
					Title:     "Test Book 2",
					Subtitle:  "Another Subtitle",
					Status:    "Coming Soon",
					Available: false,
					Summary:   "Another summary",
				},
			},
		}

		data, _ := json.Marshal(testSeries)
		os.WriteFile("series.json", data, 0644)

		books := getSeriesBooks()

		if len(books) != 2 {
			t.Errorf("Expected 2 books, got %d", len(books))
		}

		if books[0].Title != "Test Book 1" {
			t.Errorf("Expected first book title 'Test Book 1', got '%s'", books[0].Title)
		}

		if !books[0].Available {
			t.Error("Expected first book to be available")
		}

		if books[1].Available {
			t.Error("Expected second book to not be available")
		}
	})

	t.Run("returns fallback when JSON file missing", func(t *testing.T) {
		// Remove any existing series.json
		os.Remove("series.json")
		os.Remove("ssh-reader/series.json")
		os.Remove("../series.json")

		books := getSeriesBooks()

		if len(books) != 1 {
			t.Errorf("Expected 1 fallback book, got %d", len(books))
		}

		if books[0].Title != "Void Reavers" {
			t.Errorf("Expected fallback book title 'Void Reavers', got '%s'", books[0].Title)
		}
	})

	t.Run("handles malformed JSON gracefully", func(t *testing.T) {
		os.WriteFile("series.json", []byte("not valid json"), 0644)

		books := getSeriesBooks()

		// Should return fallback
		if len(books) != 1 {
			t.Errorf("Expected 1 fallback book on JSON error, got %d", len(books))
		}
	})
}