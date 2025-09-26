package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/controller/bookmark"
	"github.com/felangga/chiko/internal/controller/grpc"
	"github.com/felangga/chiko/internal/controller/storage"
	"github.com/felangga/chiko/internal/entity"
	"github.com/felangga/chiko/internal/logger"
)

type ComponentLayout struct {
	MenuList     *tview.List
	BookmarkList *tview.TreeView
	LogList      *tview.TextView
	OutputPanel  InitOutputPanelComponents
}

type UI struct {
	App    *tview.Application
	WinMan *winman.Manager
	Layout *ComponentLayout

	GRPC          *grpc.GRPC
	Bookmark      *bookmark.Bookmark
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

// setupGlobalKeyboardShortcuts sets up global keyboard shortcuts that work from any component
func (u *UI) setupGlobalKeyboardShortcuts() {
	u.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Check if user is currently typing in a text input field
		if u.isTypingInTextField() {
			return event
		}

		// Only handle global shortcuts for single key presses
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'u':
				u.ShowSetServerURLModal()
				return nil
			case 'm':
				u.ShowSetRequestMethodModal()
				return nil
			case 'a':
				u.ShowAuthorizationModal()
				return nil
			case 'd':
				u.ShowMetadataModal()
				return nil
			case 'p':
				u.ShowRequestPayloadModal()
				return nil
			case 'i':
				u.InvokeRPC()
				return nil
			case 'b':
				u.ShowSaveToBookmarkModal()
				return nil
			case 'q':
				u.QuitApplication()
				return nil
			}
		}

		return event
	})
}

// isTypingInTextField checks if the currently focused component is a text input field
func (u *UI) isTypingInTextField() bool {
	focused := u.App.GetFocus()
	if focused == nil {
		return false
	}

	// Check if the focused component is a text input field
	switch focused.(type) {
	case *tview.InputField:
		return true
	case *tview.TextArea:
		return true
	default:
		return false
	}
}

func NewUI(session entity.Session) UI {
	logger := logger.New()

	app := tview.NewApplication()
	wm := winman.NewWindowManager()
	grpc := grpc.NewGRPC(logger, &session)
	bookmark := bookmark.NewBookmark()
	storage := storage.NewStorage()

	ui := UI{
		App:           app,
		WinMan:        wm,
		GRPC:          &grpc,
		Bookmark:      &bookmark,
		Storage:       &storage,
		LogChannel:    logger.LogChannel(),
		OutputChannel: logger.OutputChannel(),
		Theme:         &entity.TerminalTheme,
	}

	ui.Layout = &ComponentLayout{
		MenuList:     ui.InitSidebarMenu(),
		BookmarkList: ui.InitBookmarkMenu(),
		LogList:      ui.InitLogList(),
		OutputPanel:  ui.InitOutputPanel(),
	}

	window := wm.NewWindow().
		Show().
		SetRoot(ui.setupAppLayout()).
		SetBorder(false)

	window.Maximize()
	app.SetRoot(wm, true)

	// Set up global keyboard shortcuts
	ui.setupGlobalKeyboardShortcuts()

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

	splitMainPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(u.Layout.OutputPanel.Layout, 0, 3, false).
		AddItem(u.Layout.LogList, 0, 1, false)

	childLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(splitSidebar, 35, 1, true).
		AddItem(splitMainPanel, 0, 4, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(setupAppTitle(), 3, 1, false).
		AddItem(childLayout, 0, 1, true)

	return layout
}
