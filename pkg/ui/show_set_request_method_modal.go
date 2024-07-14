package ui

import (
	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"

	"chiko/pkg/entity"
)

// ShowSetRequestMethodModal is used to show the RPC method selection modal dialog
func (u *UI) ShowSetRequestMethodModal() {
	if u.GRPC.Conn.ActiveConnection == nil {
		u.PrintLog(entity.Log{
			Content: "❗ no active connection",
			Type:    entity.LOG_WARNING,
		})

		return
	}

	// Set placeholder text and style
	style := tcell.StyleDefault.
		Background(u.Theme.Colors.PlaceholderColor).
		Italic(true)

	listMethods := tview.NewList().
		ShowSecondaryText(false)
	listMethods.SetBorderPadding(1, 0, 0, 0)
	listMethods.SetBackgroundColor(u.Theme.Colors.WindowColor)

	txtSearch := tview.NewInputField().
		SetText("").
		SetPlaceholder(" 🔍 Search methods...").
		SetPlaceholderStyle(style)

	txtSearch.SetFieldBackgroundColor(u.Theme.Colors.PlaceholderColor)

	wndLayer := tview.NewFlex()
	wndLayer.SetDirection(tview.FlexRow)
	wndLayer.AddItem(txtSearch, 1, 1, true)
	wndLayer.AddItem(listMethods, 0, 1, false)

	wnd := u.CreateModalDialog(CreateModalDiaLog{
		title:         " 📡 Select RPC Methods ",
		rootView:      wndLayer,
		draggable:     true,
		resizeable:    true,
		size:          winSize{0, 0, 70, 20},
		fallbackFocus: u.Layout.MenuList,
	})
	wnd.SetBorderPadding(1, 1, 1, 1)

	u.ShowSetRequestMethodModal_SetInputCapture(wnd, listMethods, txtSearch)
	u.refreshRequestMethodList(listMethods, u.GRPC.Conn.AvailableMethods)

	// Customize window to have maximize button
	var maxMinButton *winman.Button
	maxMinButton = &winman.Button{
		Symbol:    '▴',
		Alignment: winman.ButtonRight,
		OnClick: func() {
			if wnd.IsMaximized() {
				wnd.Restore()
				maxMinButton.Symbol = '▴'
			} else {
				wnd.Maximize()
				maxMinButton.Symbol = '▾'
			}
		},
	}
	wnd.AddButton(maxMinButton)
}

func (u *UI) ShowSetRequestMethodModal_SetInputCapture(wnd *winman.WindowBase, listMethods *tview.List, txtSearch *tview.InputField) {
	txtSearch.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		getText := txtSearch.GetText()

		switch event.Key() {
		case tcell.KeyDown, tcell.KeyEnter, tcell.KeyTAB:
			u.SetFocus(listMethods)
			return nil
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.Layout.MenuList)
			return nil
		case tcell.KeyRune:
			getText = getText + string(event.Rune())
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if len(getText) > 0 {
				getText = getText[:len(getText)-1]
			}
		}

		fuzzyFind := fuzzy.FindFold(getText, u.GRPC.Conn.AvailableMethods)
		u.refreshRequestMethodList(listMethods, fuzzyFind)

		return event
	})

	listMethods.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.Layout.MenuList)
		case tcell.KeyEnter:
			idx := listMethods.GetCurrentItem()
			item, _ := listMethods.GetItemText(idx)
			u.GRPC.Conn.SelectedMethod = &item

			// Reset the search bar
			txtSearch.SetText("")

			// Remove the window and restore focus to menu list
			u.PrintLog(entity.Log{
				Content: "👉 Method set to [blue]" + *u.GRPC.Conn.SelectedMethod,
				Type:    entity.LOG_INFO,
			})

			u.CloseModalDialog(wnd, u.Layout.MenuList)
		case tcell.KeyRune, tcell.KeyBackspace, tcell.KeyBackspace2:
			// Set focus to input field when user typing something
			u.SetFocus(txtSearch)
		case tcell.KeyTAB:
			u.SetFocus(txtSearch)
			return nil
		}
		return event
	})

}

func (u *UI) refreshRequestMethodList(listView *tview.List, items []string) {
	listView.Clear()
	listView.SetCurrentItem(0)

	for i, method := range items {
		listView.AddItem(method, "", 0, nil)

		// Ignore if none was selected before
		if u.GRPC.Conn.SelectedMethod != nil && *u.GRPC.Conn.SelectedMethod == method {
			listView.SetCurrentItem(i)
		}
	}

}
