package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	
	"github.com/google/uuid"
	"github.com/slavakurilyak/ctx/internal/models"
)

const (
	historyDirName = ".ctx"
	historyDirPerm = 0755
	historyFilePerm = 0644
)

type HistoryManager struct {
	enabled bool
	baseDir string
}

func NewHistoryManager() *HistoryManager {
	// Check if history is disabled
	if os.Getenv("CTX_NO_HISTORY") == "true" {
		return &HistoryManager{enabled: false}
	}
	
	// Determine base directory for history
	baseDir := getHistoryBaseDir()
	
	return &HistoryManager{
		enabled: true,
		baseDir: baseDir,
	}
}

func (h *HistoryManager) SaveRecord(output *models.Output) error {
	if !h.enabled || output == nil {
		return nil
	}
	
	// Ensure history directory exists
	historyDir := filepath.Join(h.baseDir, historyDirName)
	if err := os.MkdirAll(historyDir, historyDirPerm); err != nil {
		// Silently fail - history is not critical
		return nil
	}
	
	// Generate filename with timestamp and UUID
	filename := generateHistoryFilename()
	filePath := filepath.Join(historyDir, filename)
	
	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil
	}
	
	// Write to file
	if err := os.WriteFile(filePath, data, historyFilePerm); err != nil {
		// Silently fail - history is not critical
		return nil
	}
	
	return nil
}

func getHistoryBaseDir() string {
	// Priority order:
	// 1. CTX_HISTORY_DIR environment variable
	// 2. Current working directory
	// 3. Home directory as fallback
	
	if dir := os.Getenv("CTX_HISTORY_DIR"); dir != "" {
		return dir
	}
	
	// Try current working directory
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	
	// Fallback to home directory
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	
	// Last resort: temp directory
	return os.TempDir()
}

func generateHistoryFilename() string {
	// Format: YYYY-MM-DD_HH-MM-SS_<uuid>.json
	// Using underscores and hyphens for better readability and sorting
	now := time.Now()
	id := uuid.New()
	
	return fmt.Sprintf(
		"%04d-%02d-%02d_%02d-%02d-%02d_%s.json",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		id.String(),
	)
}

// HistoryRecord is no longer needed as we save Output directly

// createHistoryRecord is no longer needed as Output already contains all fields

// GetHistoryDir returns the path to the history directory
func (h *HistoryManager) GetHistoryDir() string {
	if !h.enabled {
		return ""
	}
	return filepath.Join(h.baseDir, historyDirName)
}

// IsEnabled returns whether history recording is enabled
func (h *HistoryManager) IsEnabled() bool {
	return h.enabled
}

// CleanOldRecords removes history files older than the specified duration
func (h *HistoryManager) CleanOldRecords(maxAge time.Duration) error {
	if !h.enabled {
		return nil
	}
	
	historyDir := filepath.Join(h.baseDir, historyDirName)
	
	// Read directory
	entries, err := os.ReadDir(historyDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	
	cutoff := time.Now().Add(-maxAge)
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		// Check if it's a JSON file
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		
		// Get file info
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		// Remove if older than cutoff
		if info.ModTime().Before(cutoff) {
			filePath := filepath.Join(historyDir, entry.Name())
			os.Remove(filePath) // Ignore errors
		}
	}
	
	return nil
}