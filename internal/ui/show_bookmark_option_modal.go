package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

// ShowBookmarkOptionsModal is used to show the bookmark options modal, to load, overwrite, or remove the bookmark
func (u UI) ShowBookmarkOptionsModal(parentWnd tview.Primitive, bookmark *entity.Session) {
	listOptions := tview.NewList()
	listOptions.ShowSecondaryText(false)
	listOptions.SetBackgroundColor(u.Theme.Colors.WindowColor)

	st := tcell.StyleDefault
	listOptions.SetMainTextStyle(st.Background(u.Theme.Colors.WindowColor).Foreground(tcell.ColorWhite))
	listOptions.SetShortcutStyle(st.Background(u.Theme.Colors.WindowColor).Foreground(tcell.ColorYellow))

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
	bookmark    *entity.Session
}

// populateBookmarkChoices is used to populate the bookmark choices on the bookmark modal
func (u *UI) populateBookmarkChoices(param populateBookmarkChoicesParam) {
	// Load bookmark to the current session
	param.listOptions.AddItem("Load Bookmark", "", 'a', func() {
		u.ApplyBookmark(*param.bookmark)

		u.PrintLog(entity.Log{
			Content: fmt.Sprintf("📚 Bookmark loaded : %s", param.bookmark.Name),
			Type:    entity.LOG_INFO,
		})

		// Close the window
		u.CloseModalDialog(param.wnd, *param.parentWnd)
	})

	// Overwrite bookmark from the current session
	param.listOptions.AddItem("Overwrite Bookmark", "", 'o', func() {
		u.ShowMessageBox(ShowMessageBoxParam{
			title:   " 📝 Overwrite Bookmark ",
			message: "Are you sure you want to overwrite this bookmark?",
			buttons: []Button{
				{
					Name: "Yes",
					OnClick: func(wnd *winman.WindowBase) {
						err := u.OverwriteBookmark(param.bookmark)
						if err != nil {
							u.PrintLog(entity.Log{
								Content: fmt.Sprintf("❌ failed to overwrite bookmark, err: %v", err),
								Type:    entity.LOG_ERROR,
							})
						}

						u.PrintLog(entity.Log{
							Content: fmt.Sprintf("✅ bookmark [blue]%s [white]overwritten", param.bookmark.Name),
							Type:    entity.LOG_INFO,
						})

						// Close the window
						u.CloseModalDialog(param.wnd, *param.parentWnd)
					},
				},
				{
					Name: "No",
					OnClick: func(wnd *winman.WindowBase) {
						u.CloseModalDialog(wnd, *param.parentWnd)
					},
				},
			},
		})
	})

	param.listOptions.AddItem("Delete Bookmark", "", 'd', func() {
		u.ShowMessageBox(ShowMessageBoxParam{
			title:   " 🗑️ Delete Bookmark ",
			message: "Are you sure you want to delete this bookmark?",
			buttons: []Button{
				{
					Name: "Yes",
					OnClick: func(wnd *winman.WindowBase) {
						err := u.DeleteBookmark(*param.bookmark)
						if err != nil {
							u.PrintLog(entity.Log{
								Content: fmt.Sprintf("❌ failed to delete bookmark, err: %v", err),
								Type:    entity.LOG_ERROR,
							})
						}

						u.PrintLog(entity.Log{
							Content: fmt.Sprintf("✅ bookmark [blue]%s [white]deleted", param.bookmark.Name),
							Type:    entity.LOG_INFO,
						})

						// Close the window
						u.CloseModalDialog(param.wnd, *param.parentWnd)
					},
				},
				{
					Name: "No",
					OnClick: func(wnd *winman.WindowBase) {
						u.CloseModalDialog(param.wnd, *param.parentWnd)
					},
				},
			},
		})

	})
}

func (u *UI) ApplyBookmark(session entity.Session) {
	// Get selected connection
	*u.GRPC.Conn = session

	go func() {
		err := u.GRPC.Connect(u.GRPC.Conn.ServerURL)
		if err != nil {
			u.PrintLog(entity.Log{
				Content: "❌ failed to connect to [blue]" + u.GRPC.Conn.ServerURL + " [red]" + err.Error(),
				Type:    entity.LOG_ERROR,
			})

			return
		}
	}()
}

// DeleteBookmark is used to delete a bookmark from the bookmark tree and save the bookmark
func (u *UI) DeleteBookmark(b entity.Session) error {
	// Helper function to remove a session from a bookmark
	removeSession := func(sessions []entity.Session, sessionID uuid.UUID) []entity.Session {
		for i, session := range sessions {
			if session.ID == sessionID {
				return append(sessions[:i], sessions[i+1:]...)
			}
		}
		return sessions
	}

	for i := 0; i < len(*u.Bookmark.Categories); i++ {
		bookmark := &(*u.Bookmark.Categories)[i] // Get the address of the bookmark to modify the original element
		updatedSessions := removeSession(bookmark.Sessions, b.ID)
		if len(updatedSessions) != len(bookmark.Sessions) {

			// Regenerate the bookmark list
			bookmark.Sessions = updatedSessions

			err := u.Bookmark.SaveBookmark()
			if err != nil {
				return err
			}

			// Refresh the bookmark list
			_ = u.RefreshBookmarkList()

			return nil
		}
	}

	return nil
}

// OverwriteBookmark is used to overwrite the bookmark with current active session
func (u *UI) OverwriteBookmark(b *entity.Session) error {
	for i := 0; i < len(*u.Bookmark.Categories); i++ {
		category := &(*u.Bookmark.Categories)[i]
		for j, session := range category.Sessions {
			if session.ID == b.ID {
				category.Sessions[j] = *u.GRPC.Conn
				category.Sessions[j].Name = session.Name

				return u.Bookmark.SaveBookmark()
			}
		}
	}
	return nil
}
