package menu

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"chiko/pkg/ui"
)

type MenuUI struct {
	ParentUI *ui.UI
}

func (u *MenuUI) InitSidebarMenu() *tview.List {
	// Setup the side bar menu
	menuList := tview.NewList().ShowSecondaryText(false)
	menuList.SetBorder(true).SetTitle(" 🐶 Menu ")
	menuList.SetBorderPadding(1, 1, 1, 1)

	menuList.AddItem("Server URL", "", 'u', u.ShowSetServerURLModal)
	// menuList.AddItem("Methods", "", 'm', c.SetRequestMethods)
	// menuList.AddItem("Authorization", "", 'a', c.SetAuthorizationModal)
	// menuList.AddItem("Metadata", "", 'd', nil)
	// menuList.AddItem("Request Payload", "", 'p', c.SetRequestPayload)
	// menuList.AddItem("Invoke", "", 'i', c.DoInvoke)
	// menuList.AddItem("[::d]"+strings.Repeat(string(tcell.RuneHLine), 25), "", 0, nil)
	// menuList.AddItem("Save to Bookmark", "", 'b', c.DoSaveBookmark)
	menuList.AddItem("Quit", "", 'q', u.ParentUI.QuitApplication)

	// Handle keypress on menu list
	u.handleInputCapture(menuList)

	return menuList
}

func (u *MenuUI) handleInputCapture(menuList *tview.List) {
	menuList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			//u.SetFocus(u.BookmarkList)
		}
		return event
	})
}
