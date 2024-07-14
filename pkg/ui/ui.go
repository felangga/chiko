package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"

	"chiko/pkg/controller/bookmark"
	"chiko/pkg/controller/grpc"
	"chiko/pkg/entity"
)

type ComponentLayout struct {
	MenuList     *tview.List
	BookmarkList *tview.TreeView
	LogList      *tview.TextView
}

type UI struct {
	App    *tview.Application
	WinMan *winman.Manager
	Layout *ComponentLayout

	GRPC       *grpc.GRPC
	Bookmark   *bookmark.Bookmark
	LogChannel chan entity.Log

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

func NewUI() UI {
	log := make(chan entity.Log)

	app := tview.NewApplication()
	wm := winman.NewWindowManager()
	grpc := grpc.NewGRPC(log)
	bookmark := bookmark.NewBookmark()

	ui := UI{
		app,
		wm,
		nil,
		&grpc,
		&bookmark,
		log,
		&entity.TerminalTheme,
	}

	ui.Layout = &ComponentLayout{
		MenuList:     ui.InitSidebarMenu(),
		BookmarkList: ui.InitBookmarkMenu(),
		LogList:      ui.InitLogList(),
	}

	window := wm.NewWindow().
		Show().
		SetRoot(ui.setupAppLayout()).
		SetBorder(false)

	window.Maximize()
	app.SetRoot(wm, true)

	ui.startupSequence()
	return ui
}

// setupTitle sets up the header title of the application.
// containing the application name and version.
func setupAppTitle() *tview.TextView {
	title := tview.NewTextView()
	title.SetBorder(true)
	title.SetText(fmt.Sprintf("Chiko v%s", entity.APP_VERSION))
	title.SetTextAlign(tview.AlignCenter)

	return title
}

// setupAppLayout sets up the main grid layout of the application.
func (u *UI) setupAppLayout() *tview.Flex {
	// Setup the main layout
	splitSidebar := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(u.Layout.MenuList, 15, 1, true).
		AddItem(u.Layout.BookmarkList, 0, 1, false)

	childLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(splitSidebar, 35, 1, true).
		AddItem(u.Layout.LogList, 0, 4, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(setupAppTitle(), 3, 1, false).
		AddItem(childLayout, 0, 1, true)

	return layout
}
