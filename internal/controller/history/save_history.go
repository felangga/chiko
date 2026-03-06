package history

import (
	"encoding/json"
	"os"
	"time"

	"github.com/felangga/chiko/internal/entity"
)

// SaveHistory persists the current history entries to disk
func (h *History) SaveHistory() error {
	data, err := json.Marshal(h.Entries)
	if err != nil {
		return err
	}
	return os.WriteFile(h.Path, data, entity.FILE_PERMISSION)
}

// AddEntry prepends a new history entry and trims the list to HISTORY_MAX_SIZE.
// It then persists the updated list to disk.
func (h *History) AddEntry(session entity.Session) error {
	entry := entity.HistoryEntry{
		Session:   session,
		InvokedAt: time.Now(),
	}

	// Prepend so the most recent is always first
	*h.Entries = append([]entity.HistoryEntry{entry}, *h.Entries...)

	// Cap the history size
	if len(*h.Entries) > entity.HISTORY_MAX_SIZE {
		*h.Entries = (*h.Entries)[:entity.HISTORY_MAX_SIZE]
	}

	return h.SaveHistory()
}
