package ui

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

type winSize struct {
	x      int
	y      int
	width  int
	height int
}

type CreateModalDiaLog struct {
	title         string
	rootView      tview.Primitive
	draggable     bool
	resizeable    bool
	size          winSize
	fallbackFocus tview.Primitive
}

// CreateModalDialog is a helper to create a modal dialog window
func (u *UI) CreateModalDialog(param CreateModalDiaLog) *winman.WindowBase {
	wnd := winman.NewWindow().Show()

	wnd.SetTitle(param.title)
	wnd.SetRoot(param.rootView)
	wnd.SetDraggable(param.draggable)
	wnd.SetResizable(param.resizeable)
	wnd.SetModal(true)
	wnd.SetBackgroundColor(u.Theme.Colors.WindowColor)

	wnd.SetRect(param.size.x, param.size.y, param.size.width, param.size.height)
	wnd.AddButton(&winman.Button{
		Symbol: 'X',
		OnClick: func() {
			// Close current window and get back focus to the fallback primitive
			u.CloseModalDialog(wnd, param.fallbackFocus)
		},
	})

	u.WinMan.AddWindow(wnd)
	u.WinMan.Center(wnd)
	u.SetFocus(wnd)

	return wnd
}

// CloseModalDialog is a helper to close the modal dialog window
func (u *UI) CloseModalDialog(wnd *winman.WindowBase, focus tview.Primitive) {
	u.WinMan.RemoveWindow(wnd)
	u.SetFocus(focus)
}
