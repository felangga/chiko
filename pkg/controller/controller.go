package controller

import (
	"chiko/pkg/entity"
	"chiko/pkg/ui"
	"context"
	"fmt"
)

type Controller struct {
	ctx       context.Context
	ui        ui.View
	conn      *entity.Session
	bookmarks *[]entity.Session
	theme     entity.Theme
}

func NewController() Controller {

	ui := ui.NewView()
	conn := entity.Session{
		// Default server URL
		ServerURL: "localhost:50051",
	}
	init := &[]entity.Session{}
	c := Controller{
		context.Background(),
		ui,
		&conn,
		init,
		entity.SelectedTheme,
	}

	return c
}

func (c Controller) initSys() {
	c.PrintLog(fmt.Sprintf("✨ Welcome to Chiko v%s", entity.APP_VERSION), LOG_INFO)

	// Load bookmarks
	c.loadBookmarks()
}

func (c Controller) Run() error {
	c.InitMenu()

	c.initSys()

	c.ui.App.EnableMouse(true)
	return c.ui.App.Run()
}
