package ui

import (
	"fmt"

	"github.com/felangga/chiko/internal/entity"
)

func (u *UI) InvokeRPC() {
	// Check if any method is selected
	if u.ActiveSession.GRPC.Conn.SelectedMethod == nil {
		u.PrintOutput(entity.Output{
			SessionID:  u.ActiveSession.ID,
			Content:    "❗ no method selected. Please connect and select a method first.",
			WithHeader: true,
		})
		return
	}

	// Check if there is no active connection
	if u.ActiveSession.GRPC.Conn.ActiveConnection == nil {
		u.PrintOutput(entity.Output{
			SessionID:  u.ActiveSession.ID,
			Content:    "❗ no active connection. Please enter a URL and press Enter to connect.",
			WithHeader: true,
		})
		return
	}

	// Invoke the RPC
	err := u.GRPC.InvokeRPC()
	if err != nil {
		u.PrintOutput(entity.Output{
			SessionID:  u.GRPC.Conn.ID,
			Content:    err.Error(),
			WithHeader: true,
		})
		return
	}

	// Record successful invocation in history
	if err := u.History.AddEntry(*u.GRPC.Conn); err != nil {
		u.PrintLog(entity.Log{
			Content: fmt.Sprintf("⚠️ could not save history entry: %v", err),
			Type:    entity.LOG_ERROR,
		})
	}

	u.RefreshHistoryPanel()
}
