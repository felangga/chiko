package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InitSidebarMenu is used to initialize and populate sidebar menu
func (u *UI) InitSidebarMenu() *tview.List {

	menuList := tview.NewList().ShowSecondaryText(false)
	menuList.SetBorder(true).SetTitle(" 🐶 Menu ")
	menuList.SetBorderPadding(1, 1, 1, 1)

	menuList.AddItem("Server URL", "", 'u', u.ShowSetServerURLModal)
	menuList.AddItem("Methods", "", 'm', u.ShowSetRequestMethodModal)
	menuList.AddItem("Authorization", "", 'a', u.ShowAuthorizationModal)
	menuList.AddItem("Metadata", "", 'd', u.ShowMetadataModal)
	menuList.AddItem("Request Payload", "", 'p', u.ShowRequestPayloadModal)
	menuList.AddItem("Invoke", "", 'i', u.InvokeRPC)
	menuList.AddItem("[::d]"+strings.Repeat(string(tcell.RuneHLine), 25), "", 0, nil)
	menuList.AddItem("Create New Bookmark", "", 'b', u.ShowSaveToBookmarkModal)
	menuList.AddItem("Quit", "", 'q', u.QuitApplication)

	// Handle keypress on menu list
	u.InitSidebarMenu_SetInputCapture(menuList)

	return menuList
}

func (u *UI) InitSidebarMenu_SetInputCapture(menuList *tview.List) {
	menuList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			u.SetFocus(u.Layout.BookmarkList)
			return nil
		}
		return event
	})
}
