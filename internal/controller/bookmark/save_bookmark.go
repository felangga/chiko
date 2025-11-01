package bookmark

import (
	"encoding/json"
	"os"

	"github.com/felangga/chiko/internal/entity"
)

// SaveBookmark is used to save the bookmark object to file by encoding the object with JSON.
func (b *Bookmark) SaveBookmark() error {
	// Encoding the object to JSON
	convert, err := json.Marshal(b.Categories)
	if err != nil {
		return err
	}

	// Saving the json to file
	return os.WriteFile(b.Path, convert, entity.FILE_PERMISSION)
}
