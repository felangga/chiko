package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/felangga/chiko/pkg/entity"
)

type GRPC struct {
	Ctx           context.Context
	Conn          *entity.Session
	LogChannel    chan entity.Log
	OutputChannel chan entity.Output
}

// NewGRPC is used to create a new grpc object
func NewGRPC(logChannel chan entity.Log, outputChannel chan entity.Output) GRPC {
	conn := entity.Session{
		// Default server URL
		ID:        uuid.New(),
		ServerURL: "localhost:20010",
	}

	g := GRPC{
		context.Background(),
		&conn,
		logChannel,
		outputChannel,
	}

	return g
}
