package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/felangga/chiko/pkg/entity"
)

// Connect is used to connect to the server and doing check if the server support server reflection
func (g *GRPC) Connect(serverURL string) error {
	context, timeout := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer timeout()

	var (
		opts  []grpc.DialOption
		creds credentials.TransportCredentials
		certs tls.Certificate
	)

	creds = insecure.NewCredentials()
	caCertPool := x509.NewCertPool()
	opts = append(opts, grpc.WithUserAgent("chiko/"+entity.APP_VERSION))

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

	if g.Conn.SSLCert != nil {
		if g.Conn.SSLCert.ClientCert_Path != nil && g.Conn.SSLCert.ClientKey_Path != nil {
			var err error

			// Load the client's certificate and private key
			certs, err = tls.LoadX509KeyPair("./cert/client-cert.pem", "./cert/client-key.pem")
			if err != nil {
				return err
			}
		}

		if g.Conn.SSLCert.CA_Path != nil {
			// Load the CA certificate
			caCert, err := os.ReadFile(*g.Conn.SSLCert.CA_Path)
			if err != nil {
				log := entity.Log{
					Content: err.Error(),
					Type:    entity.LOG_INFO,
				}
				log.DumpLogToChannel(g.LogChannel)
				return err
			}

			if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
				return err
			}
		}

		creds = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: g.Conn.InsecureSkipVerify,
			Certificates:       []tls.Certificate{certs},
			RootCAs:            caCertPool,
		})
	}

	conn, err := grpcurl.BlockingDial(context, "tcp", serverURL, creds, opts...)
	if err != nil {
		return err
	}

	g.Conn.ActiveConnection = conn

	log = entity.Log{
		Content: "✅ connected to [blue]" + g.Conn.ServerURL,
		Type:    entity.LOG_INFO,
	}
	log.DumpLogToChannel(g.LogChannel)

	refClient := grpcreflect.NewClientV1Alpha(context, reflectpb.NewServerReflectionClient(conn))
	reflSource := grpcurl.DescriptorSourceFromServer(context, refClient)
	svcs, err := grpcurl.ListServices(reflSource)
	if err != nil {
		log = entity.Log{
			Content: "❗️ connected but failed to get services from server reflection",
			Type:    entity.LOG_WARNING,
		}
		log.DumpLogToChannel(g.LogChannel)

		return nil
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
