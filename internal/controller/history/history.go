package history

import (
	"path/filepath"
	"sync"

	"github.com/felangga/chiko/internal/entity"
	"github.com/felangga/chiko/internal/utils"
)

type History struct {
	Entries *[]entity.HistoryEntry
	Path    string
	Mu      sync.RWMutex
}

// NewHistory creates a new History controller
func NewHistory() History {
	entries := []entity.HistoryEntry{}
	return History{
		Entries: &entries,
		Path:    filepath.Join(utils.GetOSConfigDir(), "Chiko", entity.HISTORY_FILE_NAME),
	}
}
