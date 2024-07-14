package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"

	"chiko/pkg/entity"
)

// ShowBookmarkOptionsModal is used to show the bookmark options modal, to load, overwrite, or remove the bookmark
func (u UI) ShowBookmarkOptionsModal(parentWnd tview.Primitive, bookmark entity.Session) {
	listOptions := tview.NewList()
	listOptions.ShowSecondaryText(false)
	listOptions.SetBackgroundColor(u.Theme.Colors.WindowColor)

	wnd := u.CreateModalDialog(CreateModalDiaLog{
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
		u.ApplyBookmark(param.bookmark)

		u.PrintLog(entity.Log{
			Content: fmt.Sprintf("📚 Bookmark loaded : %s", param.bookmark.Name),
			Type:    entity.LOG_INFO,
		})

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
						err := u.DeleteBookmark(param.bookmark)
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
					OnClick: func() {
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
		err := u.GRPC.CheckGRPC(u.GRPC.Conn.ServerURL)
		if err != nil {
			u.PrintLog(entity.Log{
				Content: err.Error(),
				Type:    entity.LOG_ERROR,
			})
			return
		}
	}()
}

// DeleteBookmark is used to delete a bookmark from the bookmark tree and save the bookmark
func (u *UI) DeleteBookmark(b entity.Session) error {
	// Helper function to remove a session from a bookmark
	removeSession := func(sessions []entity.Session, sessionID *uuid.UUID) []entity.Session {
		for i, session := range sessions {
			if session.ID == *sessionID {
				return append(sessions[:i], sessions[i+1:]...)
			}
		}
		return sessions
	}

	// Regenerate the boomark tree
	for i, bookmark := range *&u.Bookmark.Bookmarks {
		updatedSessions := removeSession(bookmark.Sessions, &b.ID)
		if len(updatedSessions) != len(bookmark.Sessions) {

			// Regenerate the bookmark list
			bookmark.Sessions = updatedSessions
			if len(bookmark.Sessions) == 0 {
				u.Bookmark.Bookmarks = append((u.Bookmark.Bookmarks)[:i], (u.Bookmark.Bookmarks)[i+1:]...)
			}

			err := u.Bookmark.SaveBookmark()
			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}
