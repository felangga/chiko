package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"chiko/pkg/entity"
)

// ShowBookmarkOptionsModal is used to show the bookmark options modal, to load, overwrite, or remove the bookmark
func (u UI) ShowBookmarkOptionsModal(parentWnd tview.Primitive, bookmark entity.Session) {
	listOptions := tview.NewList()
	listOptions.ShowSecondaryText(false)
	listOptions.SetBackgroundColor(u.Theme.Colors.WindowColor)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " 📚 Bookmark Options ",
		rootView:      listOptions,
		draggable:     true,
		resizeable:    false,
		size:          winSize{0, 0, 50, 7},
		fallbackFocus: parentWnd,
	})

	u.populateBookmarkChoices(populateBookmarkChoicesParam{
		listOptions: listOptions,
		wnd:         wnd,
		parentWnd:   &parentWnd,
		bookmark:    bookmark,
	})

	u.ShowBookmarkOptionsModal_SetInputCapture(wnd, parentWnd, listOptions)

}

// ShowBookmarkOptionsModal_SetInputCapture sets the input capture for the bookmark options modal
func (u *UI) ShowBookmarkOptionsModal_SetInputCapture(wnd *winman.WindowBase, parentWnd tview.Primitive, listOptions *tview.List) {
	listOptions.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			u.CloseModalDialog(wnd, parentWnd)
		}
		return event
	})
}

type populateBookmarkChoicesParam struct {
	listOptions *tview.List
	wnd         *winman.WindowBase
	parentWnd   *tview.Primitive
	bookmark    entity.Session
}

// populateBookmarkChoices is used to populate the bookmark choices on the bookmark modal
func (u *UI) populateBookmarkChoices(param populateBookmarkChoicesParam) {
	// Load bookmark to the current session
	param.listOptions.AddItem("Load Bookmark", "", 'a', func() {
		u.Controller.ApplyBookmark(param.bookmark)

		// Close the window
		u.CloseModalDialog(param.wnd, *param.parentWnd)
	})

	// Overwrite bookmark from the current session
	param.listOptions.AddItem("Overwrite Bookmark", "", 'o', func() {
		// c.overwriteBookmark(index)
	})

	param.listOptions.AddItem("Delete Bookmark", "", 'd', func() {
		u.ShowMessageBox(ShowMessageBoxParam{
			title:   " 🗑️ Delete Bookmark ",
			message: "Are you sure you want to delete this bookmark?",
			buttons: []Button{
				{
					Name: "Yes",
					OnClick: func() {
						err := u.Controller.DeleteBookmark(param.bookmark)
						if err != nil {
							u.PrintLog(entity.LogParam{
								Content: fmt.Sprintf("❌ failed to delete bookmark, err: %v", err),
								Type:    entity.LOG_ERROR,
							})
						}

						u.PrintLog(entity.LogParam{
							Content: fmt.Sprintf("✅ bookmark [blue]%s [white]deleted", param.bookmark.Name),
							Type:    entity.LOG_INFO,
						})

						// Close the window
						u.CloseModalDialog(param.wnd, *param.parentWnd)
					},
				},
				{
					Name: "No",
					OnClick: func() {
						u.CloseModalDialog(param.wnd, *param.parentWnd)
					},
				},
			},
		})

	})
}
