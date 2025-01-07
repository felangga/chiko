package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/felangga/chiko/pkg/entity"
)

// PrintOutput used to print output to the output panel
func (u *UI) PrintOutput(param entity.Output) {
	var (
		metadata  string
		newBuffer string
	)

	out := u.Layout.OutputPanel
	_, _, width, _ := out.TextArea.GetRect()

	if param.WithHeader {
		if u.GRPC != nil && len(u.GRPC.Conn.ParseMetadata()) > 0 {
			for _, meta := range u.GRPC.Conn.ParseMetadata() {
				metadata += "- " + meta + "\n"
			}

			metaHeader := strings.Repeat(string(tcell.RuneCkBoard), 2) + "[ Request Metadata ]" + (strings.Repeat(string(tcell.RuneCkBoard), width-22)) + "\n\n"
			newBuffer = metaHeader + metadata + "\n"
		}

		payloadHeader := strings.Repeat(string(tcell.RuneCkBoard), 2) + "[ Request Payload ]" + (strings.Repeat(string(tcell.RuneCkBoard), width-21)) + "\n\n"
		newBuffer += payloadHeader + u.GRPC.Conn.RequestPayload

		responseHeader := "\n\n" + strings.Repeat(string(tcell.RuneCkBoard), 2) + "[ Response ]" + (strings.Repeat(string(tcell.RuneCkBoard), width-14)) + "\n\n"
		newBuffer += responseHeader + param.Content
	}

	out.TextArea.SetText(newBuffer, true)
}
