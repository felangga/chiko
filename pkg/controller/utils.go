package controller

// import (
// 	"fmt"

// 	"github.com/epiclabs-io/winman"
// 	"github.com/gdamore/tcell/v2"
// 	"github.com/rivo/tview"
// )

// type Button struct {
// 	Name    string
// 	OnClick func()
// }

// // ShowMessageBox is used to show message box with fixed size 50x11
// func (c Controller) ShowMessageBox(title, message string, buttons []Button) {
// 	// Message box content
// 	root := tview.NewForm()
// 	txtMessage := tview.NewTextView()
// 	fmt.Fprint(txtMessage, message)
// 	root.AddFormItem(txtMessage)

// 	// Populate buttons
// 	for _, button := range buttons {
// 		root.AddButton(button.Name, button.OnClick)
// 	}

// 	root.SetButtonsAlign(tview.AlignCenter)
// 	root.SetButtonBackgroundColor(c.theme.Colors.ButtonColor)
// 	root.SetBackgroundColor(c.theme.Colors.ModalColor)

// 	wnd := winman.NewWindow().
// 		Show().
// 		SetRoot(root).
// 		SetDraggable(true).
// 		SetResizable(true).
// 		SetTitle(title)

// 	wnd.SetModal(true)
// 	wnd.AddButton(&winman.Button{
// 		Symbol: 'X',
// 		OnClick: func() {
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		},
// 	})

// 	// TODO: need to create auto-sized window to accomodate long message text
// 	wnd.SetRect(0, 0, 50, 11)
// 	wnd.SetResizable(false)
// 	wnd.SetBackgroundColor(c.theme.Colors.ModalColor)

// 	// Handle keypress on window
// 	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		switch event.Key() {
// 		case tcell.KeyEscape:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		}
// 		return event
// 	})

// 	c.ui.WinMan.AddWindow(wnd)
// 	c.ui.WinMan.Center(wnd)
// 	c.ui.SetFocus(wnd)

// }
