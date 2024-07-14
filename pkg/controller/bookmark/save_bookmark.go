package bookmark

import (
	"encoding/json"
	"os"

	"chiko/pkg/entity"
)

// SaveBookmark is used to save the bookmark object to file by encoding the object with JSON.
func (b Bookmark) SaveBookmark() error {
	// Encoding the object to JSON
	convert, err := json.Marshal(b.Categories)
	if err != nil {
		return err
	}

	// Saving the json to file
	err = os.WriteFile(entity.BOOKMARKS_FILE_NAME, convert, 0644)
	if err != nil {
		return err
	}

	return nil
}
