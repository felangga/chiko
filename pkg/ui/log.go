package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InitLogList initializes the log panel on the main screen
func (u *UI) InitLogList() *tview.TextView {
	logPanel := tview.NewTextView()
	logPanel.SetDynamicColors(true)
	logPanel.SetTitle(" 📃 Logs ")
	logPanel.SetBorder(true)
	logPanel.SetWordWrap(true)
	logPanel.SetBorderPadding(1, 1, 1, 1)

	logPanel.SetScrollable(true).SetChangedFunc(func() {
		u.App.Draw()
	})

	u.handleLogInputCapture(logPanel)

	return logPanel
}

func (u *UI) handleLogInputCapture(logPanel *tview.TextView) {
	logPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			u.App.SetFocus(u.Layout.MenuList)
		}
		return event
	})
}
