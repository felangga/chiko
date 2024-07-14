package bookmark

import "github.com/felangga/chiko/pkg/entity"

type Bookmark struct {
	Categories *[]entity.Category
}

// NewBookmark is used to create a new bookmark object
func NewBookmark() Bookmark {
	category := []entity.Category{}
	b := Bookmark{
		&category,
	}

	return b
}
