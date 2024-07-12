package controller

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"

// 	"github.com/epiclabs-io/winman"
// 	"github.com/gdamore/tcell/v2"
// 	"github.com/google/uuid"
// 	"github.com/rivo/tview"

// 	"chiko/pkg/entity"
// )

// func (c Controller) applyBookmark(session entity.Session) {
// 	// Get selected connection
// 	*c.conn = session

// 	c.PrintLog(fmt.Sprintf("📚 Bookmark loaded : %s", c.conn.ServerURL), LOG_INFO)

// 	go c.CheckGRPC(c.conn.ServerURL)
// }

// // loadBookmarks is used to load bookmarks from bookmark file
// func (c Controller) loadBookmarks() {
// 	if _, err := os.Stat(entity.BOOKMARKS_FILE_NAME); err != nil {
// 		c.PrintLog(err.Error(), LOG_INFO)
// 		return
// 	}

// 	// Read bookmark file
// 	file, err := os.ReadFile(entity.BOOKMARKS_FILE_NAME)
// 	if err != nil {
// 		c.PrintLog(err.Error(), LOG_ERROR)
// 		return
// 	}

// 	err = json.Unmarshal([]byte(file), &c.bookmarks)
// 	if err != nil {
// 		c.PrintLog(err.Error(), LOG_ERROR)
// 		return
// 	}

// 	// Populate bookmarks list
// 	for _, b := range *c.bookmarks {
// 		categoryNode := tview.NewTreeNode("📁 " + b.CategoryName)
// 		categoryNode.SetReference(b)

// 		for _, session := range b.Sessions {
// 			sessionNode := tview.NewTreeNode("📗 " + session.Name)
// 			sessionNode.SetReference(session)
// 			categoryNode.AddChild(sessionNode)
// 		}

// 		c.ui.BookmarkList.GetRoot().AddChild(categoryNode)
// 	}

// 	c.PrintLog(fmt.Sprintf("📚 %d bookmark(s) loaded", len(*c.bookmarks)), LOG_INFO)
// }

// // SaveBookmark is used to save the bookmark object to file by encoding the object with JSON.
// func (c Controller) SaveBookmark() {
// 	// Encoding the object to JSON
// 	convert, err := json.Marshal(c.bookmarks)
// 	if err != nil {
// 		c.PrintLog("❌ failed to prepare bookmarks data", LOG_ERROR)
// 		return
// 	}

// 	// Saving the json to file
// 	err = os.WriteFile(entity.BOOKMARKS_FILE_NAME, convert, 0644)
// 	if err != nil {
// 		c.PrintLog("💾 failed to write bookmark configuration, maybe write-protected?", LOG_ERROR)
// 		return
// 	}
// }

// func (c Controller) ShowBookmarkCategoryModal(onSelectedCategory func(wnd winman.Window, b *entity.Bookmark)) {
// 	list := tview.NewList().
// 		ShowSecondaryText(false)

// 	wnd := winman.NewWindow().
// 		Show().
// 		SetModal(true).
// 		SetRoot(list).
// 		SetDraggable(true).
// 		SetTitle(" 📚 Select Bookmark Category ")

// 	wnd.SetBorderPadding(1, 1, 1, 1)

// 	createCategoryModal := func() {
// 		catName := tview.NewInputField().SetText(c.conn.Name)
// 		mdlNewCategory := winman.NewWindow().
// 			SetModal(true).
// 			Show().
// 			SetRoot(catName).
// 			SetDraggable(true).
// 			SetTitle(" ✏️ Enter New Category Name ")

// 		mdlNewCategory.SetBackgroundColor(c.theme.Colors.WindowColor)
// 		catName.SetFieldBackgroundColor(wnd.GetBackgroundColor())
// 		catName.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 			switch event.Key() {
// 			case tcell.KeyEscape:
// 				c.ui.WinMan.RemoveWindow(mdlNewCategory)
// 				c.ui.SetFocus(wnd)
// 			case tcell.KeyEnter:
// 				// Add new category to bookmark list
// 				newCategory := entity.Bookmark{
// 					CategoryName: catName.GetText(),
// 					Sessions:     []entity.Session{},
// 				}
// 				*c.bookmarks = append(*c.bookmarks, newCategory)

// 				list.AddItem("📁 "+newCategory.CategoryName, "", 0, func() {
// 					onSelectedCategory(wnd, &newCategory)
// 				})
// 				// Remove the window and restore focus to previous window
// 				c.ui.WinMan.RemoveWindow(mdlNewCategory)
// 				c.ui.SetFocus(wnd)
// 			}
// 			return event
// 		})

// 		mdlNewCategory.SetRect(0, 0, 50, 1)
// 		mdlNewCategory.AddButton(&winman.Button{
// 			Symbol: 'X',
// 			OnClick: func() {
// 				c.ui.WinMan.RemoveWindow(mdlNewCategory)
// 				c.ui.SetFocus(wnd)
// 			},
// 		})

// 		c.ui.WinMan.AddWindow(mdlNewCategory)
// 		c.ui.WinMan.Center(mdlNewCategory)
// 		c.ui.SetFocus(mdlNewCategory)
// 	}

// 	// If user wants to create a new category
// 	list.AddItem("📖 Create New Category", "", 0, createCategoryModal)

// 	// Populate the list
// 	for i := range *c.bookmarks {
// 		b := &((*c.bookmarks)[i])
// 		list.AddItem("📁 "+b.CategoryName, "", 0, func() {
// 			onSelectedCategory(wnd, b)
// 		})
// 	}

