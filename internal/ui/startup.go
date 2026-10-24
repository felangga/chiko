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

	"github.com/felangga/chiko/internal/controller/vimmode"
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
	if u.VimEnabled {
		u.registerDefaultVimBindings()
	}
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
		if u.VimEnabled && u.VimMode != nil {
			event = u.VimMode.TranslateInput(event)
			if event == nil {
				return nil
			}
		}

		// Only fire global shortcuts when no modal windows are open.
		// WindowCount == 1 means only the main app window exists.
		if u.WinMan.WindowCount() > 1 {
			return event
		}

		// Don't intercept shortcuts while typing (InputField, or vim Edit mode).
		if _, ok := u.App.GetFocus().(*tview.InputField); ok {
			return event
		}
		if u.VimEnabled && u.VimMode != nil && u.VimMode.ActiveMode == entity.VimModeEdit {
			return event
		}

		switch r := event.Rune(); r {
		case 'u':
			u.ShowSetServerURLModal()
		case 'm':
			u.ShowSetRequestMethodModal()
		case 'd':
			u.ShowMetadataModal()
		case 'p':
			u.ShowRequestPayloadModal()
		case 'q':
			u.QuitApplication()
		default:
			switch {
			case r == u.keyAuthorization():
				u.ShowAuthorizationModal()
			case r == u.keyInvoke():
				u.InvokeRPC()
			case r == u.keyBookmark():
				u.ShowSaveToBookmarkModal()
			case r == u.keyHistory():
				u.ShowHistoryModal()
			default:
				return event
			}
		}

		return nil
	})
}

func (u *UI) registerDefaultVimBindings() {
	// nvim-like: Normal moves without selecting; Visual extends selection;
	// Edit leaves motions unbound so keys insert as text.
	bindMotion := func(km vimmode.KeyMap, withSelection bool) {
		mod := tcell.ModNone
		wordMod := tcell.ModCtrl
		// TextArea's Alt-f / Alt-b are native word jumps; Shift keeps selection.
		altWordMod := tcell.ModAlt
		if withSelection {
			mod = tcell.ModShift
			wordMod = tcell.ModCtrl | tcell.ModShift
			altWordMod = tcell.ModAlt | tcell.ModShift
		}

		// Character / line motions.
		km.Bind(
			tcell.NewEventKey(tcell.KeyRune, 'h', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyLeft, 0, mod),
		)
		km.Bind(
			tcell.NewEventKey(tcell.KeyRune, 'j', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyDown, 0, mod),
		)
		km.Bind(
			tcell.NewEventKey(tcell.KeyRune, 'k', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyUp, 0, mod),
		)
		km.Bind(
			tcell.NewEventKey(tcell.KeyRune, 'l', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyRight, 0, mod),
		)

		// Word motions (nvim w/b/e).
		// w / Ctrl-Right — forward by one word
		// b / Ctrl-Left  — back by one word
		// e / Alt-f      — end of current/next word (TextArea native)
		km.Bind(
			tcell.NewEventKey(tcell.KeyRune, 'w', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyRight, 0, wordMod),
		)
		km.Bind(
			tcell.NewEventKey(tcell.KeyRune, 'b', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyLeft, 0, wordMod),
		)
		km.Bind(
			tcell.NewEventKey(tcell.KeyRune, 'e', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyRune, 'f', altWordMod),
		)

		// Ctrl-hjkl are handled in VimMode.translateCtrlHJKL so KeyCtrlL
		// can never reach TextArea's select-all on key-repeat.
	}

	bindMotion(vimmode.NormalKeyMap, false)
	bindMotion(vimmode.VisualKeyMap, true)
	// EditKeyMap intentionally empty: motions insert as characters; Esc → Normal.
}
