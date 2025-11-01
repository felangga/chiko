package bookmark

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/felangga/chiko/internal/entity"
)

func TestMigrateBookmark(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	oldFilePath := entity.BOOKMARKS_FILE_NAME
	newFilePath := filepath.Join(tempDir, ".bookmarks")

	// Create an old bookmark file
	if err := os.WriteFile(oldFilePath, []byte("old bookmarks data"), entity.FILE_PERMISSION); err != nil {
		t.Fatalf("failed to create old bookmark file: %v", err)
	}

	bookmark := &Bookmark{Path: newFilePath, Categories: &[]entity.Category{}}

	// Test case 1: Old bookmark file does not exist
	os.Remove(oldFilePath)
	err := bookmark.MigrateBookmark()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test case 2: Old bookmark file exists, should move it
	if err := os.WriteFile(oldFilePath, []byte("old bookmarks data"), entity.FILE_PERMISSION); err != nil {
		t.Fatalf("failed to create old bookmark file: %v", err)
	}
	err = bookmark.MigrateBookmark()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
		t.Errorf("expected new bookmark file to exist, got none")
	}

	// Test case 3: Error in moving the file
	// Simulate an error by removing the old file
	if err := os.WriteFile(oldFilePath, []byte("old bookmarks data"), entity.FILE_PERMISSION); err != nil {
		t.Fatalf("failed to create old bookmark file: %v", err)
	}
	err = bookmark.MigrateBookmark()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(oldFilePath); err == nil {
		t.Errorf("expected the file should be moved, but still exists")
	}

}
