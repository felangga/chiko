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
	conn      *entity.Connection
	bookmarks *[]entity.Connection
}

func NewController() Controller {

	ui := ui.NewView()
	conn := entity.Connection{
		ServerURL: "localhost:50051",
	}
	init := &[]entity.Connection{}
	c := Controller{
		context.Background(),
		ui,
		&conn,
		init,
	}

	return c
}

func (c Controller) initSys() {
	c.PrintLog(fmt.Sprintf("✨ Welcome to Chiko v%s", entity.APP_VERSION), LOG_INFO)
	// Load bookmarks
	c.loadBookmark()
}

func (c Controller) Run() error {
	c.initMenu()

	c.initSys()

	c.ui.App.EnableMouse(true)
	return c.ui.App.Run()
}
