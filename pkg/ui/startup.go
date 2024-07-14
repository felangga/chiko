package ui

import (
	"chiko/pkg/entity"
	"fmt"

	"github.com/rivo/tview"
)

func (u *UI) startupSequence() {
	u.PrintLog(entity.Log{
		Content: fmt.Sprintf("✨ Welcome to Chiko v%s", entity.APP_VERSION),
		Type:    entity.LOG_INFO,
	})

	u.loadBookmarks()
	u.logDumper()
}

func (u *UI) loadBookmarks() {
	// Load bookmarks from file
	err := u.Bookmark.LoadBookmarks()
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
			}
		}
	}()
}
