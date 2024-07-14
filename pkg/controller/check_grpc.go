package controller

import (
	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"chiko/pkg/entity"
)

// CheckGRPC will check if the server supports server reflection and list all available services and methods
func (c Controller) CheckGRPC(serverURL string) error {
	c.Conn.ServerURL = serverURL

	// Close active connection if we are going to connect to another server
	if c.Conn.ActiveConnection != nil {
		err := c.Conn.ActiveConnection.Close()
		if err != nil {
			return err
		}
		c.Conn.ActiveConnection = nil
	}

	c.PrintLog(entity.LogParam{
		Content: "🌏 server URL set to [blue]" + c.Conn.ServerURL,
		Type:    entity.LOG_INFO,
	})

	c.PrintLog(entity.LogParam{
		Content: "🔍 checking server reflection...",
		Type:    entity.LOG_INFO,
	})

	conn, err := grpcurl.BlockingDial(c.Ctx, "tcp", serverURL, nil)
	if err != nil {
		return err
	}
	c.Conn.ActiveConnection = conn
	c.PrintLog(entity.LogParam{
		Content: "✅ connected to [blue]" + c.Conn.ServerURL,
		Type:    entity.LOG_INFO,
	})

	refClient := grpcreflect.NewClientV1Alpha(c.Ctx, reflectpb.NewServerReflectionClient(conn))
	reflSource := grpcurl.DescriptorSourceFromServer(c.Ctx, refClient)
	svcs, err := grpcurl.ListServices(reflSource)
	if err != nil {
		return err
	}
	c.Conn.DescriptorSource = reflSource
	c.PrintLog(entity.LogParam{
		Content: "🤩 this server support server reflection",
		Type:    entity.LOG_INFO,
	})

	c.Conn.AvailableMethods = []string{} // Reset available methods
	for _, svc := range svcs {
		c.Conn.AvailableServices = append(c.Conn.AvailableServices, svc)
		methods, err := grpcurl.ListMethods(reflSource, svc)
		if err != nil {
			return err
		}
		c.Conn.AvailableMethods = append(c.Conn.AvailableMethods, methods...)
	}

	return nil
}
