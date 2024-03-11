package entity

// Bookmarks configuration file
const BOOKMARKS_FILE_NAME = ".bookmarks"

type Bookmark struct {
	CategoryName string    `json:"category_name"`
	Sessions     []Session `json:"sessions"`
}
