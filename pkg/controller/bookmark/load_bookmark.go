package bookmark

import (
	"encoding/json"
	"os"

	"chiko/pkg/entity"
)

// LoadBookmarks is used to load bookmarks from bookmark file
func (b *Bookmark) LoadBookmarks() error {
	if _, err := os.Stat(entity.BOOKMARKS_FILE_NAME); err != nil {
		return err
	}

	// Read bookmark file and dump to array of bookmarks
	file, err := os.ReadFile(entity.BOOKMARKS_FILE_NAME)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(file), &b.Bookmarks)
	if err != nil {
		return err
	}

	return nil
}