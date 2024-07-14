package controller

import (
	"encoding/json"
	"os"

	"chiko/pkg/entity"
)

// SaveBookmark is used to save the bookmark object to file by encoding the object with JSON.
func (c Controller) SaveBookmark() error {
	// Encoding the object to JSON
	convert, err := json.Marshal(c.Bookmarks)
	if err != nil {
		return err
	}

	// Saving the json to file
	err = os.WriteFile(entity.BOOKMARKS_FILE_NAME, convert, 0644)
	if err != nil {
		// c.PrintLog("💾 failed to write bookmark configuration, maybe write-protected?", LOG_ERROR)
		return err
	}

	return nil
}
