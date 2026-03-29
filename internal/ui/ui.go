/*
 * Copyright (c) PT Pintu Kemana Saja 2026 All Rights Reserved.
 */

package ui

import (
	"sync"

	"github.com/google/uuid"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/controller/bookmark"
	"github.com/felangga/chiko/internal/controller/grpc"
	"github.com/felangga/chiko/internal/controller/history"
	"github.com/felangga/chiko/internal/controller/storage"
	"github.com/felangga/chiko/internal/entity"
	"github.com/felangga/chiko/internal/logger"
)

type ComponentLayout struct {
	BookmarkList *tview.TreeView
	LogList      *tview.TextView
	HistoryPanel *tview.TreeView
}

// SessionWindow represents a discrete, floating or maximized winman.Window containing a Request context
type SessionWindow struct {
	ID          uuid.UUID
	GRPC        *grpc.GRPC
	WinBase     *winman.WindowBase
	MenuList    *tview.List // kept for modal fallback focus target
	OutputPanel InitOutputPanelComponents

	// Postman-style inline UI components
	URLField        *tview.InputField
	MethodField     *tview.InputField
	SendBtn         *tview.Button
	RequestBodyArea *tview.TextArea
	AuthInput       *tview.InputField
	MetadataArea    *tview.TextArea
	RequestPages   *tview.Pages
	RequestTabBar  *tview.TextView
	RefreshTopBar  func()
	RefreshRequestTabs func(string)
}

type UI struct {
	App    *tview.Application
	WinMan *winman.Manager
	Layout *ComponentLayout

	Sessions       []*SessionWindow
	ActiveSession  *SessionWindow

	Logger        *logger.Logger
	GRPC          *grpc.GRPC
	Bookmark      *bookmark.Bookmark
	History       *history.History
	Storage       *storage.Storage
	LogChannel    chan entity.Log
	OutputChannel chan entity.Output

	Theme *entity.Theme

	// Dashboard window reference for Z-index manipulation
	HomeWindow *winman.WindowBase
}

func (u *UI) SetFocus(p tview.Primitive) {
	go u.App.QueueUpdateDraw(func() {
		u.App.SetFocus(p)
	})
}

// activeSessionFocus returns the best focus target for the current session window.
// In the Postman-style layout the URLField is the primary interactive element.
func (u *UI) activeSessionFocus() tview.Primitive {
	if u.ActiveSession != nil {
		if u.ActiveSession.URLField != nil {
			return u.ActiveSession.URLField
		}
		if u.ActiveSession.MenuList != nil {
			return u.ActiveSession.MenuList
		}
	}
	return u.Layout.LogList
}

func (u UI) Run() error {
	u.App.EnableMouse(true)
	return u.App.Run()
}

func (u UI) QuitApplication() {
	u.App.Stop()
}

func NewUI(session entity.Session) UI {
	logger := logger.New()

	app := tview.NewApplication()
	wm := winman.NewWindowManager()
	bookmark := bookmark.NewBookmark()
	history := history.NewHistory()
	storage := storage.NewStorage()

	instance := &UI{
		App:           app,
		WinMan:        wm,
		Logger:        logger,
		Bookmark:      &bookmark,
		History:       &history,
		Storage:       &storage,
		LogChannel:    logger.LogChannel(),
		OutputChannel: logger.OutputChannel(),
		Theme:         &entity.TerminalTheme,
	}

	instance.Layout = &ComponentLayout{
		BookmarkList: instance.InitBookmarkMenu(),
		LogList:      instance.InitLogList(),
		HistoryPanel: instance.InitHistoryPanel(),
	}

	// Create the Floating Home Dashboard Window
	homeMenuLayout := instance.InitHomeMenu()
	bgWindow := wm.NewWindow().
		Show().
		SetRoot(homeMenuLayout).
		SetBorder(true).
		SetDraggable(true).
		SetResizable(true).
		SetTitle(" 🚀 Chiko Dashboard ")
	
	// Default starting size for the main menu
	bgWindow.SetRect(2, 2, 60, 20)
	instance.HomeWindow = bgWindow

	// If the user launched Chiko with CLI arguments, start an RPC session automatically
	if session.ID != uuid.Nil {
		instance.CreateSessionWindow(&session)
	}

	app.SetRoot(wm, true)

	// Trigger startupSequence only after the first draw so the
	// event loop is guaranteed to be running and can drain QueueUpdateDraw.
	var once sync.Once
	app.SetAfterDrawFunc(func(screen tcell.Screen) {
		once.Do(func() {
			app.SetAfterDrawFunc(nil) // unregister immediately
			go instance.startupSequence()
		})
	})

	return *instance
}


