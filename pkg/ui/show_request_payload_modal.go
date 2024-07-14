package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"chiko/pkg/entity"
)

// ShowRequestPayloadModal is used to show the request payload modal dialog
func (u *UI) ShowRequestPayloadModal() {
	if u.Controller.Conn.ActiveConnection == nil {
		u.PrintLog(entity.LogParam{
			Content: "❗ no active connection",
			Type:    entity.LOG_WARNING,
		})

		return
	}

	if u.Controller.Conn.SelectedMethod == nil {
		u.PrintLog(entity.LogParam{
			Content: "❗ please select rpc method first",
			Type:    entity.LOG_WARNING,
		})

		return
	}

	requestPayload := u.Controller.Conn.RequestPayload

	// Create text area for filling the payload
	txtPayload := tview.NewTextArea().SetText(requestPayload, true)
	txtPayload.SetSize(9, 100)

	form := tview.NewForm()
	form.SetBackgroundColor(u.Theme.Colors.WindowColor)
	form.SetBorderPadding(1, 1, 0, 1)
	form.SetFieldBackgroundColor(u.Theme.Colors.FieldColor)
	form.AddFormItem(txtPayload)
	form.SetButtonsAlign(tview.AlignRight)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " 📦 Request Payload ",
		rootView:      form,
		draggable:     true,
		resizeable:    false,
		size:          winSize{0, 0, 70, 15},
		fallbackFocus: u.Layout.MenuList,
	})

	u.ShowRequestPayloadModal_SetInputCapture(wnd, form)
	u.ShowRequestPayloadModal_SetComponentActions(wnd, form, txtPayload)

}

func (u *UI) ShowRequestPayloadModal_SetInputCapture(wnd *winman.WindowBase, form *tview.Form) {
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.Layout.MenuList)
			return nil
		}
		return event
	})
}

func (u *UI) ShowRequestPayloadModal_SetComponentActions(wnd *winman.WindowBase, form *tview.Form, txtPayload *tview.TextArea) {
	form.AddButton("Generate Sample", func() {
		out, err := u.Controller.GenerateRPCPayloadSample()
		if err != nil {
			u.PrintLog(entity.LogParam{
				Content: fmt.Sprintf("❗ failed to generate sample: %s", err.Error()),
				Type:    entity.LOG_ERROR,
			})
		}

		txtPayload.SetText(out, true)
	})

	form.AddButton("Apply", func() {
		u.Controller.Conn.RequestPayload = txtPayload.GetText()

		// Remove the window and restore focus to menu list
		u.PrintLog(entity.LogParam{
			Content: "\nRequest Payload:\n[yellow]" + u.Controller.Conn.RequestPayload,
			Type:    entity.LOG_INFO,
		})
		u.CloseModalDialog(wnd, u.Layout.MenuList)

	})
}
