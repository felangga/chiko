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
	u.loadWorkspace()
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
	if u.GRPC == nil || u.GRPC.Conn.ID == uuid.Nil {
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

// loadWorkspace restores previous session windows
func (u *UI) loadWorkspace() {
	ws, err := u.Workspace.LoadWorkspace()
	if err != nil {
		u.PrintLog(entity.Log{
			Content: fmt.Sprintf("⚠️ could not load workspace: %v", err),
			Type:    entity.LOG_ERROR,
		})
		return
	}

	if ws == nil || len(ws.ActiveSessions) == 0 {
		return
	}

	u.PrintLog(entity.Log{
		Content: fmt.Sprintf("📂 restoring %d session(s) from previous workspace", len(ws.ActiveSessions)),
		Type:    entity.LOG_INFO,
	})

	var activeSession *SessionWindow
	for _, session := range ws.ActiveSessions {
		u.CreateSessionWindow(session)
		if session.ID == ws.ActiveSessionID {
			activeSession = u.Sessions[len(u.Sessions)-1]
		}
	}

	if activeSession != nil {
		u.ActiveSession = activeSession
		u.GRPC = activeSession.GRPC
		u.SetFocus(u.activeSessionFocus())
	}
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
		// Modal windows are any windows that exceed our known Dashboard (1) + Sessions
		if u.WinMan.WindowCount() > len(u.Sessions)+1 {
			return event
		}

		switch event.Key() {
		case tcell.KeyCtrlN:
			u.CreateSessionWindow(nil)
			return nil
		case tcell.KeyCtrlW:
			if u.ActiveSession != nil {
				u.CloseSession(u.ActiveSession)
			}
			return nil
		case tcell.KeyCtrlL: // next tab
			u.FocusNextSession()
			return nil
		case tcell.KeyCtrlH: // prev tab
			u.FocusPrevSession()
			return nil
		case tcell.KeyTab: // cycle forward
			if u.ActiveSession != nil {
				// Don't intercept tabs natively when inside the payload body JSON editor
				if u.App.GetFocus() != u.ActiveSession.RequestBodyArea {
					u.ActiveSession.CycleFocus(u, 1)
					return nil
				}
			}
		case tcell.KeyBacktab: // cycle backward
			if u.ActiveSession != nil {
				u.ActiveSession.CycleFocus(u, -1)
				return nil
			}
		}

		// Don't intercept single-key shortcuts while the user is typing in an input field or text area
		_, isInput := u.App.GetFocus().(*tview.InputField)
		_, isTextArea := u.App.GetFocus().(*tview.TextArea)
		if isInput || isTextArea {
			return event
		}

		switch event.Rune() {
		case 'u':
			u.ShowSetServerURLModal()
		case 'm':
			if u.ActiveSession != nil && u.ActiveSession.MethodField != nil {
				u.SetFocus(u.ActiveSession.MethodField)
			}
		case 'a':
			if u.ActiveSession != nil && u.ActiveSession.AuthInput != nil {
				u.ActiveSession.RequestPages.SwitchToPage(reqPanelAuth)
				u.ActiveSession.RefreshRequestTabs(reqPanelAuth)
				u.SetFocus(u.ActiveSession.AuthInput)
			}
		case 'd':
			if u.ActiveSession != nil && u.ActiveSession.MetadataArea != nil {
				u.ActiveSession.RequestPages.SwitchToPage(reqPanelMetadata)
				u.ActiveSession.RefreshRequestTabs(reqPanelMetadata)
				u.SetFocus(u.ActiveSession.MetadataArea)
			}
		case 'p':
			if u.ActiveSession != nil && u.ActiveSession.RequestBodyArea != nil {
				u.SetFocus(u.ActiveSession.RequestBodyArea)
			}
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
