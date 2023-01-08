package controller

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

func (c Controller) ShowMessageBox(title, message string) {
	// Message box content
	root := tview.NewFlex()

	txtMessage := tview.NewTextView()
	txtMessage.SetWordWrap(true)
	txtMessage.SetLabel(message)
	// txtMessage.SetBorderPadding(1, 1, 1, 1)
	root.AddItem(txtMessage, 0, 0, true)

	wnd := winman.NewWindow().
		Show().
		SetRoot(root).
		SetDraggable(true).
		SetTitle(title)

	wnd.SetModal(true)
	c.PrintLog(string(txtMessage.GetFieldWidth()), LOG_INFO)
	wnd.SetRect(0, 0, 50, 5)
	wnd.AddButton(&winman.Button{
		Symbol: '❌',
		OnClick: func() {
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		},
	})

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)

}
