package grpc

import (
	"context"
	"time"

	"github.com/felangga/chiko/internal/entity"
	"github.com/felangga/chiko/internal/logger"
)

const GRPC_TIMEOUT = time.Second * 10

type GRPC struct {
	Ctx    context.Context
	Conn   *entity.Session
	Logger *logger.Logger
}

// NewGRPC is used to create a new grpc object
func NewGRPC(logger *logger.Logger, session *entity.Session) GRPC {
	g := GRPC{
		Ctx:    context.Background(),
		Conn:   session,
		Logger: logger,
	}

	return g
}
