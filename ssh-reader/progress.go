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
	"os"
	"path/filepath"
	"time"
)

type UserProgress struct {
	Username       string            `json:"username"`
	CurrentChapter int               `json:"current_chapter"`
	ScrollOffset   int               `json:"scroll_offset"`
	LastRead       time.Time         `json:"last_read"`
	ChapterProgress map[int]bool     `json:"chapter_progress"` // tracks completed chapters
	Bookmarks      []Bookmark        `json:"bookmarks"`
	ReadingTime    time.Duration     `json:"reading_time"`
	SessionStart   time.Time         `json:"-"` // not persisted
}

type Bookmark struct {
	Chapter      int       `json:"chapter"`
	ScrollOffset int       `json:"scroll_offset"`
	Note         string    `json:"note"`
	Created      time.Time `json:"created"`
}

type ProgressManager struct {
	dataDir string
}

func NewProgressManager() *ProgressManager {
	dataDir := "../.void_reader_data"
	os.MkdirAll(dataDir, 0755)
	return &ProgressManager{dataDir: dataDir}
}

func (pm *ProgressManager) LoadProgress(username string) (*UserProgress, error) {
	if username == "" {
		username = "anonymous"
	}

	filename := filepath.Join(pm.dataDir, username+".json")
	
	// If file doesn't exist, return new progress
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &UserProgress{
			Username:        username,
			CurrentChapter:  0,
			ScrollOffset:    0,
			LastRead:        time.Now(),
			ChapterProgress: make(map[int]bool),
			Bookmarks:       []Bookmark{},
			ReadingTime:     0,
			SessionStart:    time.Now(),
		}, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var progress UserProgress
	err = json.Unmarshal(data, &progress)
	if err != nil {
		return nil, err
	}

	progress.SessionStart = time.Now()
	return &progress, nil
}

func (pm *ProgressManager) SaveProgress(progress *UserProgress) error {
	if progress.Username == "" {
		progress.Username = "anonymous"
	}

	// Update reading time
	if !progress.SessionStart.IsZero() {
		progress.ReadingTime += time.Since(progress.SessionStart)
		progress.SessionStart = time.Now()
	}

	progress.LastRead = time.Now()

	filename := filepath.Join(pm.dataDir, progress.Username+".json")
	
	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (p *UserProgress) AddBookmark(chapter, scrollOffset int, note string) {
	bookmark := Bookmark{
		Chapter:      chapter,
		ScrollOffset: scrollOffset,
		Note:         note,
		Created:      time.Now(),
	}
	
	p.Bookmarks = append(p.Bookmarks, bookmark)
}

func (p *UserProgress) RemoveBookmark(index int) {
	if index >= 0 && index < len(p.Bookmarks) {
		p.Bookmarks = append(p.Bookmarks[:index], p.Bookmarks[index+1:]...)
	}
}

func (p *UserProgress) MarkChapterComplete(chapter int) {
	if p.ChapterProgress == nil {
		p.ChapterProgress = make(map[int]bool)
	}
	p.ChapterProgress[chapter] = true
}

func (p *UserProgress) IsChapterComplete(chapter int) bool {
	if p.ChapterProgress == nil {
		return false
	}
	return p.ChapterProgress[chapter]
}

func (p *UserProgress) GetCompletionPercentage(totalChapters int) float64 {
	if totalChapters == 0 {
		return 0
	}
	
	completed := 0
	for i := 0; i < totalChapters; i++ {
		if p.IsChapterComplete(i) {
			completed++
		}
	}
	
	return float64(completed) / float64(totalChapters) * 100
}

func (p *UserProgress) GetReadingStats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_reading_time": p.ReadingTime,
		"last_read":         p.LastRead,
		"current_chapter":   p.CurrentChapter + 1, // Display as 1-based
		"bookmarks_count":   len(p.Bookmarks),
		"chapters_completed": len(p.ChapterProgress),
	}
	
	return stats
}