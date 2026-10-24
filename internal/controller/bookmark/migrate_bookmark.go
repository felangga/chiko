package bookmark

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"

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

	// Move the file; fall back to copy+remove across devices (e.g. ext4 -> tmpfs)
	if err := os.Rename(entity.BOOKMARKS_FILE_NAME, b.Path); err != nil {
		if isCrossDevice(err) {
			return moveFileCrossDevice(entity.BOOKMARKS_FILE_NAME, b.Path)
		}
		return err
	}

	return nil
}

// isCrossDevice reports whether err is an EXDEV / cross-device rename error.
func isCrossDevice(err error) bool {
	return err != nil && errors.Is(err, syscall.EXDEV)
}

// moveFileCrossDevice copies src to dst and removes src.
func moveFileCrossDevice(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dst, data, entity.FILE_PERMISSION); err != nil {
		return err
	}
	return os.Remove(src)
}
