package ui

import (
	"fmt"
	"time"

	"github.com/felangga/chiko/internal/entity"
)

// PrintLog is used to print a log message to the log window
func (u *UI) PrintLog(param entity.Log) {
	warnaLog := "white"
	switch param.Type {
	case entity.LOG_ERROR:
		warnaLog = "red"
	case entity.LOG_WARNING:
		warnaLog = "yellow"
	}

	formatLog := fmt.Sprintf("[green][%s] [%s]%s [white]\n", time.Now().Format(time.RFC822), warnaLog, param.Content)

	go u.App.QueueUpdateDraw(func() {
		fmt.Fprint(u.Layout.LogList, formatLog)
		u.Layout.LogList.ScrollToEnd()
	})
}
