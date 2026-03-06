package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/felangga/chiko/internal/entity"
	"github.com/felangga/chiko/internal/utils"
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

	decryptedData, err := utils.Decrypt(file)
	if err != nil {
		// If decryption fails, it could be an old unencrypted file or corrupt.
		// For safety, we try unmarshaling the original file in case it was not encrypted.
		if errJson := json.Unmarshal(file, h.Entries); errJson == nil {
			return nil
		}
		return fmt.Errorf("failed to decrypt history file: %v", err)
	}

	if err := json.Unmarshal(decryptedData, h.Entries); err != nil {
		return fmt.Errorf("failed to read history file, maybe corrupted: %v", err)
	}

	return nil
}

func (h *History) prepareHistory() error {
	if _, err := os.Stat(h.Path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(h.Path), entity.DIR_PERMISSION); err != nil {
		return err
	}

	return h.SaveHistory()
}
