package entity

// Bookmarks configuration file
const BOOKMARKS_FILE_NAME = ".bookmarks"

type Category struct {
	Name     string    `json:"category_name"`
	Sessions []Session `json:"sessions"`
}