// 	wnd.SetRect(0, 0, 50, 10)
// 	wnd.AddButton(&winman.Button{
// 		Symbol: 'X',
// 		OnClick: func() {
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		},
// 	})

// 	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		switch event.Key() {
// 		case tcell.KeyEscape:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)

// 		}
// 		return event
// 	})

// 	c.ui.WinMan.AddWindow(wnd)
// 	c.ui.WinMan.Center(wnd)
// 	c.ui.SetFocus(wnd)

// }

// // ShowBookmarkNameModal is used to show modal with text box to change the bookmark name
// func (c Controller) ShowBookmarkNameModal(parentWND winman.Window, onEnter func(bookmarkName string)) {
// 	// Create Window
// 	bookmarkName := tview.NewInputField().SetText(c.conn.Name)
// 	wnd := winman.NewWindow().
// 		Show().
// 		SetRoot(bookmarkName).
// 		SetDraggable(true).
// 		SetTitle(" ✏️ Enter Bookmark Name ")

// 	wnd.SetBackgroundColor(c.theme.Colors.WindowColor)
// 	bookmarkName.SetFieldBackgroundColor(wnd.GetBackgroundColor())

// 	bookmarkName.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		switch event.Key() {
// 		case tcell.KeyEscape:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(parentWND)
// 		case tcell.KeyEnter:
// 			onEnter(bookmarkName.GetText())

// 			// Remove the window and restore focus to menu list
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		}
// 		return event
// 	})

// 	wnd.SetRect(0, 0, 80, 1)
// 	wnd.AddButton(&winman.Button{
// 		Symbol: 'X',
// 		OnClick: func() {
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(parentWND)
// 		},
// 	})

// 	c.ui.WinMan.AddWindow(wnd)
// 	c.ui.WinMan.Center(wnd)
// 	c.ui.SetFocus(wnd)
// }

// // DoSaveBookmark used to open the save bookmark dialog to save the current payload to bookmark
// func (c Controller) DoSaveBookmark() {
// 	// Initiate saving bookmark sequence
// 	c.ShowBookmarkCategoryModal(func(wnd winman.Window, b *entity.Bookmark) {
// 		c.ShowBookmarkNameModal(wnd, func(bookmarkName string) {
// 			genID := uuid.New()
// 			c.conn.ID = &genID
// 			c.conn.Name = bookmarkName
// 			b.Sessions = append(b.Sessions, *c.conn)

// 			c.SaveBookmark()
// 		})

// 		c.ui.WinMan.RemoveWindow(wnd)
// 	})
// }

// func (c Controller) ShowBookmarkOptionsModal(parentWnd tview.Primitive, b entity.Session) {
// 	listOptions := tview.NewList().
// 		ShowSecondaryText(false)

// 	wnd := winman.NewWindow().
// 		Show().
// 		SetRoot(listOptions).
// 		SetDraggable(true).
// 		SetResizable(false).
// 		SetTitle(" 📚 Bookmark Options ")

// 	wnd.SetBackgroundColor(c.theme.Colors.WindowColor)
// 	listOptions.SetBackgroundColor(wnd.GetBackgroundColor())

// 	// Load bookmark to the current session
// 	listOptions.AddItem("Load Bookmark", "", 'a', func() {
// 		c.applyBookmark(b)

// 		// Close the window
// 		c.ui.WinMan.RemoveWindow(wnd)
// 		c.ui.SetFocus(parentWnd)
// 	})

// 	// Overwrite bookmark from the current session
// 	listOptions.AddItem("Overwrite Bookmark", "", 'o', func() {
// 		// c.overwriteBookmark(index)
// 	})

// 	listOptions.AddItem("Delete Bookmark", "", 'd', func() {
// 		c.deleteBookmark(b)

// 		// Close the window
// 		c.ui.WinMan.RemoveWindow(wnd)
// 		c.ui.SetFocus(parentWnd)
// 	})

// 	listOptions.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		switch event.Key() {
// 		case tcell.KeyEscape:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(parentWnd)
// 		case tcell.KeyEnter:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(parentWnd)
// 		}
// 		return event
// 	})

// 	wnd.SetModal(true)
// 	wnd.SetRect(0, 0, 50, 7)
// 	wnd.AddButton(&winman.Button{
// 		Symbol: 'X',
// 		OnClick: func() {
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(parentWnd)
// 		},
// 	})

// 	c.ui.WinMan.AddWindow(wnd)
// 	c.ui.WinMan.Center(wnd)
// 	c.ui.SetFocus(wnd)
// }

// func (c Controller) deleteBookmark(b entity.Session) {
// 	// Helper function to remove a session from a bookmark
// 	removeSession := func(sessions []entity.Session, sessionID *uuid.UUID) []entity.Session {
// 		for i, session := range sessions {
// 			if session.ID == sessionID {
// 				return append(sessions[:i], sessions[i+1:]...)
// 			}
// 		}
// 		return sessions
// 	}

// 	for i, bookmark := range *c.bookmarks {
// 		updatedSessions := removeSession(bookmark.Sessions, b.ID)
// 		if len(updatedSessions) != len(bookmark.Sessions) {
// 			bookmark.Sessions = updatedSessions
// 			if len(bookmark.Sessions) == 0 {
// 				*c.bookmarks = append((*c.bookmarks)[:i], (*c.bookmarks)[i+1:]...)
// 			}
// 			c.SaveBookmark()
// 			c.PrintLog("🗑️ Bookmark deleted", LOG_INFO)
// 			return
// 		}
// 	}
// }
