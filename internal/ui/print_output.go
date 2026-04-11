package ui

import (
	"github.com/google/uuid"

	"github.com/felangga/chiko/internal/entity"
)

// PrintOutput used to print output to the output panel
func (u *UI) PrintOutput(param entity.Output) {
	// Determine which SessionWindow output should go to
	var targetTab *SessionWindow
	if param.SessionID == uuid.Nil {
		if len(u.Sessions) > 0 {
			targetTab = u.ActiveSession
		}
	} else {
		for _, tab := range u.Sessions {
			if tab.ID == param.SessionID {
				targetTab = tab
				break
			}
		}
	}

	// Fallback
	if targetTab == nil {
		if len(u.Sessions) == 0 {
			return // Nothing to print to
		}
		targetTab = u.ActiveSession
	}

	go u.App.QueueUpdateDraw(func() {
		payload := param.ResponsePayload
		if payload == "" && param.Content != "" {
			payload = param.Content
		}

		// Update payload tab
		targetTab.OutputPanel.TextArea.SetText(payload, param.CursorAtEnd)
		// Update headers tab
		targetTab.OutputPanel.HeaderArea.SetText(param.ResponseHeaders, param.CursorAtEnd)
	})
}
