package grpc

import (
	"chiko/pkg/entity"
	"context"

	"github.com/google/uuid"
)

type GRPC struct {
	Ctx        context.Context
	Conn       *entity.Session
	LogChannel chan entity.Log
}

// NewGRPC is used to create a new grpc object
func NewGRPC(logChannel chan entity.Log) GRPC {
	conn := entity.Session{
		// Default server URL
		ID:        uuid.New(),
		ServerURL: "localhost:20010",
	}

	g := GRPC{
		context.Background(),
		&conn,
		logChannel,
	}

	return g
}
