package controller

import (
	"chiko/pkg/entity"
	"fmt"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InvokeRPC will invoke the configured payload and try to hit the server with it
func (c Controller) InvokeRPC() error {
	if c.Conn.SelectedMethod == nil {
		return fmt.Errorf("❗ no method selected")
	}

	if c.Conn.ActiveConnection == nil {
		return fmt.Errorf("❗ no active connection")
	}

	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		AllowUnknownFields:    true,
	}
	rf, formatter, err := grpcurl.RequestParserAndFormatter(
		grpcurl.Format("json"),
		c.Conn.DescriptorSource,
		strings.NewReader(c.Conn.RequestPayload),
		options,
	)
	if err != nil {
		return err
	}

	h := &handler{
		controller: c,
	}

	err = grpcurl.InvokeRPC(
		c.Ctx,
		c.Conn.DescriptorSource,
		c.Conn.ActiveConnection,
		*c.Conn.SelectedMethod,
		c.Conn.ParseMetadata(),
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

		c.PrintLog(entity.LogParam{
			Content: formattedStatus,
			Type:    entity.LOG_ERROR,
		})
	}

	return nil
}
