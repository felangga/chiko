package storage

import (
	"os"
)

func (s *Storage) SaveToFile(path string, payload []byte) error {
	return os.WriteFile(path, payload, 0644)
}
