package controller

import (
	"chiko/pkg/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func (c Controller) applyBookmark() {
	// Get selected connection
	selIndex := c.ui.BookmarkList.GetCurrentItem()
	selConn := *(c.bookmarks)
	*c.conn = selConn[selIndex]

	c.PrintLog(fmt.Sprintf("📚 bookmark applied: %s", c.conn.ServerURL), LOG_INFO)
}

func (c Controller) loadBookmark() {
	// Get selected connection
	file, err := ioutil.ReadFile("bookmarks.cfg")
	if err != nil {
		c.PrintLog(err.Error(), LOG_INFO)
		return
	}

	err = json.Unmarshal([]byte(file), &c.bookmarks)
	if err != nil {
		c.PrintLog(err.Error(), LOG_INFO)
		return
	}

	// Populate bookmarks list
	for _, b := range *c.bookmarks {
		c.ui.BookmarkList.AddItem(b.Name, "", 0, c.applyBookmark)
	}

	c.PrintLog(fmt.Sprintf("%d 📚 bookmarks loaded", len(*c.bookmarks)), LOG_INFO)

}

func (c Controller) saveBookmark(conn entity.Connection) {

	*c.bookmarks = append(*c.bookmarks, conn)
	c.ui.BookmarkList.AddItem(conn.Name, "", 0, c.applyBookmark)

	// Save to file
	convert, err := json.Marshal(c.bookmarks)
	if err != nil {
		c.PrintLog("❌ failed to prepare bookmarks data", LOG_INFO)
		return
	}

	err = ioutil.WriteFile("bookmarks.cfg", convert, 0644)
	if err != nil {
		c.PrintLog("💾 failed to write bookmark configuration, maybe write-protected?", LOG_INFO)
		return
	}
}

func (c Controller) doSaveBookmark() {
	conID := c.conn.ID

	c.ShowMessageBox(" Overwrite Bookmark ", "Do you want to create new bookmark or overwrite existing bookmark?")
	return
	if conID == nil {
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
			Symbol: '❌',
			OnClick: func() {
				c.ui.WinMan.RemoveWindow(wnd)
				c.ui.SetFocus(c.ui.MenuList)
			},
		})

		c.ui.WinMan.AddWindow(wnd)
		c.ui.WinMan.Center(wnd)
		c.ui.SetFocus(wnd)
	}
}
