package history

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/felangga/chiko/internal/entity"
)

func TestSaveAndLoadHistory(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "history.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	entries := []entity.HistoryEntry{}
	h := &History{
		Entries: &entries,
		Path:    tempFile.Name(),
	}

	// Test case 1: Add entries
	session1 := entity.Session{ID: uuid.New(), Name: "Session 1", ServerURL: "localhost:50051"}
	session2 := entity.Session{ID: uuid.New(), Name: "Session 2", ServerURL: "localhost:50052"}

	if err := h.AddEntry(session1); err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}
	if err := h.AddEntry(session2); err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}

	if len(*h.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(*h.Entries))
	}

	// Session 2 should be first (most recent)
	if (*h.Entries)[0].Session.ID != session2.ID {
		t.Errorf("expected recent entry to be session 2, got %v", (*h.Entries)[0].Session.ID)
	}

	// Test case 2: Load History from file
	h2_entries := []entity.HistoryEntry{}
	h2 := &History{
		Entries: &h2_entries,
		Path:    tempFile.Name(),
	}

	if err := h2.LoadHistory(); err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}

	if len(*h2.Entries) != 2 {
		t.Errorf("expected 2 loaded entries, got %d", len(*h2.Entries))
	}

	if (*h2.Entries)[0].Session.ServerURL != "localhost:50052" {
		t.Errorf("expected first entry server URL to be localhost:50052, got %s", (*h2.Entries)[0].Session.ServerURL)
	}

	// Test case 3: Max history size
	h3_entries := []entity.HistoryEntry{}
	h3 := &History{
		Entries: &h3_entries,
		Path:    tempFile.Name(),
	}

	// Add more than HISTORY_MAX_SIZE
	for i := 0; i < entity.HISTORY_MAX_SIZE+5; i++ {
		err := h3.AddEntry(entity.Session{ID: uuid.New(), Name: "Spam"})
		if err != nil {
			t.Fatalf("AddEntry loop failed at %d: %v", i, err)
		}
	}

	if len(*h3.Entries) != entity.HISTORY_MAX_SIZE {
		t.Errorf("expected max %d entries, got %d", entity.HISTORY_MAX_SIZE, len(*h3.Entries))
	}

	// Test case 4: Corrupted file
	err = os.WriteFile(tempFile.Name(), []byte("invalid json data"), 0644)
	if err != nil {
		t.Fatalf("failed to corrupt file: %v", err)
	}

	h4_entries := []entity.HistoryEntry{}
	h4 := &History{Entries: &h4_entries, Path: tempFile.Name()}
	if err := h4.LoadHistory(); err == nil {
		t.Error("expected error loading corrupted file, got nil")
	}
}

func TestHistoryEncryption(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "history_encrypted.json")
	defer os.Remove(tempFile.Name())

	entries := []entity.HistoryEntry{}
	h := &History{Entries: &entries, Path: tempFile.Name()}

	session := entity.Session{Name: "Secret Session"}
	if err := h.AddEntry(session); err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}

	// Verify the file on disk is NOT plain text JSON
	data, _ := os.ReadFile(tempFile.Name())
	var testMap map[string]interface{}
	if err := json.Unmarshal(data, &testMap); err == nil {
		t.Error("expected file to be encrypted (invalid JSON), but unmarshal succeeded")
	}

	// Verify we can still load it
	h2_entries := []entity.HistoryEntry{}
	h2 := &History{Entries: &h2_entries, Path: tempFile.Name()}
	if err := h2.LoadHistory(); err != nil {
		t.Fatalf("LoadHistory failed on encrypted file: %v", err)
	}

	if len(*h2.Entries) != 1 || (*h2.Entries)[0].Session.Name != "Secret Session" {
		t.Error("failed to load session from encrypted file correctly")
	}
}

func TestHistoryTimeRecording(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "history_time.json")
	defer os.Remove(tempFile.Name())

	entries := []entity.HistoryEntry{}
	h := &History{Entries: &entries, Path: tempFile.Name()}

	before := time.Now().Truncate(time.Second)
	_ = h.AddEntry(entity.Session{Name: "Timed"})
	after := time.Now()

	entryTime := (*h.Entries)[0].InvokedAt
	if entryTime.Before(before) || entryTime.After(after) {
		t.Errorf("entry time %v outside range [%v, %v]", entryTime, before, after)
	}
}
