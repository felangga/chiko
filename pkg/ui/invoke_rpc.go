package ui

import (
	"chiko/pkg/entity"
	"fmt"
)

func (u *UI) InvokeRPC() {
	err := u.GRPC.InvokeRPC()
	if err != nil {
		u.PrintLog(entity.Log{
			Content: fmt.Sprintf("❌ failed to invoke RPC, err: %v", err),
			Type:    entity.LOG_ERROR,
		})
		return
	}
}
