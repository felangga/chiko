package grpc

import (
	"fmt"
	"strings"

	"github.com/felangga/chiko/pkg/entity"
	"github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InvokeRPC will invoke the configured payload and try to hit the server with it
func (g *GRPC) InvokeRPC() error {
	if g.Conn.SelectedMethod == nil {
		return fmt.Errorf("❗ no method selected")
	}

	if g.Conn.ActiveConnection == nil {
		return fmt.Errorf("❗ no active connection")
	}

	// Construct metadata info
	var metadata string
	for _, meta := range g.Conn.ParseMetadata() {
		metadata += "- " + meta + "\n"
	}

	log := entity.Log{
		Content: "\nRequest Metadata:\n" + metadata + "\nRequest Payload:\n[yellow]" + g.Conn.RequestPayload + "\n",
		Type:    entity.LOG_INFO,
	}
	log.DumpLogToChannel(g.LogChannel)

	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		AllowUnknownFields:    true,
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

	h := &handler{
		grpc: *g,
	}

	err = grpcurl.InvokeRPC(
		g.Ctx,
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

		return fmt.Errorf(formattedStatus)
	}

	return nil
}
