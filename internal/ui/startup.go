package ui

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

func (u *UI) startupSequence() {
	u.loadStartupUI()
	u.loadBookmarks()
	u.startLogDumper()
	u.startArgsConnection()
}

// loadStartupUI displays the welcome message and banner
func (u *UI) loadStartupUI() {
	u.PrintLog(entity.Log{
		Content: fmt.Sprintf("✨ Welcome to Chiko v%s", entity.APP_VERSION),
		Type:    entity.LOG_INFO,
	})

	banner, _ := base64.StdEncoding.DecodeString(entity.BANNER)
	u.PrintOutput(entity.Output{
		Content:     string(banner),
		WithHeader:  false,
		CursorAtEnd: false,
	})
}

// startArgsConnection handle the connection request from command line arguments
func (u *UI) startArgsConnection() {
	if u.GRPC.Conn.ID == uuid.Nil {
		return
	}

	go func() {
		err := u.GRPC.Connect()
		if err != nil {
			u.PrintLog(entity.Log{
				Content: err.Error(),
				Type:    entity.LOG_ERROR,
			})
			return
		}
	}()
}

// loadBookmarks loads the bookmarks from file
func (u *UI) loadBookmarks() {
	// Load bookmarks from file
	err := u.Bookmark.LoadBookmarks()
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, create a new one
			err = u.Bookmark.SaveBookmark()
			if err != nil {
				u.PrintLog(entity.Log{
					Content: fmt.Sprintf("❌ failed to create new bookmark file, err: %v", err),
					Type:    entity.LOG_ERROR,
				})
				return
			}

			u.PrintLog(entity.Log{
				Content: "📁 new bookmark file created",
				Type:    entity.LOG_INFO,
			})
			return
		}

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
func (u *UI) startLogDumper() {
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
