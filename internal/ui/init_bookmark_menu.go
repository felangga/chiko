package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
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
			u.SetFocus(u.ActiveSession.OutputPanel.Layout)
			return nil
		case tcell.KeyDelete:
			// Delete key directly shows the delete confirmation for the focused node
			currentNode := bookmarkList.GetCurrentNode()
			if currentNode == nil {
				return nil
			}
			switch ref := currentNode.GetReference().(type) {
			case entity.Category:
				u.ShowDeleteCategoryConfirmation(ref)
			case *entity.Session:
				u.ShowDeleteBookmarkConfirmation(ref)
			}
			return nil
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

		switch ref := node.GetReference().(type) {
		case *entity.Session:
			u.ShowBookmarkOptionsModal(u.Layout.BookmarkList, ref)
		case entity.Category:
			node.SetExpanded(!node.IsExpanded())
		}
	})
}
