package entity

import (
	"github.com/fullstorydev/grpcurl"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const APP_VERSION = "0.0.1"

type Connection struct {
	ID                *uuid.UUID               `json:"id"`
	Name              string                   `json:"name"`
	ServerURL         string                   `json:"server_url"`
	ActiveConnection  *grpc.ClientConn         `json:"-"`
	AvailableServices []string                 `json:"-"`
	SelectedMethod    *string                  `json:"-"`
	AvailableMethods  []string                 `json:"-"`
	RequestPayload    string                   `json:"request_payload"`
	DescriptorSource  grpcurl.DescriptorSource `json:"-"`
}

var Bookmarks []Connection
