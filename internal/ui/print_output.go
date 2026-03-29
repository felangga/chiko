package ui

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gdamore/tcell/v2"

	"github.com/felangga/chiko/internal/entity"
)

// PrintOutput used to print output to the output panel
func (u *UI) PrintOutput(param entity.Output) {
	var (
		metadata  string
		newBuffer string
	)

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

	out := targetTab.OutputPanel
	targetGRPC := targetTab.GRPC
	_, _, width, _ := out.TextArea.GetRect()

	timeHeader := time.Now().Format("15:04:05 02/01/2006")

	if param.WithHeader {
		if targetGRPC != nil && len(targetGRPC.Conn.ParseMetadata()) > 0 {
			for _, meta := range targetGRPC.Conn.ParseMetadata() {
				metadata += "  ► " + meta + "\n"
			}

			metaHeader := strings.Repeat(string(tcell.RuneCkBoard), 2) + "[ Request Metadata ]" + (strings.Repeat(string(tcell.RuneCkBoard), width-47)) + "[ " + timeHeader + " ]" + strings.Repeat(string(tcell.RuneCkBoard), 2) + "\n\n"
			newBuffer = metaHeader + metadata + "\n"
		}



		responseHeader := "\n\n" + strings.Repeat(string(tcell.RuneCkBoard), 2) + "[ Response Payload ]" + (strings.Repeat(string(tcell.RuneCkBoard), width-47)) + "[ " + timeHeader + " ]" + strings.Repeat(string(tcell.RuneCkBoard), 2) + "\n"
		newBuffer += responseHeader + param.Content
	} else {
		newBuffer = param.Content
	}

	go u.App.QueueUpdateDraw(func() {
		out.TextArea.SetText(newBuffer, param.CursorAtEnd)
	})
}

