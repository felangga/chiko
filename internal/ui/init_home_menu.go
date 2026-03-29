package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InitHomeMenu initializes the starting dashboard menu
func (u *UI) InitHomeMenu() tview.Primitive {
	menuList := tview.NewList().ShowSecondaryText(true)
	// Padding ensures the text isn't stuck right against the window border
	menuList.SetBorderPadding(1, 1, 2, 2)

	menuList.AddItem("Return to Workspace", "Return to your active request tabs", 'r', func() {
		if len(u.Sessions) > 0 {
			u.Sessions[0].WinBase.Show()
		} else {
			// Show an error/modal or create one implicitly
			u.CreateSessionWindow(nil)
		}
	})

	menuList.AddItem("Open New Request", "Start a new RPC request session", 'n', func() {
		u.CreateSessionWindow(nil)
	})

	menuList.AddItem("Bookmarks", "Manage saved requests", 'b', u.ShowBookmarkManager)
	menuList.AddItem("History", "Review past request history", 'h', u.ShowHistoryManager)
	menuList.AddItem("Settings", "Application settings (Coming Soon)", 's', nil)
	menuList.AddItem("Quit", "Exit the application", 'q', u.QuitApplication)

    // Ensure keyboard shortcuts for the Home menu work
	menuList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// Do nothing or maybe show prompt to quit
		}
		return event
	})

	return menuList
}

// Stubs for showing the global History/Bookmarks manually from the Dashboard
func (u *UI) ShowBookmarkManager() {
	u.ShowSaveToBookmarkModal() // TODO: Actually make a global bookmark browser modal
}

func (u *UI) ShowHistoryManager() {
	u.ShowHistoryModal() // TODO: Create a dashboard version of History
}
