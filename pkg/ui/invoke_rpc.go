package ui

import (
	"github.com/felangga/chiko/pkg/entity"
)

func (u *UI) InvokeRPC() {
	// Check if any method is selected
	if u.GRPC.Conn.SelectedMethod == nil {
		u.PrintLog(entity.Log{
			Content: "❗ no method selected",
			Type:    entity.LOG_ERROR,
		})
		return
	}

	// Check if there is no active connection
	if u.GRPC.Conn.ActiveConnection == nil {
		u.PrintLog(entity.Log{
			Content: "❗ no active connection",
			Type:    entity.LOG_ERROR,
		})
		return
	}

	// Invoke the RPC
	err := u.GRPC.InvokeRPC()
	if err != nil {
		u.PrintOutput(entity.Output{
			Content:        err.Error(),
			ShowTimeHeader: false,
		})
		return
	}
}
