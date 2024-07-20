package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InitOutputPanel initializes the output panel on the main screen
func (u *UI) InitOutputPanel() *tview.TextArea {
	outPanel := tview.NewTextArea()
	outPanel.SetTitle(" ⚙️ Output ")
	outPanel.SetBorder(true)
	outPanel.SetWordWrap(true)
	outPanel.SetBorderPadding(1, 1, 1, 1)

	u.InitOutputPanel_SetInputCapture(outPanel)

	return outPanel
}

func (u *UI) InitOutputPanel_SetInputCapture(outPanel *tview.TextArea) {
	outPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			u.App.SetFocus(u.Layout.MenuList)
		}
		return event
	})
}
