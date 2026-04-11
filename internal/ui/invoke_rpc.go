package ui

import (
	"fmt"

	"github.com/felangga/chiko/internal/entity"
)

func (u *UI) InvokeRPC() {
	sw := u.ActiveSession

	// Check if any method is selected
	if sw.GRPC.Conn.SelectedMethod == nil {
		u.PrintOutput(entity.Output{
			SessionID:  sw.ID,
			Content:    "❗ no method selected. Please connect and select a method first.",
			WithHeader: true,
		})
		return
	}

	// Check if there is no active connection
	if sw.GRPC.Conn.ActiveConnection == nil {
		u.PrintOutput(entity.Output{
			SessionID:  sw.ID,
			Content:    "❗ no active connection. Please enter a URL and press Enter to connect.",
			WithHeader: true,
		})
		return
	}

	if sw.Loading {
		return // Already loading
	}

	u.SetLoading(sw, true)

	go func() {
		defer u.SetLoading(sw, false)

		// Invoke the RPC
		err := sw.GRPC.InvokeRPC()
		if err != nil {
			u.PrintOutput(entity.Output{
				SessionID:  sw.GRPC.Conn.ID,
				Content:    err.Error(),
				WithHeader: true,
			})
			return
		}

		// Record successful invocation in history
		if err := u.History.AddEntry(*sw.GRPC.Conn); err != nil {
			u.PrintLog(entity.Log{
				Content: fmt.Sprintf("⚠️ could not save history entry: %v", err),
				Type:    entity.LOG_ERROR,
			})
		}

		u.RefreshHistoryPanel()
	}()
}
