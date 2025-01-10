package storage

import (
	"testing"
)

func TestGetWorkingDirectory(t *testing.T) {
	storage := &Storage{}

	// Test case: Get current working directory successfully
	dir, err := storage.GetWorkingDirectory()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify that the directory is not empty
	if dir == "" {
		t.Error("expected non-empty directory, got empty")
	}
}
