package ui

import (
	"chiko/pkg/entity"
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type View struct {
	App          *tview.Application
	Pages        *tview.Pages
	MainLayout   *tview.Flex
	MenuList     *tview.List
	BookmarkList *tview.TreeView
	OutputPanel  *tview.TextView
	WinMan       *winman.Manager
}

func (v View) SetFocus(p tview.Primitive) {
	go v.App.QueueUpdateDraw(func() {
		v.App.SetFocus(p)
	})
}

func NewView() View {
	app := tview.NewApplication()
	pages := tview.NewPages()
	wm := winman.NewWindowManager()

	title := tview.NewTextView()
	title.SetBorder(true)
	title.SetText(fmt.Sprintf("Chiko v%s", entity.APP_VERSION))
	title.SetTextAlign(tview.AlignCenter)

	// Setup the side bar menu
	menuList := tview.NewList().ShowSecondaryText(false)
	menuList.SetBorder(true).SetTitle(" 🐶 Menu ")
	menuList.SetBorderPadding(1, 1, 1, 1)

	// Initialize the bookmarks tree view
	root := tview.NewTreeNode("📚 Library")
	bookmarkList := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	bookmarkList.SetBorder(true).SetTitle(" 📚 Bookmarks Library ").SetBorderPadding(1, 1, 1, 1)

	outputPanel := tview.NewTextView()
	outputPanel.SetDynamicColors(true)
	outputPanel.SetTitle(" 📃 Output Logs ")
	outputPanel.SetBorder(true)
	outputPanel.SetWordWrap(true)
	outputPanel.SetBorderPadding(1, 1, 1, 1)
	outputPanel.SetScrollable(true).SetChangedFunc(func() {
		app.Draw()
	})

	// Handle keypress on menu list
	menuList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(bookmarkList)
		}
		return event
	})
	bookmarkList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(outputPanel)
		}
		return event
	})
	outputPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(menuList)
		}
		return event
	})

	// Setup the main layout
	splitSidebar := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(menuList, 15, 1, true).
		AddItem(bookmarkList, 0, 1, false)

	childLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(splitSidebar, 35, 1, true).
		AddItem(outputPanel, 0, 4, false)

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 3, 1, false).
		AddItem(childLayout, 0, 1, true)

	window := wm.NewWindow().
		Show().
		SetRoot(mainLayout).
		SetBorder(false)

	window.Maximize()
	app.SetRoot(wm, true)

	v := View{
		app,
		pages,
		mainLayout,
		menuList,
		bookmarkList,
		outputPanel,
		wm,
	}

	return v
}
