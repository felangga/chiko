package grpc

import (
	"time"

	"github.com/google/uuid"

	"github.com/felangga/chiko/pkg/entity"
)

const GRPC_TIMEOUT = time.Second * 10

type GRPC struct {
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
		&conn,
		logChannel,
	}

	return g
}
