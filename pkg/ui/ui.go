package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"

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
}

func (u *UI) SetFocus(p tview.Primitive) {
	go u.App.QueueUpdateDraw(func() {
		u.App.SetFocus(p)
	})
}

func NewUI() UI {
	app := tview.NewApplication()
	wm := winman.NewWindowManager()

	// Handle keypress on menu list

	// bookmarkList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	switch event.Key() {
	// 	case tcell.KeyTAB:
	// 		app.SetFocus(outputPanel)
	// 	}
	// 	return event
	// })
	// outputPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	switch event.Key() {
	// 	case tcell.KeyTAB:
	// 		app.SetFocus(menuList)
	// 	}
	// 	return event
	// })

	ui := UI{
		app,
		wm,
		nil,
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

func (u UI) Run() error {
	return u.App.Run()
}

func (u UI) QuitApplication() {
	u.App.Stop()
}
