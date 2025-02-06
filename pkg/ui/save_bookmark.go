package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/pkg/entity"
)

// ShowSaveToBookmarkModal used to open the save bookmark dialog to save the current payload to bookmark
func (u *UI) ShowSaveToBookmarkModal() {
	// Initiate saving bookmark sequence
	u.ShowBookmarkCategoryModal(func(wnd winman.Window, b *entity.Category) {
		u.ShowBookmarkNameModal(wnd, func(bookmarkName string) {
			genID := uuid.New()
			u.GRPC.Conn.ID = genID
			u.GRPC.Conn.Name = bookmarkName

			b.Sessions = append(b.Sessions, *u.GRPC.Conn)

			err := u.Bookmark.SaveBookmark()
			if err != nil {
				u.PrintLog(entity.Log{
					Content: "❌ failed to save bookmark",
					Type:    entity.LOG_ERROR,
				})
				return
			}

			u.PrintLog(entity.Log{
				Content: fmt.Sprintf("✅ [white]bookmark [blue]%s [white]saved", bookmarkName),
				Type:    entity.LOG_ERROR,
			})

			u.RefreshBookmarkList()
		})

		u.WinMan.RemoveWindow(wnd)
	})
}

// ShowBookmarkNameModal is used to show modal with text box to change the bookmark name
func (u *UI) ShowBookmarkNameModal(parentWND winman.Window, onEnter func(bookmarkName string)) {
	bookmarkName := tview.NewInputField().SetText(u.GRPC.Conn.Name)
	bookmarkName.SetFieldBackgroundColor(u.Theme.Colors.WindowColor)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " Enter Bookmark Name ",
		rootView:      bookmarkName,
		draggable:     true,
		size:          winSize{0, 0, 80, 1},
		fallbackFocus: parentWND,
	})

	bookmarkName.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, parentWND)
		case tcell.KeyEnter:
			onEnter(bookmarkName.GetText())

			// Remove the window and restore focus to menu list
			u.CloseModalDialog(wnd, u.Layout.MenuList)
		}
		return event
	})
}

func (u *UI) ShowBookmarkCategoryModal(onSelectedCategory func(wnd winman.Window, b *entity.Category)) {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBackgroundColor(u.Theme.Colors.WindowColor)
	list.SetMainTextStyle(u.Theme.Style.ListMainTextStyle)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " 📚 Select Bookmark Category ",
		rootView:      list,
		draggable:     true,
		size:          winSize{0, 0, 50, 10},
		fallbackFocus: u.Layout.MenuList,
	})

	wnd.SetBorderPadding(1, 1, 1, 1)

	// Add selection for user if the user wants to create a new category
	list.AddItem("📖 Create New Category", "", 0, func() {
		u.ShowCreateNewCategoryModal(wnd, list, onSelectedCategory)
	})

	for i := range *u.Bookmark.Categories {
		// Get the address of the newly appended category
		categoryPtr := &(*u.Bookmark.Categories)[i]
		list.AddItem("📁 "+categoryPtr.Name, "", 0, func() {
			onSelectedCategory(wnd, categoryPtr)
		})
	}

	u.ShowBookmarkCategoryModal_SetInputCapture(wnd, list)
}

func (u *UI) ShowBookmarkCategoryModal_SetInputCapture(wnd *winman.WindowBase, list *tview.List) {
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.Layout.MenuList)
		}
		return event
	})
}

// ShowCreateNewCategoryModal is used to show modal with text box to create new category
func (u *UI) ShowCreateNewCategoryModal(parentWND *winman.WindowBase, list *tview.List, onSelectedCategory func(wnd winman.Window, b *entity.Category)) {
	catName := tview.NewInputField().SetText(u.GRPC.Conn.Name)
	catName.SetFieldBackgroundColor(u.Theme.Colors.WindowColor)

	mdlNewCategory := u.CreateModalDialog(CreateModalDialogParam{
		title:         " 📁 Enter New Category Name ",
		rootView:      catName,
		draggable:     true,
		size:          winSize{0, 0, 50, 1},
		fallbackFocus: parentWND,
	})

	catName.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(mdlNewCategory, parentWND)
		case tcell.KeyEnter:
			// Add new category to bookmark list
			newCategory := entity.Category{
				Name:     catName.GetText(),
				Sessions: []entity.Session{},
			}
			*u.Bookmark.Categories = append(*u.Bookmark.Categories, newCategory)

			// Get the address of the newly appended category
			newCategoryPtr := &(*u.Bookmark.Categories)[len(*u.Bookmark.Categories)-1]

			list.AddItem("📁 "+newCategory.Name, "", 0, func() {
				onSelectedCategory(parentWND, newCategoryPtr)
			})

			// Remove the window and restore focus to previous window
			u.CloseModalDialog(mdlNewCategory, parentWND)
		}
		return event
	})

}
