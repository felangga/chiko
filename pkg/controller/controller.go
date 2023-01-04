package controller

import (
	"chiko/pkg/ui"
	"context"

	"github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc"
)

type Controller struct {
	ctx  context.Context
	ui   ui.View
	conn *Connection
}

type Connection struct {
	ServerURL         string
	ActiveConnection  *grpc.ClientConn
	AvailableServices []string
	SelectedMethod    *string
	AvailableMethods  []string
	RequestPayload    string
	DescriptorSource  grpcurl.DescriptorSource
}

func NewController() Controller {

	ui := ui.NewView()
	conn := Connection{
		ServerURL: "localhost:50051",
	}

	c := Controller{
		context.Background(),
		ui,
		&conn,
	}

	return c
}

func (c Controller) Run() error {
	c.initMenu()

	c.initSys()

	c.ui.App.EnableMouse(true)
	return c.ui.App.Run()
}
