package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/felangga/chiko/internal/entity"
	"github.com/felangga/chiko/internal/utils"
)

type Workspace struct {
	ActiveSessions *[]*entity.Session
	Path           string
	Mu             sync.RWMutex
}

func NewWorkspace() Workspace {
	sessions := []*entity.Session{}
	return Workspace{
		ActiveSessions: &sessions,
		Path:           filepath.Join(utils.GetOSConfigDir(), "Chiko", entity.WORKSPACE_FILE_NAME),
	}
}

func (w *Workspace) SaveWorkspace(ws entity.Workspace) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	data, err := json.Marshal(ws)
	if err != nil {
		return err
	}

	encryptedData, err := utils.Encrypt(data)
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(w.Path), entity.DIR_PERMISSION); err != nil {
		return err
	}

	return os.WriteFile(w.Path, encryptedData, entity.FILE_PERMISSION)
}

func (w *Workspace) LoadWorkspace() (*entity.Workspace, error) {
	w.Mu.RLock()
	defer w.Mu.RUnlock()

	if _, err := os.Stat(w.Path); os.IsNotExist(err) {
		return nil, nil // No workspace yet
	}

	encryptedData, err := os.ReadFile(w.Path)
	if err != nil {
		return nil, err
	}

	data, err := utils.Decrypt(encryptedData)
	if err != nil {
		return nil, err
	}

	var ws entity.Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, err
	}

	return &ws, nil
}
