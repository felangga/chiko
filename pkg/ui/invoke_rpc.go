package ui

import (
	"chiko/pkg/entity"
	"fmt"
)

func (u *UI) InvokeRPC() {
	err := u.Controller.InvokeRPC()
	if err != nil {
		u.PrintLog(entity.LogParam{
			Content: fmt.Sprintf("❌ failed to invoke RPC, err: %v", err),
			Type:    entity.LOG_ERROR,
		})
		return
	}
}
