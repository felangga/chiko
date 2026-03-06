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
	defer h.Mu.RUnlock()
	return h.saveHistory()
}

func (h *History) saveHistory() error {
	data, err := json.Marshal(h.Entries)
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
	defer h.Mu.Unlock()

	entry := entity.HistoryEntry{
		Session:   session,
		InvokedAt: time.Now(),
	}

	*h.Entries = append([]entity.HistoryEntry{entry}, *h.Entries...)

	if len(*h.Entries) > entity.HISTORY_MAX_SIZE {
		*h.Entries = (*h.Entries)[:entity.HISTORY_MAX_SIZE]
	}

	return h.saveHistory()
}
