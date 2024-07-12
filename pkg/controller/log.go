package controller

// import (
// 	"fmt"
// 	"time"
// )

// type LogType int8

// const (
// 	LOG_INFO    LogType = 0
// 	LOG_ERROR   LogType = 1
// 	LOG_WARNING LogType = 2
// )

// func (c Controller) PrintLog(log string, logType LogType) {
// 	// Get last log message
// 	lastLog := c.ui.OutputPanel.GetText(false)
// 	warnaLog := "white"
// 	switch logType {
// 	case LOG_ERROR:
// 		warnaLog = "red"
// 	case LOG_WARNING:
// 		warnaLog = "yellow"
// 	}
// 	formatLog := fmt.Sprintf("[green][%s] [%s]%s [white]\n", time.Now().Format(time.RFC822), warnaLog, log)
// 	c.ui.OutputPanel.SetWordWrap(true).SetText(lastLog + formatLog)

// 	// Scroll log window to bottom
// 	c.ui.OutputPanel.ScrollToEnd()
// }
