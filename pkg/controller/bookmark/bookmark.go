package bookmark

import (
	"os"

	"github.com/felangga/chiko/pkg/entity"
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
		GetOSConfigDir() + string(os.PathSeparator) + "chiko" + string(os.PathSeparator) + entity.BOOKMARKS_FILE_NAME,
	}

	return b
}
