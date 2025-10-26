package entity

// Bookmarks configuration file
const (
	DIR_PERMISSION      = 0755
	FILE_PERMISSION     = 0644
	BOOKMARKS_FILE_NAME = ".bookmarks"
)

type Category struct {
	Name     string    `json:"category_name"`
	Sessions []Session `json:"sessions"`
}
