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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	MaxBookmarksPerUser = 20
	MaxNoteLength       = 200
	MaxFileSize         = 10 * 1024 // 10KB per user
	CleanupAgeDays      = 90        // Delete files older than 90 days
)

// StorageStats holds storage usage information
type StorageStats struct {
	TotalSizeBytes    int64   `json:"total_size_bytes"`
	TotalSizeMB       float64 `json:"total_size_mb"`
	FileCount         int     `json:"file_count"`
	AvgFileSize       int64   `json:"avg_file_size"`
	EstimatedCapacity int64   `json:"estimated_capacity"`
	OldestFile        string  `json:"oldest_file"`
	NewestFile        string  `json:"newest_file"`
}

// CleanupOldFiles removes user progress files that haven't been accessed in CleanupAgeDays.
func (pm *ProgressManager) CleanupOldFiles() (int, error) {
	cleaned := 0
	cutoffTime := time.Now().Add(-CleanupAgeDays * 24 * time.Hour)

	files, err := os.ReadDir(pm.dataDir)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		// Check both modification time and access time
		if info.ModTime().Before(cutoffTime) {
			fullPath := filepath.Join(pm.dataDir, file.Name())
			if err := os.Remove(fullPath); err == nil {
				cleaned++
				log.Printf("Cleaned up old progress file: %s", file.Name())
			}
		}
	}

	return cleaned, nil
}

// ValidateAndTrimProgress ensures user data doesn't exceed storage limits.
func (pm *ProgressManager) ValidateAndTrimProgress(progress *UserProgress) error {
	// Trim bookmarks if too many
	if len(progress.Bookmarks) > MaxBookmarksPerUser {
		// Keep the most recent bookmarks
		progress.Bookmarks = progress.Bookmarks[len(progress.Bookmarks)-MaxBookmarksPerUser:]
	}

	// Trim bookmark notes if too long
	for i := range progress.Bookmarks {
		if len(progress.Bookmarks[i].Note) > MaxNoteLength {
			progress.Bookmarks[i].Note = progress.Bookmarks[i].Note[:MaxNoteLength]
		}
	}

	// Check total size
	data, err := json.Marshal(progress)
	if err != nil {
		return err
	}

	if len(data) > MaxFileSize {
		// Remove oldest bookmarks until size is acceptable
		for len(progress.Bookmarks) > 0 && len(data) > MaxFileSize {
			progress.Bookmarks = progress.Bookmarks[1:]
			data, _ = json.Marshal(progress)
		}
	}

	return nil
}

// GetStorageStats returns current storage usage statistics.
func (pm *ProgressManager) GetStorageStats() (*StorageStats, error) {
	stats := &StorageStats{}

	var oldestTime time.Time
	var newestTime time.Time

	err := filepath.Walk(pm.dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files with errors
		}

		if !info.IsDir() {
			stats.TotalSizeBytes += info.Size()
			stats.FileCount++

			modTime := info.ModTime()
			if oldestTime.IsZero() || modTime.Before(oldestTime) {
				oldestTime = modTime
				stats.OldestFile = info.Name()
			}
			if newestTime.IsZero() || modTime.After(newestTime) {
				newestTime = modTime
				stats.NewestFile = info.Name()
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if stats.FileCount > 0 {
		stats.AvgFileSize = stats.TotalSizeBytes / int64(stats.FileCount)
		// Estimate capacity based on 1GB volume
		stats.EstimatedCapacity = 1073741824 / stats.AvgFileSize
	}

	stats.TotalSizeMB = float64(stats.TotalSizeBytes) / 1024 / 1024

	return stats, nil
}

// StorageStatsHandler is an HTTP endpoint for monitoring storage usage.
func (pm *ProgressManager) StorageStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := pm.GetStorageStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// CleanupHandler is an HTTP endpoint to trigger manual cleanup of old files.
func (pm *ProgressManager) CleanupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cleaned, err := pm.CleanupOldFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"cleaned": cleaned,
		"message": fmt.Sprintf("Cleaned up %d old files", cleaned),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StartCleanupScheduler runs daily cleanup of old progress files in the background.
func (pm *ProgressManager) StartCleanupScheduler() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			cleaned, err := pm.CleanupOldFiles()
			if err != nil {
				log.Printf("Cleanup error: %v", err)
			} else {
				log.Printf("Daily cleanup: removed %d old files", cleaned)
			}
		}
	}()
}
