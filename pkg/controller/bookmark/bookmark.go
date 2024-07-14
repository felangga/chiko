package bookmark

import (
	"chiko/pkg/entity"
)

type Bookmark struct {
	Bookmarks []entity.Bookmark
}

// NewBookmark is used to create a new bookmark object
func NewBookmark() Bookmark {
	bookmarks := []entity.Bookmark{}
	b := Bookmark{
		bookmarks,
	}

	return b
}
