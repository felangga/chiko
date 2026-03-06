package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/felangga/chiko/internal/entity"
)

// LoadHistory reads history entries from the history file
func (h *History) LoadHistory() error {
	if err := h.prepareHistory(); err != nil {
		return err
	}

	file, err := os.ReadFile(h.Path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, h.Entries); err != nil {
		return fmt.Errorf("failed to read history file, maybe corrupted: %v", err)
	}

	return nil
}

func (h *History) prepareHistory() error {
	if _, err := os.Stat(h.Path); err == nil {
		return nil // file already exists
	}

	if err := os.MkdirAll(filepath.Dir(h.Path), entity.DIR_PERMISSION); err != nil {
		return err
	}

	return h.SaveHistory()
}
