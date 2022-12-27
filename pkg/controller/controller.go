package controller

import (
	"chiko/pkg/ui"
)

type Controller struct {
	ui ui.View
}

func NewController() Controller {

	ui := ui.NewView()

	// ui.Frame.AddText("Chiko - your gRPC client companion", true, tview.AlignCenter, tcell.ColorWhite)

	c := Controller{
		ui,
	}

	return c
}

func (c Controller) Run() error {

	return c.ui.App.Run()
}
