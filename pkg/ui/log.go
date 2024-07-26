package ui

import (
	"fmt"
	"time"

	"github.com/felangga/chiko/pkg/entity"
)

// PrintLog is used to print a log message to the log window
func (u *UI) PrintLog(param entity.Log) {
	// Get last log message
	lastLog := u.Layout.LogList.GetText(false)

	warnaLog := "white"
	switch param.Type {
	case entity.LOG_ERROR:
		warnaLog = "red"
	case entity.LOG_WARNING:
		warnaLog = "yellow"
	}

	formatLog := fmt.Sprintf("[green][%s] [%s]%s [white]\n", time.Now().Format(time.RFC822), warnaLog, param.Content)
	u.Layout.LogList.SetWordWrap(true).SetText(lastLog + formatLog)

	// Scroll log window to bottom
	u.Layout.LogList.ScrollToEnd()
}

// PrintOutput used to print output to the output panel
func (u *UI) PrintOutput(param entity.Log) {
	// u.Layout.OutputPanel.Buf.Insert(u.Layout.OutputPanel.Cursor.Loc, param.Content)

}
