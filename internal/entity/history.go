package entity

import "time"

const (
	HISTORY_FILE_NAME = ".history"
	HISTORY_MAX_SIZE  = 50
)

// HistoryEntry represents a single recorded RPC invocation
type HistoryEntry struct {
	Session   Session   `json:"session"`
	InvokedAt time.Time `json:"invoked_at"`
}
