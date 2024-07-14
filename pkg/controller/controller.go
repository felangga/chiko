package controller

import (
	"context"

	"github.com/google/uuid"

	"chiko/pkg/entity"
)

type Controller struct {
	Ctx       context.Context
	Conn      *entity.Session
	Bookmarks *[]entity.Bookmark
	LogDump   chan entity.LogParam
}

func NewController() Controller {
	conn := entity.Session{
		// Default server URL
		ID:        uuid.New(),
		ServerURL: "localhost:20010",
	}

	bookmarks := []entity.Bookmark{}
	logDump := make(chan entity.LogParam)

	c := Controller{
		context.Background(),
		&conn,
		&bookmarks,
		logDump,
	}

	return c
}
