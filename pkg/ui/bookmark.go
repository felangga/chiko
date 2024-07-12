package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InitBookmarkMenu initializes the bookmark sidebar menu
func (u *UI) InitBookmarkMenu() *tview.TreeView {
	treeRoot := tview.NewTreeNode("📚 Library")
	bookmarkList := tview.NewTreeView().
		SetRoot(treeRoot).
		SetCurrentNode(treeRoot)

	bookmarkList.SetBorder(true)
	bookmarkList.SetBorderPadding(1, 1, 1, 1)
	bookmarkList.SetTitle(" 📚 Bookmarks Library ")

	u.handleBookmarkInputCapture(bookmarkList)
	return bookmarkList
}

func (u *UI) handleBookmarkInputCapture(bookmarkList *tview.TreeView) {
	bookmarkList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			u.SetFocus(u.Layout.LogList)
		}
		return event
	})
}
