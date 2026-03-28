/*
 * Copyright (c) PT Pintu Kemana Saja 2026 All Rights Reserved.
 */

package ui

import (
	"encoding/base64"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

func (u *UI) startupSequence() {
	u.loadStartupUI()
	u.loadBookmarks()
	u.loadHistory()
	u.startLogDumper()
	go u.checkForUpdates()
	u.startArgsConnection()
	u.setupGlobalInputCapture()
}

// loadStartupUI displays the welcome message and banner
func (u *UI) loadStartupUI() {
	u.PrintLog(entity.Log{
		Content: fmt.Sprintf("✨ Welcome to Chiko %s", entity.APP_VERSION),
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

// loadHistory loads the request history from file
func (u *UI) loadHistory() {
	if err := u.History.LoadHistory(); err != nil {
		u.PrintLog(entity.Log{
			Content: fmt.Sprintf("⚠️ could not load history: %v", err),
			Type:    entity.LOG_ERROR,
		})
		return
	}

	u.PrintLog(entity.Log{
		Content: fmt.Sprintf("🕓 %d history entry(s) loaded", len(*u.History.Entries)),
		Type:    entity.LOG_INFO,
	})

	u.RefreshHistoryPanel()
}

func (u *UI) RefreshBookmarkList() int16 {
	var totalBookmarks int16
	newChildren := []*tview.TreeNode{}

	for _, b := range *u.Bookmark.Categories {
		categoryNode := tview.NewTreeNode("📁 " + b.Name)
		categoryNode.SetReference(b)

		for _, session := range b.Sessions {
			sessionNode := tview.NewTreeNode("📗 " + session.Name)
			sessionNode.SetReference(&session)
			categoryNode.AddChild(sessionNode)
			totalBookmarks++
		}
		newChildren = append(newChildren, categoryNode)
	}

	go u.App.QueueUpdateDraw(func() {
		root := u.Layout.BookmarkList.GetRoot()
		root.ClearChildren()
		for _, child := range newChildren {
			root.AddChild(child)
		}
	})

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

// setupGlobalInputCapture sets up application-wide key bindings that mirror
// the sidebar menu shortcuts, so they work regardless of which panel is focused.
// Keys are suppressed when an InputField is focused to avoid interfering with text entry.
func (u *UI) setupGlobalInputCapture() {
	u.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Only fire global shortcuts when no modal windows are open.
		// WindowCount == 1 means only the main app window exists.
		if u.WinMan.WindowCount() > 1 {
			return event
		}

		// Don't intercept shortcuts while the user is typing in an input field
		if _, ok := u.App.GetFocus().(*tview.InputField); ok {
			return event
		}

		switch event.Rune() {
		case 'u':
			u.ShowSetServerURLModal()
		case 'm':
			u.ShowSetRequestMethodModal()
		case 'a':
			u.ShowAuthorizationModal()
		case 'd':
			u.ShowMetadataModal()
		case 'p':
			u.ShowRequestPayloadModal()
		case 'i':
			u.InvokeRPC()
		case 'h':
			u.ShowHistoryModal()
		case 'b':
			u.ShowSaveToBookmarkModal()
		case 'q':
			u.QuitApplication()
		default:
			return event
		}

		return nil
	})
}
