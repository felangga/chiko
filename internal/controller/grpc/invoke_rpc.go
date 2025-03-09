package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GRPC) Validate() error {
	if g.Conn.SelectedMethod == nil {
		return fmt.Errorf("❗ no method selected")
	}
	if g.Conn.ActiveConnection == nil {
		return fmt.Errorf("❗ no active connection")
	}

	return nil
}

// InvokeRPC will invoke the configured payload and try to hit the server with it
func (g *GRPC) InvokeRPC() error {
	if err := g.Validate(); err != nil {
		return err
	}

	ctx, timeout := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer timeout()

	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		AllowUnknownFields:    g.Conn.AllowUnknownFields,
	}

	rf, formatter, err := grpcurl.RequestParserAndFormatter(
		grpcurl.Format("json"),
		g.Conn.DescriptorSource,
		strings.NewReader(g.Conn.RequestPayload),
		options,
	)
	if err != nil {
		return err
	}

	h := &handler{grpc: *g}

	err = grpcurl.InvokeRPC(
		ctx,
		g.Conn.DescriptorSource,
		g.Conn.ActiveConnection,
		*g.Conn.SelectedMethod,
		g.Conn.ParseMetadata(),
		h,
		rf.Next,
	)
	if err != nil {
		if errStatus, ok := status.FromError(err); ok {
			h.respStatus = errStatus
		} else {
			return err
		}
	}

	if h.respStatus.Code() != codes.OK {
		formattedStatus, err := formatter(h.respStatus.Proto())
		if err != nil {
			return err
		}
		return fmt.Errorf("%s", formattedStatus)
	}

	return nil
}
