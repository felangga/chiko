package controller

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"

	"chiko/pkg/entity"
)

func (c Controller) applyBookmark(idx int) {

	// Get selected connection
	selConn := *(c.bookmarks)
	*c.conn = selConn[idx]

	c.PrintLog(fmt.Sprintf("📚 bookmark loaded : %s", c.conn.ServerURL), LOG_INFO)
	go c.CheckGRPC(c.conn.ServerURL)
	c.PrintLog(c.conn.RequestPayload, LOG_INFO)
}

// loadBookmarks is used to load bookmarks from bookmark file
func (c Controller) loadBookmarks() {
	if _, err := os.Stat(entity.BOOKMARKS_FILE_NAME); err != nil {
		c.PrintLog(err.Error(), LOG_INFO)
		return
	}

	// Read bookmark file
	file, err := os.ReadFile(entity.BOOKMARKS_FILE_NAME)
	if err != nil {
		c.PrintLog(err.Error(), LOG_ERROR)
		return
	}

	err = json.Unmarshal([]byte(file), &c.bookmarks)
	if err != nil {
		c.PrintLog(err.Error(), LOG_ERROR)
		return
	}

	// Populate bookmarks list
	for i, b := range *c.bookmarks {
		c.ui.BookmarkList.AddItem("📗 "+b.Name, "", 0, func() {
			c.applyBookmark(i)
		})
	}

	c.PrintLog(fmt.Sprintf("📚 %d bookmark(s) loaded", len(*c.bookmarks)), LOG_INFO)

}

// saveBookmark is used to save the bookmark to file
func (c Controller) saveBookmark(conn entity.Session) {

	lastIdx := c.ui.BookmarkList.GetItemCount()
	*c.bookmarks = append(*c.bookmarks, conn)
	c.ui.BookmarkList.AddItem("📗 "+conn.Name, "", 0, func() {
		c.applyBookmark(lastIdx)
	})

	// Save to file
	convert, err := json.Marshal(c.bookmarks)
	if err != nil {
		c.PrintLog("❌ failed to prepare bookmarks data", LOG_ERROR)
		return
	}

	err = os.WriteFile("bookmarks.cfg", convert, 0644)
	if err != nil {
		c.PrintLog("💾 failed to write bookmark configuration, maybe write-protected?", LOG_ERROR)
		return
	}
}

func (c Controller) showNewBookmarkModal() {
	genID := uuid.New()
	// Create Window
	bookmarkName := tview.NewInputField().SetText(c.conn.Name)
	wnd := winman.NewWindow().
		Show().
		SetRoot(bookmarkName).
		SetDraggable(true).
		SetTitle(" 📚 Enter Bookmark Name ")

	wnd.SetBackgroundColor(tcell.GetColor(entity.SelectedTheme.Colors["WindowColor"]))
	bookmarkName.SetFieldBackgroundColor(wnd.GetBackgroundColor())

	bookmarkName.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		case tcell.KeyEnter:
			c.conn.ID = &genID
			c.conn.Name = bookmarkName.GetText()
			c.saveBookmark(*c.conn)

			// Remove the window and restore focus to menu list
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		}
		return event
	})

	wnd.SetRect(0, 0, 50, 1)
	wnd.AddButton(&winman.Button{
		Symbol: 'X',
		OnClick: func() {
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		},
	})

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)
}

// DoSaveBookmark used to open the save bookmark dialog to save the current payload to bookmark
func (c Controller) DoSaveBookmark() {
	conID := c.conn.ID

	if conID == nil {
		c.showNewBookmarkModal()
	}

	// // If bookmark already exist?
	// c.ShowMessageBox(" Overwrite Bookmark ", "Do you want to create new bookmark or overwrite existing bookmark?", []Button{
	// 	{"Overwrite", nil},
	// 	{"Create New", nil},
	// })
}
