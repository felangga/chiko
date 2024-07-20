package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InitOutputPanel initializes the output panel on the main screen
func (u *UI) InitOutputPanel() *tview.Flex {
	outPanel := tview.NewTextArea()
	outPanel.SetTitle(" ⚙️ Output ")
	outPanel.SetWordWrap(true)
	outPanel.SetBorderPadding(1, 1, 1, 1)
	outPanel.SetMaxLength(1)

	// buttonPanel := u.buttonPanel()

	layout := tview.NewFlex()
	layout.SetDirection(tview.FlexRow)
	layout.AddItem(outPanel, 0, 1, false)
	// layout.AddItem(buttonPanel, 0, 1, false)
	layout.SetTitle(" Output ")
	layout.SetBorder(true)

	u.InitOutputPanel_SetInputCapture(outPanel)

	return layout
}

func (u *UI) buttonPanel() *tview.Grid {
	btnClearLog := tview.NewButton("Clear Log")
	btnClearLog.SetSelectedFunc(func() {
		u.Layout.LogList.SetText("...")
	})

	btnCopyClipboard := tview.NewButton("Copy to Clipboard")

	btnCopyClipboard.SetSelectedFunc(func() {
		// Copy to clipboard
	})

	buttonPanel := tview.NewGrid()
	buttonPanel.SetColumns(0, 0, 0)
	buttonPanel.SetRows(5)
	buttonPanel.SetBorderPadding(1, 1, 1, 1)

	buttonPanel.AddItem(btnClearLog, 0, 0, 1, 1, 0, 0, false)
	buttonPanel.AddItem(btnCopyClipboard, 0, 1, 1, 1, 0, 0, false)
	return buttonPanel
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
