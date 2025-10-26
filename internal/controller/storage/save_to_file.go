package storage

import (
	"os"

	"github.com/felangga/chiko/internal/entity"
)

func (s *Storage) SaveToFile(path string, payload []byte) error {
	return os.WriteFile(path, payload, entity.FILE_PERMISSION)
}
