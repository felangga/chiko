package ui

import (
	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/pkg/entity"
)

var (
	txtServerURL  *tview.InputField
	chkSkipSecure *tview.Checkbox
	btnConnect    *tview.Button
	btnCertPath   *tview.Button
)

func (u *UI) ShowSetServerURLModal() {

	txtServerURL = tview.NewInputField()
	txtServerURL.SetBackgroundColor(u.Theme.Colors.WindowColor)
	txtServerURL.SetFieldBackgroundColor(u.Theme.Colors.PlaceholderColor)
	txtServerURL.SetText(u.GRPC.Conn.ServerURL)

	chkSkipSecure = tview.NewCheckbox().SetLabel("Skip Server Verification ")
	chkSkipSecure.SetBackgroundColor(u.Theme.Colors.WindowColor)
	chkSkipSecure.SetChecked(u.GRPC.Conn.InsecureSkipVerify)

	btnConnect = tview.NewButton("Connect")
	btnCertPath = tview.NewButton("SSL Certicates")

	
	layout := tview.NewGrid()
	layout.SetBorderPadding(1, 1, 1, 1)
	layout.SetBackgroundColor(u.Theme.Colors.WindowColor)
	layout.SetColumns(0, 0, 0)
	layout.SetRows(2, 1, 0, 1)
	layout.AddItem(txtServerURL, 0, 0, 1, 3, 0, 0, true)
	layout.AddItem(chkSkipSecure, 1, 0, 1, 2, 0, 0, false)
	layout.AddItem(btnCertPath, 1, 2, 1, 1, 0, 0, false)
	layout.AddItem(btnConnect, 3, 2, 1, 1, 0, 0, false)

	wnd := u.CreateModalDialog(CreateModalDiaLog{
		title:         " 🌏 Server URL ",
		rootView:      layout,
		draggable:     true,
		size:          winSize{0, 0, 70, 10},
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
			u.SetFocus(chkSkipSecure)
			return nil
		}

		return event
	})

	chkSkipSecure.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			u.SetFocus(btnCertPath)
			return nil
		}

		return event
	})

	btnCertPath.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

	btnCertPath.SetSelectedFunc(func() {
		u.ShowCertificatePathModal(wnd)
	})
}

func (u *UI) doConnect(wnd *winman.WindowBase) {
	go func() {
		u.GRPC.Conn.InsecureSkipVerify = chkSkipSecure.IsChecked()
		err := u.GRPC.Connect(txtServerURL.GetText())
		if err != nil {
			u.PrintLog(entity.Log{
				Content: "❌ failed to connect to [blue]" + txtServerURL.GetText() + " [red]" + err.Error(),
				Type:    entity.LOG_ERROR,
			})
		}
	}()

	// Remove the window and restore focus to menu list
	u.WinMan.RemoveWindow(wnd)
	u.SetFocus(u.Layout.MenuList)
}
