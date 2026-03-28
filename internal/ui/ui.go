/*
 * Copyright (c) PT Pintu Kemana Saja 2026 All Rights Reserved.
 */

package ui

import (
	"fmt"
	"sync"

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
	MenuList     *tview.List
	BookmarkList *tview.TreeView
	LogList      *tview.TextView
	OutputPanel  InitOutputPanelComponents
	HistoryPanel *tview.TreeView
}

type UI struct {
	App    *tview.Application
	WinMan *winman.Manager
	Layout *ComponentLayout

	GRPC          *grpc.GRPC
	Bookmark      *bookmark.Bookmark
	History       *history.History
	Storage       *storage.Storage
	LogChannel    chan entity.Log
	OutputChannel chan entity.Output

	Theme *entity.Theme
}

func (u *UI) SetFocus(p tview.Primitive) {
	go u.App.QueueUpdateDraw(func() {
		u.App.SetFocus(p)
	})
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
	grpc := grpc.NewGRPC(logger, &session)
	bookmark := bookmark.NewBookmark()
	history := history.NewHistory()
	storage := storage.NewStorage()

	instance := &UI{
		App:           app,
		WinMan:        wm,
		GRPC:          &grpc,
		Bookmark:      &bookmark,
		History:       &history,
		Storage:       &storage,
		LogChannel:    logger.LogChannel(),
		OutputChannel: logger.OutputChannel(),
		Theme:         &entity.TerminalTheme,
	}

	instance.Layout = &ComponentLayout{
		MenuList:     instance.InitSidebarMenu(),
		BookmarkList: instance.InitBookmarkMenu(),
		LogList:      instance.InitLogList(),
		OutputPanel:  instance.InitOutputPanel(),
		HistoryPanel: instance.InitHistoryPanel(),
	}

	// Background window (lowest Z) — containing the main app layout
	bgWindow := wm.NewWindow().
		Show().
		SetRoot(instance.setupAppLayout()).
		SetBorder(false)
	bgWindow.Maximize()

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

// setupAppTitle sets up the header title of the application.
func setupAppTitle() *tview.TextView {
	title := tview.NewTextView()
	title.SetBorder(true)
	title.SetText(fmt.Sprintf("Chiko %s", entity.APP_VERSION))
	title.SetTextAlign(tview.AlignCenter)

	return title
}

func (u *UI) setupAppLayout() *tview.Flex {
	sidebar := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(u.Layout.MenuList, 0, 1, true).
		AddItem(u.Layout.BookmarkList, 0, 1, false)

	mainContent := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(u.Layout.OutputPanel.Layout, 0, 3, true).
		AddItem(u.Layout.LogList, 10, 1, false)

	content := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(sidebar, 35, 1, true).
		AddItem(mainContent, 0, 1, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(setupAppTitle(), 3, 1, false).
		AddItem(content, 0, 1, true)

	return layout
}
