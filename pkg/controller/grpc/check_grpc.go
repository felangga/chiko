package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/felangga/chiko/pkg/entity"
)

// CheckGRPC will check if the server supports server reflection and list all available services and methods
func (g *GRPC) CheckGRPC(serverURL string) error {
	var (
		opts  []grpc.DialOption
		creds credentials.TransportCredentials
	)

	opts = append(opts, grpc.WithUserAgent(fmt.Sprintf("chiko/%d", entity.APP_VERSION)))

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
		Content: "🌏 server URL set to [blue]" + g.Conn.ServerURL + ", connecting...",
		Type:    entity.LOG_INFO,
	}
	log.DumpLogToChannel(g.LogChannel)

	// Load the client's certificate and private key
	cert, err := tls.LoadX509KeyPair("./tls/client.crt", "./tls/client.key")
	if err != nil {

		log := entity.Log{
			Content: err.Error(),
			Type:    entity.LOG_INFO,
		}
		log.DumpLogToChannel(g.LogChannel)
		return err
	}

	// Load the CA certificate
	caCert, err := ioutil.ReadFile("./tls/ca.crt")
	if err != nil {
		log := entity.Log{
			Content: err.Error(),
			Type:    entity.LOG_INFO,
		}
		log.DumpLogToChannel(g.LogChannel)
		return err
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return err
	}

	if g.IsSecure {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		})
	}

	conn, err := grpcurl.BlockingDial(g.Ctx, "tcp", serverURL, creds, opts...)
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
