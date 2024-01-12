package controller

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

type Button struct {
	Name    string
	OnClick func()
}

func (c Controller) ShowMessageBox(title, message string, buttons []Button) {
	// Message box content
	root := tview.NewForm()
	root.SetRect(0, 0, 20, 20)
	// root.SetDirection(tview.NewGrid())

	txtMessage := tview.NewTextView()
	txtMessage.SetWordWrap(true)
	txtMessage.SetLabel(message)
	root.AddFormItem(txtMessage)

	// Populate button
	for _, button := range buttons {
		root.AddButton(button.Name, button.OnClick)
	}
	root.SetButtonsAlign(tview.AlignCenter)

	wnd := winman.NewWindow().
		Show().
		SetRoot(root).
		SetDraggable(true).
		SetResizable(true).
		SetTitle(title)

	wnd.SetModal(true)
	wnd.SetRect(0, 0, 50, 20)
	wnd.AddButton(&winman.Button{
		Symbol: 'X',
		OnClick: func() {
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		},
	})

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)

}
