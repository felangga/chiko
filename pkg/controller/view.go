package controller

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

func (c Controller) initMenu() {
	c.ui.MenuList.AddItem("Server URL", "", 'u', c.setServerURL)
	c.ui.MenuList.AddItem("Method", "", 'm', c.setMethod)
	c.ui.MenuList.AddItem("Metadata", "", 'd', nil)
}

func (c Controller) initSys() {
	c.PrintLog("✨ Chiko ", LOG_INFO)
}

func (c Controller) setServerURL() {
	tmpURL := c.conn.ServerURL

	// Create Set Server URL From
	form := tview.NewForm()
	wnd := winman.NewWindow().
		Show().
		SetRoot(form).
		SetDraggable(true)

	form.AddInputField("Server URL", c.conn.ServerURL, 0, nil, func(txt string) {
		tmpURL = txt
	})

	form.AddButton("Set", func() {
		c.conn.ServerURL = tmpURL
		// Remove the window and restore focus to menu list
		c.PrintLog("Server URL set to [blue]"+c.conn.ServerURL, LOG_INFO)
		c.ui.WinMan.RemoveWindow(wnd)
		c.ui.SetFocus(c.ui.MenuList)

		c.CheckGRPC()
	})

	form.AddButton("Cancel", func() {
		// Remove the window and restore focus to menu list
		c.ui.WinMan.RemoveWindow(wnd)
		c.ui.SetFocus(c.ui.MenuList)
	})
	form.SetButtonsAlign(tview.AlignRight)

	wnd.SetModal(true)
	wnd.SetRect(0, 0, 50, 7)

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)
}

func (c Controller) setMethod() {
	selectedMethod := c.conn.SelectedMethod

	// Create Set Server URL From
	form := tview.NewForm()
	wnd := winman.NewWindow().
		Show().
		SetRoot(form).
		SetDraggable(true)

	form.AddDropDown("Select Method", c.conn.AvailableMethods, 0, func(option string, optionIndex int) {
		selectedMethod = option
	}).SetBorderPadding(1, 1, 1, 1)

	form.AddButton("Set", func() {
		c.conn.SelectedMethod = selectedMethod

		// Remove the window and restore focus to menu list
		c.PrintLog("👉 Method set to [blue]"+c.conn.SelectedMethod, LOG_INFO)
		c.ui.WinMan.RemoveWindow(wnd)
		c.ui.SetFocus(c.ui.MenuList)
	})

	form.AddButton("Cancel", func() {
		// Remove the window and restore focus to menu list
		c.ui.WinMan.RemoveWindow(wnd)
		c.ui.SetFocus(c.ui.MenuList)
	})
	form.SetButtonsAlign(tview.AlignRight)

	wnd.SetModal(true)
	wnd.SetRect(0, 0, 70, 7)

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)
}
