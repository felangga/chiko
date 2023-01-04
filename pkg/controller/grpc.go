package controller

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc/codes"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

func (c Controller) CheckGRPC() {
	conn, err := grpcurl.BlockingDial(c.ctx, "tcp", c.conn.ServerURL, nil)
	if err != nil {
		c.PrintLog(" ⛔️ "+err.Error(), LOG_ERROR)
		return
	}
	c.conn.ActiveConnection = conn
	c.PrintLog(" ✅ connected to "+c.conn.ServerURL, LOG_INFO)
	refClient := grpcreflect.NewClientV1Alpha(c.ctx, reflectpb.NewServerReflectionClient(conn))
	reflSource := grpcurl.DescriptorSourceFromServer(c.ctx, refClient)
	svcs, err := grpcurl.ListServices(reflSource)
	if err != nil {
		c.PrintLog(err.Error(), LOG_ERROR)
		return
	}
	c.conn.DescriptorSource = reflSource
	c.PrintLog(" 🤩 this server support server reflection", LOG_INFO)
	for _, svc := range svcs {
		c.conn.AvailableServices = append(c.conn.AvailableServices, svc)
		methods, err := grpcurl.ListMethods(reflSource, svc)
		if err != nil {
			c.PrintLog(err.Error(), LOG_ERROR)
			return
		}
		c.conn.AvailableMethods = append(c.conn.AvailableMethods, methods...)
	}

}

func (c Controller) parseRequestResponse(text string) [][]string {
	var re = regexp.MustCompile(`\(([^()]*)\)`)
	return re.FindAllStringSubmatch(text, -1)
}

func (c Controller) doInvoke() {
	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		AllowUnknownFields:    true,
	}
	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), c.conn.DescriptorSource, strings.NewReader(c.conn.RequestPayload), options)
	if err != nil {
		c.PrintLog(err.Error(), LOG_ERROR)
		return
	}
	h := &handler{
		controller: c,
	}

	err = grpcurl.InvokeRPC(c.ctx, c.conn.DescriptorSource, c.conn.ActiveConnection, *c.conn.SelectedMethod, nil, h, rf.Next)
	if err != nil {
		if errStatus, ok := status.FromError(err); ok {
			h.respStatus = errStatus
		} else {
			c.PrintLog(fmt.Sprintf("error invoking method %s", err.Error()), LOG_ERROR)
			return
		}
	}

	if h.respStatus.Code() != codes.OK {
		formattedStatus, err := formatter(h.respStatus.Proto())
		if err != nil {
			c.PrintLog(err.Error(), LOG_ERROR)
			return
		}
		c.PrintLog(formattedStatus, LOG_ERROR)
	}
}
