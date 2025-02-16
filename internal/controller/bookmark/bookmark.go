package bookmark

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/felangga/chiko/internal/entity"
)

type Bookmark struct {
	Categories *[]entity.Category
	Path       string
}

// NewBookmark is used to create a new bookmark object
func NewBookmark() Bookmark {
	category := []entity.Category{}

	b := Bookmark{
		&category,
		filepath.Join(GetOSConfigDir(), "Chiko", entity.BOOKMARKS_FILE_NAME),
	}

	return b
}

// GetOSConfigDir is used to get the path of the config directory based on the operating system.
func GetOSConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("APPDATA")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "linux":
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			return xdgConfigHome
		}
		return filepath.Join(os.Getenv("HOME"), ".config")
	default:
		return ""
	}
}
