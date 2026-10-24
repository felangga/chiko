package events

import "github.com/felangga/chiko/internal/entity"

const (
	EventTypeVimModeChange EventType = "VimModeChange"
)

type VimModeChanged struct {
	NewMode entity.VimMode
}

// Type implements [IEvent].
func (e *VimModeChanged) Type() EventType {
	return EventTypeVimModeChange
}

var _ IEvent = (*VimModeChanged)(nil)
