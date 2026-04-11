package entity

import "github.com/google/uuid"

const (
	WORKSPACE_FILE_NAME = ".workspace"
)

type Workspace struct {
	ActiveSessions  []*Session `json:"active_sessions"`
	ActiveSessionID uuid.UUID   `json:"active_session_id"`
}
