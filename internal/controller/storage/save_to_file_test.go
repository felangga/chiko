package storage

import (
	"os"
	"testing"
)

func TestSaveToFile(t *testing.T) {
	storage := &Storage{}

	// Test case 1: Save data successfully
	testPath := "testfile.txt"
	testData := []byte("This is a test.")

	err := storage.SaveToFile(testPath, testData)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify that the file was created and contains the correct data
	content, err := os.ReadFile(testPath)
	if err != nil {
		t.Errorf("expected to read file, got error: %v", err)
	}

	if string(content) != string(testData) {
		t.Errorf("expected file content to be %q, got %q", testData, content)
	}

	// Clean up the test file
	os.Remove(testPath)

	// Test case 2: Attempt to save to an invalid path
	invalidPath := "/invalid/path/testfile.txt"
	err = storage.SaveToFile(invalidPath, testData)
	if err == nil {
		t.Error("expected an error when saving to an invalid path, got none")
	}
}
