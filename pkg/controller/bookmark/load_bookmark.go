package bookmark

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadBookmarks is used to load bookmarks from bookmark file
func (b *Bookmark) LoadBookmarks() error {

	if _, err := os.Stat(b.Path); err != nil {
		// If no bookmarks file found, then create a new one
		if err := b.SaveBookmark(); err != nil {
			return err
		}
	}

	// Read bookmark file and dump to array of bookmarks
	file, err := os.ReadFile(b.Path)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(file), &b.Categories)
	if err != nil {
		return fmt.Errorf("failed to read file, maybe corrupted, error: %v", err)
	}

	return nil
}
