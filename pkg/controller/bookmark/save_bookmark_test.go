package bookmark

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/felangga/chiko/pkg/entity"
	"github.com/google/uuid"
)

func TestSaveBookmark(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, ".bookmarks")

	bookmark := &Bookmark{
		Path: filePath,
		Categories: &[]entity.Category{
			{Name: "Category 1", Sessions: []entity.Session{{ID: uuid.New(), Name: "Session 1"}}},
			{Name: "Category 2", Sessions: []entity.Session{{ID: uuid.New(), Name: "Session 2"}}},
		},
	}

	// Test case 1: Save successfully
	err := bookmark.SaveBookmark()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify that the file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("expected bookmark file to exist, got none")
	}

	// Test case 2: Save with an invalid path
	bookmark.Path = "/invalid/path/bookmarks.json"
	err = bookmark.SaveBookmark()
	if err == nil {
		t.Errorf("expected an error, got none")
	}
}
