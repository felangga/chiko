package ui

import (
	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/pkg/entity"
)

func (u *UI) ShowSetServerURLModal() {
	txtServerURL := tview.NewInputField().SetText(u.GRPC.Conn.ServerURL)
	txtServerURL.SetFieldBackgroundColor(u.Theme.Colors.WindowColor)

	wnd := u.CreateModalDialog(CreateModalDiaLog{
		title:         " 🌏 Enter Server URL ",
		rootView:      txtServerURL,
		draggable:     true,
		size:          winSize{0, 0, 70, 1},
		fallbackFocus: u.Layout.MenuList,
	})

	u.ShowSetServerURLModal_SetInputCapture(wnd, txtServerURL)
}

// ShowSetServerURLModal_SetInputCapture handle the input capture from keyboard
func (u *UI) ShowSetServerURLModal_SetInputCapture(wnd *winman.WindowBase, inp *tview.InputField) {
	inp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.WinMan.RemoveWindow(wnd)
			u.SetFocus(u.Layout.MenuList)
			return nil

		case tcell.KeyEnter:
			go func() {
				err := u.GRPC.Connect(inp.GetText())
				if err != nil {
					u.PrintLog(entity.Log{
						Content: "❌ failed to connect to [blue]" + inp.GetText() + " [red]" + err.Error(),
						Type:    entity.LOG_ERROR,
					})
				}
			}()

			// Remove the window and restore focus to menu list
			u.WinMan.RemoveWindow(wnd)
			u.SetFocus(u.Layout.MenuList)
		}
		return event
	})

}
