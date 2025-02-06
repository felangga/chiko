package ui

import (
	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/pkg/entity"
)

var (
	txtServerURL *tview.InputField
	btnConnect   *tview.Button
	btnTLS       *tview.Button
)

func (u *UI) ShowSetServerURLModal() {
	txtServerURL = tview.NewInputField()
	txtServerURL.SetBackgroundColor(u.Theme.Colors.WindowColor)
	txtServerURL.SetFieldStyle(u.Theme.Style.FieldStyle)
	txtServerURL.SetText(u.GRPC.Conn.ServerURL)

	btnConnect = tview.NewButton("Connect")
	btnConnect.SetStyle(u.Theme.Style.ButtonStyle)
	btnTLS = tview.NewButton("🔐 TLS")
	btnTLS.SetStyle(u.Theme.Style.ButtonStyle)

	layout := tview.NewGrid()
	layout.SetBorderPadding(1, 1, 1, 1)
	layout.SetBackgroundColor(u.Theme.Colors.WindowColor)
	layout.SetColumns(10, 0, 20)
	layout.SetRows(2, 1)
	layout.AddItem(txtServerURL, 0, 0, 1, 3, 0, 0, true)
	layout.AddItem(btnTLS, 1, 0, 1, 1, 0, 0, false)
	layout.AddItem(btnConnect, 1, 2, 1, 1, 0, 0, false)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " 🌏 Server URL ",
		rootView:      layout,
		draggable:     true,
		size:          winSize{0, 0, 70, 7},
		fallbackFocus: u.Layout.MenuList,
	})

	u.ShowSetServerURLModal_SetInputCapture(wnd)
}

// ShowSetServerURLModal_SetInputCapture handle the input capture from keyboard
func (u *UI) ShowSetServerURLModal_SetInputCapture(wnd *winman.WindowBase) {

	txtServerURL.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.WinMan.RemoveWindow(wnd)
			u.SetFocus(u.Layout.MenuList)
			return nil

		case tcell.KeyEnter:
			u.SetFocus(btnConnect)
			return nil

		case tcell.KeyTab:
			u.SetFocus(btnTLS)
			return nil
		}

		return event
	})

	btnTLS.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			u.SetFocus(btnConnect)
			return nil
		}

		return event
	})

	btnConnect.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			u.SetFocus(txtServerURL)
			return nil
		case tcell.KeyEnter:
			u.doConnect(wnd)
			return nil
		}

		return event
	})

	btnConnect.SetSelectedFunc(func() {
		u.doConnect(wnd)
	})

	btnTLS.SetSelectedFunc(func() {
		u.ShowCertificatePathModal(wnd)
	})
}

func (u *UI) doConnect(wnd *winman.WindowBase) {
	go func() {
		err := u.GRPC.Connect(txtServerURL.GetText())
		if err != nil {
			u.PrintLog(entity.Log{
				Content: "❌ failed to connect to [blue]" + txtServerURL.GetText() + " [red]" + err.Error(),
				Type:    entity.LOG_ERROR,
			})
		}
	}()

	// Remove the window and restore focus to menu list
	u.CloseModalDialog(wnd, u.Layout.MenuList)
}
