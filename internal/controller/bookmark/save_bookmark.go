package bookmark

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SaveBookmark is used to save the bookmark object to file by encoding the object with JSON.
func (b Bookmark) SaveBookmark() error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(b.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Encoding the object to JSON
	convert, err := json.Marshal(b.Categories)
	if err != nil {
		return err
	}

	// Saving the json to file
	return os.WriteFile(b.Path, convert, 0644)
}
