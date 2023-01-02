package controller

import (
	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
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

	c.PrintLog(" 🤩 this server support server reflection", LOG_INFO)
	for _, svc := range svcs {
		c.conn.AvailableServices = append(c.conn.AvailableServices, svc)
		methods, err := grpcurl.ListMethods(reflSource, svc)
		if err != nil {
			c.PrintLog(err.Error(), LOG_ERROR)
		}
		for _, method := range methods {
			c.PrintLog(method, LOG_WARNING)
			c.conn.AvailableMethods = append(c.conn.AvailableMethods, method)
		}
	}

}
