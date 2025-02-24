package ui

import (
	"encoding/base64"
	"fmt"

	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

func (u *UI) startupSequence() {
	u.PrintLog(entity.Log{
		Content: fmt.Sprintf("✨ Welcome to Chiko v%s", entity.APP_VERSION),
		Type:    entity.LOG_INFO,
	})

	u.loadBookmarks()
	u.logDumper()

	banner, _ := base64.StdEncoding.DecodeString(entity.BANNER)
	u.PrintOutput(entity.Output{
		Content:     string(banner),
		WithHeader:  false,
		CursorAtEnd: false,
	})

}

func (u *UI) loadBookmarks() {
	// Before v.0.0.4, the default bookmark file location is at binary folder
	// Thus we need to move the bookmark file to the new location to the OS default config folder
	err := u.Bookmark.MigrateBookmark()
	if err != nil {
		u.PrintLog(entity.Log{
			Content: fmt.Sprintf("❌ failed to migrate bookmarks, err: %v", err),
			Type:    entity.LOG_ERROR,
		})
		return
	}

	// Load bookmarks from file
	err = u.Bookmark.LoadBookmarks()
	if err != nil {
		u.PrintLog(entity.Log{
			Content: fmt.Sprintf("❌ failed to load bookmarks, err: %v", err),
			Type:    entity.LOG_ERROR,
		})
		return
	}

	totalBookmarks := u.RefreshBookmarkList()

	u.PrintLog(entity.Log{
		Content: fmt.Sprintf("📚 %d bookmark(s) loaded", totalBookmarks),
		Type:    entity.LOG_INFO,
	})
}

func (u *UI) RefreshBookmarkList() int16 {
	u.Layout.BookmarkList.GetRoot().ClearChildren()

	// Populate bookmarks list
	var totalBookmarks int16
	for _, b := range *u.Bookmark.Categories {
		categoryNode := tview.NewTreeNode("📁 " + b.Name)
		categoryNode.SetReference(b)

		for _, session := range b.Sessions {
			sessionNode := tview.NewTreeNode("📗 " + session.Name)
			sessionNode.SetReference(&session)
			categoryNode.AddChild(sessionNode)

			totalBookmarks++
		}

		u.Layout.BookmarkList.GetRoot().AddChild(categoryNode)
	}

	return totalBookmarks
}

// logDumper is used to dump log messages from channels to log window
func (u *UI) logDumper() {
	go func() {
		for {
			select {
			case log := <-u.LogChannel:
				u.PrintLog(log)
			case output := <-u.OutputChannel:
				u.PrintOutput(output)
			}
		}
	}()
}
