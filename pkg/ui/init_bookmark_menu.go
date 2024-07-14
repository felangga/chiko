package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"chiko/pkg/entity"
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

	u.InitBookmarkMenu_SetInputCapture(bookmarkList)
	u.InitBookmarkMenu_SetSelection(bookmarkList)

	return bookmarkList
}

func (u *UI) InitBookmarkMenu_SetInputCapture(bookmarkList *tview.TreeView) {
	bookmarkList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			u.SetFocus(u.Layout.LogList)
		}
		return event
	})
}

func (u *UI) InitBookmarkMenu_SetSelection(bookmarkList *tview.TreeView) {
	// Initialize bookmark tree view
	bookmarkList.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			switch node.GetReference().(type) {
			case entity.Session:
				getSession := node.GetReference().(entity.Session)
				u.ShowBookmarkOptionsModal(u.Layout.MenuList, getSession)
			}

		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

}
