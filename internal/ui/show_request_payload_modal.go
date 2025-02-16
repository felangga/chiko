package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

// ShowRequestPayloadModal is used to show the request payload modal dialog
func (u *UI) ShowRequestPayloadModal() {
	if u.GRPC.Conn.ActiveConnection == nil {
		u.PrintLog(entity.Log{
			Content: "❗ no active connection",
			Type:    entity.LOG_WARNING,
		})

		return
	}

	if u.GRPC.Conn.SelectedMethod == nil {
		u.PrintLog(entity.Log{
			Content: "❗ please select rpc method first",
			Type:    entity.LOG_WARNING,
		})

		return
	}

	requestPayload := u.GRPC.Conn.RequestPayload

	// Create text area for filling the payload
	txtPayload := tview.NewTextArea().SetText(requestPayload, true)
	txtPayload.SetSize(9, 100)

	form := tview.NewForm()
	form.SetBackgroundColor(u.Theme.Colors.WindowColor)
	form.SetBorderPadding(1, 1, 0, 1)
	form.SetFieldStyle(u.Theme.Style.FieldStyle)
	form.AddFormItem(txtPayload)
	form.SetButtonsAlign(tview.AlignRight)
	form.SetButtonStyle(u.Theme.Style.ButtonStyle)

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
		out, err := u.GRPC.GenerateRPCPayloadSample()
		if err != nil {
			u.PrintLog(entity.Log{
				Content: fmt.Sprintf("❗ failed to generate sample: %s", err.Error()),
				Type:    entity.LOG_ERROR,
			})
		}

		txtPayload.SetText(out, true)
	})

	form.AddButton("Apply", func() {
		u.GRPC.Conn.RequestPayload = txtPayload.GetText()
		u.CloseModalDialog(wnd, u.Layout.MenuList)
	})
}
