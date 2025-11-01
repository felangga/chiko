package bookmark

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/felangga/chiko/internal/entity"
)

func TestLoadBookmarks(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "bookmarks.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	bookmark := &Bookmark{Path: tempFile.Name(), Categories: &[]entity.Category{}}

	// Test case 1: File does not exist, should create a new one
	err = os.Remove(tempFile.Name())
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove temp file: %v", err)
	}
	if err := bookmark.LoadBookmarks(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test case 2: Valid bookmarks file
	categories := []entity.Category{
		{Name: "Category 1", Sessions: []entity.Session{{ID: uuid.New(), Name: "Session 1"}}},
		{Name: "Category 2", Sessions: []entity.Session{{ID: uuid.New(), Name: "Session 2"}}},
	}
	bookmarkJSON, err := json.Marshal(categories)
	if err != nil {
		t.Fatalf("failed to marshal categories: %v", err)
	}
	err = os.WriteFile(tempFile.Name(), bookmarkJSON, entity.FILE_PERMISSION)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	err = bookmark.LoadBookmarks()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(*bookmark.Categories) != 2 {
		t.Errorf("expected 2 bookmarks, got %d", len(*bookmark.Categories))
	}

	// Test case 3: Corrupted bookmarks file
	err = os.WriteFile(tempFile.Name(), []byte("corrupted data"), entity.FILE_PERMISSION)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = bookmark.LoadBookmarks()
	if err == nil {
		t.Errorf("expected an error, got none")
	}
}
