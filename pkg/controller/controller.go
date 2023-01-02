package controller

import (
	"chiko/pkg/ui"
)

type Controller struct {
	ui ui.View
}

func NewController() Controller {
	ui := ui.NewView()
	c := Controller{
		ui,
	}

	return c
}

func (c Controller) setServerURL() {
	c.PrintLog("Set server URL")
}

func (c Controller) initMenu() {
	c.ui.MenuList.AddItem("Server URL", "", 'u', c.setServerURL)
	c.ui.MenuList.AddItem("Method", "", 'm', nil)
	c.ui.MenuList.AddItem("Metadata", "", 'd', nil)
}

func (c Controller) initSys() {
	c.ui.OutputPanel.SetDynamicColors(true)
	c.ui.OutputPanel.SetText("[blue]Welcome to [white]Chiko v.0.0.1")
}

func (c Controller) Run() error {
	c.initMenu()

	c.initSys()

	c.ui.App.EnableMouse(true)
	return c.ui.App.Run()
}
