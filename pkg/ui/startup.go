package ui

import (
	"chiko/pkg/entity"
	"fmt"

	"github.com/rivo/tview"
)

func (u *UI) startupSequence() {
	u.PrintLog(entity.LogParam{
		Content: fmt.Sprintf("✨ Welcome to Chiko v%s", entity.APP_VERSION),
		Type:    entity.LOG_INFO,
	})

	u.loadBookmarks()
	u.logDumper()
}

func (u *UI) loadBookmarks() {
	// Load bookmarks from file
	err := u.Controller.LoadBookmarks()
	if err != nil {
		u.PrintLog(entity.LogParam{
			Content: fmt.Sprintf("❌ failed to load bookmarks, err: %v", err),
			Type:    entity.LOG_ERROR,
		})
		return
	}

	// Populate bookmarks list
	for _, b := range *u.Controller.Bookmarks {
		categoryNode := tview.NewTreeNode("📁 " + b.CategoryName)
		categoryNode.SetReference(b)

		for _, session := range b.Sessions {
			sessionNode := tview.NewTreeNode("📗 " + session.Name)
			sessionNode.SetReference(session)
			categoryNode.AddChild(sessionNode)
		}

		u.Layout.BookmarkList.GetRoot().AddChild(categoryNode)
	}

	u.PrintLog(entity.LogParam{
		Content: fmt.Sprintf("📚 %d bookmark(s) loaded", len(*u.Controller.Bookmarks)),
		Type:    entity.LOG_INFO,
	})
}

// logDumper is used to dump log messages from channels to log window
func (u *UI) logDumper() {
	go func() {
		for {
			select {
			case log := <-u.Controller.LogDump:
				u.PrintLog(log)
			}
		}
	}()
}