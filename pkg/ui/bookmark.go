package ui

import "github.com/rivo/tview"

// InitBookmarkMenu initializes the bookmark sidebar menu
func (u *UI) InitBookmarkMenu() *tview.TreeView {
	treeRoot := tview.NewTreeNode("📚 Library")
	bookmarkList := tview.NewTreeView().
		SetRoot(treeRoot).
		SetCurrentNode(treeRoot)

	bookmarkList.SetBorder(true)
	bookmarkList.SetBorderPadding(1, 1, 1, 1)
	bookmarkList.SetTitle(" 📚 Bookmarks Library ")

	return bookmarkList
}
