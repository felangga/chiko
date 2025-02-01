package storage

import (
	"os"
)

func (s *Storage) GetWorkingDirectory() (string, error) {
	// Show dialog box so user can input the file path
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dir, nil
}
