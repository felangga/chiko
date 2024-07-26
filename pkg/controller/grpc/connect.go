package grpc

import (
	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/felangga/chiko/pkg/entity"
)

// Connect will connect to the server and try to get the available services and methods
func (g *GRPC) Connect(serverURL string) error {
	g.Conn.ServerURL = serverURL

	// Close active connection if we are going to connect to another server
	if g.Conn.ActiveConnection != nil {
		err := g.Conn.ActiveConnection.Close()
		if err != nil {
			return err
		}
		g.Conn.ActiveConnection = nil
	}

	log := entity.Log{
		Content: "🌏 server URL set to [blue]" + g.Conn.ServerURL,
		Type:    entity.LOG_INFO,
	}
	log.DumpLogToChannel(g.LogChannel)

	conn, err := grpcurl.BlockingDial(g.Ctx, "tcp", serverURL, nil)
	if err != nil {
		return err
	}
	g.Conn.ActiveConnection = conn

	log = entity.Log{
		Content: "✅ connected to [blue]" + g.Conn.ServerURL,
		Type:    entity.LOG_INFO,
	}
	log.DumpLogToChannel(g.LogChannel)

	refClient := grpcreflect.NewClientV1Alpha(g.Ctx, reflectpb.NewServerReflectionClient(conn))
	reflSource := grpcurl.DescriptorSourceFromServer(g.Ctx, refClient)
	svcs, err := grpcurl.ListServices(reflSource)
	if err != nil {
		return err
	}
	g.Conn.DescriptorSource = reflSource
	log = entity.Log{
		Content: "✅ this server support server reflection",
		Type:    entity.LOG_INFO,
	}
	log.DumpLogToChannel(g.LogChannel)

	g.Conn.AvailableMethods = []string{} // Reset available methods
	for _, svc := range svcs {
		g.Conn.AvailableServices = append(g.Conn.AvailableServices, svc)
		methods, err := grpcurl.ListMethods(reflSource, svc)
		if err != nil {
			return err
		}
		g.Conn.AvailableMethods = append(g.Conn.AvailableMethods, methods...)
	}

	return nil
}
