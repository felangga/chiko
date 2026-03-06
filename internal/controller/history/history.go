package history

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/felangga/chiko/internal/entity"
)

type History struct {
	Entries *[]entity.HistoryEntry
	Path    string
}

// NewHistory creates a new History controller
func NewHistory() History {
	entries := []entity.HistoryEntry{}
	return History{
		Entries: &entries,
		Path:    filepath.Join(getOSConfigDir(), "Chiko", entity.HISTORY_FILE_NAME),
	}
}

func getOSConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("APPDATA")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "linux":
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return xdg
		}
		return filepath.Join(os.Getenv("HOME"), ".config")
	default:
		return ""
	}
}
