package history

import (
	"encoding/json"
	"os"
	"time"

	"github.com/felangga/chiko/internal/entity"
	"github.com/felangga/chiko/internal/utils"
)

func (h *History) SaveHistory() error {
	h.Mu.RLock()
	entries := *h.Entries
	h.Mu.RUnlock()
	return h.writeHistory(entries)
}

func (h *History) writeHistory(entries []entity.HistoryEntry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	encryptedData, err := utils.Encrypt(data)
	if err != nil {
		return err
	}

	return os.WriteFile(h.Path, encryptedData, entity.FILE_PERMISSION)
}

func (h *History) AddEntry(session entity.Session) error {
	h.Mu.Lock()
	entry := entity.HistoryEntry{
		Session:   session,
		InvokedAt: time.Now(),
	}
	*h.Entries = append([]entity.HistoryEntry{entry}, *h.Entries...)

	if len(*h.Entries) > entity.HISTORY_MAX_SIZE {
		*h.Entries = (*h.Entries)[:entity.HISTORY_MAX_SIZE]
	}

	entries := *h.Entries
	h.Mu.Unlock()

	return h.writeHistory(entries)
}
