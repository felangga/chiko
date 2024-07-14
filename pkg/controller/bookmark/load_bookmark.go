package bookmark

import (
	"encoding/json"
	"fmt"
	"os"

	"chiko/pkg/entity"
)

// LoadBookmarks is used to load bookmarks from bookmark file
func (b *Bookmark) LoadBookmarks() error {
	if _, err := os.Stat(entity.BOOKMARKS_FILE_NAME); err != nil {
		// If no bookmarks file found, then create a new one
		b.SaveBookmark()
	}

	// Read bookmark file and dump to array of bookmarks
	file, err := os.ReadFile(entity.BOOKMARKS_FILE_NAME)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(file), &b.Categories)
	if err != nil {
		return fmt.Errorf("failed to read file, maybe corrupted")
	}

	return nil
}
