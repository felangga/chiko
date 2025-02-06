package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/felangga/chiko/pkg/entity"
)

const GRPC_TIMEOUT = time.Second * 10

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
		Ctx:           context.Background(),
		Conn:          &conn,
		LogChannel:    logChannel,
		OutputChannel: outputChannel,
	}

	return g
}
