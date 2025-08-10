# Storage Capacity Analysis for Void Chronicles

## Current Storage Situation

**PROBLEM**: User progress files are stored in the container's ephemeral storage, which:
1. Gets wiped on every deployment
2. Has no configured persistent volume
3. Could fill up with many users

## Storage Calculations

### Per-User Storage
Based on the UserProgress structure:
```json
{
  "username": "user123",                    // ~20 bytes
  "current_chapter": 15,                    // ~20 bytes
  "scroll_offset": 1250,                    // ~20 bytes
  "last_read": "2024-01-15T10:30:00Z",     // ~30 bytes
  "chapter_progress": {                     // ~200 bytes (20 chapters)
    "1": true, "2": true, ...
  },
  "bookmarks": [                           // ~150 bytes per bookmark
    {
      "chapter": 5,
      "scroll_offset": 500,
      "note": "Important scene",
      "created": "2024-01-15T10:30:00Z"
    }
  ],
  "reading_time": 3600000000000            // ~20 bytes
}
```

**Average user file size**: ~500 bytes - 2KB (with 5-10 bookmarks)

## Fly.io Free Tier Limits

- **Ephemeral storage**: ~8GB in container (but lost on redeploy!)
- **Persistent volumes**: 3GB free (across all apps)
- **RAM**: 256MB (current config)

## Storage Capacity

### Without Persistent Volume (Current - PROBLEM!)
- Storage: Ephemeral, lost on every deployment
- Users lose all progress when app redeploys!

### With 1GB Persistent Volume (Recommended)
- **1GB = 1,073,741,824 bytes**
- Average file size: 2KB per user
- **Maximum users**: ~500,000 users
- With heavy usage (10KB per user): ~100,000 users

## Solutions

### Option 1: Add Persistent Volume (Recommended)
```toml
# Add to fly.toml
[[mounts]]
  source = "void_data"
  destination = "/data"
  initial_size = "1gb"
```

Then update progress.go:
```go
func NewProgressManager() *ProgressManager {
    // Use persistent volume path
    dataDir := "/data/void_reader_data"
    if _, err := os.Stat("/data"); os.IsNotExist(err) {
        // Fallback for local development
        dataDir = "../.void_reader_data"
    }
    os.MkdirAll(dataDir, 0755)
    return &ProgressManager{dataDir: dataDir}
}
```

### Option 2: Client-Side Storage (Alternative)
Store progress in SSH client using escape sequences:
```go
// Store in client terminal title or environment
fmt.Printf("\033]0;PROGRESS:%s\007", encodedProgress)
```

### Option 3: Cleanup Policy
Implement automatic cleanup:
```go
func (pm *ProgressManager) CleanupOldFiles() error {
    files, _ := os.ReadDir(pm.dataDir)
    for _, file := range files {
        info, _ := file.Info()
        // Delete files not accessed in 90 days
        if time.Since(info.ModTime()) > 90*24*time.Hour {
            os.Remove(filepath.Join(pm.dataDir, file.Name()))
        }
    }
    return nil
}
```

### Option 4: Storage Limits
Add per-user limits:
```go
const (
    MaxBookmarksPerUser = 20
    MaxNoteLength = 200
    MaxFileSize = 10 * 1024 // 10KB
)

func (pm *ProgressManager) SaveProgress(progress *UserProgress) error {
    // Trim bookmarks if too many
    if len(progress.Bookmarks) > MaxBookmarksPerUser {
        progress.Bookmarks = progress.Bookmarks[:MaxBookmarksPerUser]
    }
    
    // Check file size before saving
    data, _ := json.Marshal(progress)
    if len(data) > MaxFileSize {
        return fmt.Errorf("user data too large")
    }
    
    // Save...
}
```

## Monitoring Storage Usage

Add a monitoring endpoint:
```go
func getStorageStats(w http.ResponseWriter, r *http.Request) {
    dataDir := "/data/void_reader_data"
    
    var totalSize int64
    var fileCount int
    
    filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            totalSize += info.Size()
            fileCount++
        }
        return nil
    })
    
    stats := map[string]interface{}{
        "total_size_bytes": totalSize,
        "total_size_mb": float64(totalSize) / 1024 / 1024,
        "file_count": fileCount,
        "avg_file_size": totalSize / int64(fileCount),
        "estimated_capacity": 1073741824 / (totalSize / int64(fileCount)),
    }
    
    json.NewEncoder(w).Encode(stats)
}
```

## Recommended Implementation

1. **Add 1GB persistent volume** (immediate fix)
2. **Implement 90-day cleanup** (prevent accumulation)
3. **Add storage monitoring** (track usage)
4. **Set bookmark limits** (prevent abuse)

This gives you capacity for 100,000+ active users with room to grow!