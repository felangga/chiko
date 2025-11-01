package bookmark

import (
	"os"
	"path/filepath"

	"github.com/felangga/chiko/internal/entity"
)

// MigrateBookmark is used to move the old bookmark file location (same as binary folder) to the OS default config folder
func (b *Bookmark) MigrateBookmark() error {
	// Check if the old bookmark file exists
	// If found then we need to move the bookmark file to new location
	if _, err := os.Stat(entity.BOOKMARKS_FILE_NAME); err != nil {
		return nil
	}

	// Ensure the new directory exists
	newDir := filepath.Dir(b.Path)
	if err := os.MkdirAll(newDir, entity.DIR_PERMISSION); err != nil {
		return err
	}

	// Move the file
	if err := os.Rename(entity.BOOKMARKS_FILE_NAME, b.Path); err != nil {
		return err
	}

	return nil
}
