package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Button struct {
	Name    string
	OnClick func()
}

type ShowMessageBoxParam struct {
	title   string
	message string
	buttons []Button
}

// ShowMessageBox is used to show message box with fixed size 50x11
func (u *UI) ShowMessageBox(param ShowMessageBoxParam) {

	txtMessage := tview.NewTextView()

	form := tview.NewForm()
	fmt.Fprint(txtMessage, param.message)
	form.AddFormItem(txtMessage)
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetButtonBackgroundColor(u.Theme.Colors.ButtonColor)
	form.SetBackgroundColor(u.Theme.Colors.WindowColor)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         param.title,
		rootView:      form,
		draggable:     true,
		resizeable:    false,
		size:          winSize{0, 0, 50, 11},
		fallbackFocus: u.Layout.MenuList,
	})

	u.ShowMessageBox_SetComponentActions(form, param.buttons)
	u.ShowMessageBox_SetInputCapture(wnd, form)
}

func (u *UI) ShowMessageBox_SetComponentActions(form *tview.Form, buttons []Button) {
	// Populate buttons
	for _, button := range buttons {
		form.AddButton(button.Name, button.OnClick)
	}
}

func (u *UI) ShowMessageBox_SetInputCapture(wnd *winman.WindowBase, parent *tview.Form) {
	parent.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.Layout.MenuList)
		}
		return event
	})

}
