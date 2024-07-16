package ui

import (
	"fmt"

	"github.com/felangga/chiko/pkg/entity"
)

func (u *UI) InvokeRPC() {

	// Construct metadata info
	var metadata string
	for _, meta := range u.GRPC.Conn.ParseMetadata() {
		metadata += "- " + meta + "\n"
	}

	u.PrintLog(entity.Log{
		Content: "\nRequest Metadata:\n" + metadata + "\nRequest Payload:\n[yellow]" + u.GRPC.Conn.RequestPayload + "\n",
		Type:    entity.LOG_INFO,
	})

	err := u.GRPC.InvokeRPC()
	if err != nil {
		u.PrintLog(entity.Log{
			Content: fmt.Sprintf("❌ failed to invoke RPC, err: %v", err),
			Type:    entity.LOG_ERROR,
		})
		return
	}
}
