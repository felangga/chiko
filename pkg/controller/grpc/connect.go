package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
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
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()

	var creds credentials.TransportCredentials

	// Reset connection state
	if err := g.resetActiveConnection(); err != nil {
		return fmt.Errorf("failed to reset connection: %w", err)
	}

	g.Conn.ServerURL = serverURL
	g.logInfo("🌏 server URL set to [blue]" + serverURL + ", connecting...")

	if g.Conn.EnableTLS {
		// Configure TLS credentials
		var err error
		creds, err = g.configureTLSCredentials()
		if err != nil {
			return fmt.Errorf("TLS configuration error: %w", err)
		}
	}

	// Dial options
	opts := []grpc.DialOption{
		grpc.WithUserAgent("chiko/" + entity.APP_VERSION),
	}

	// Establish gRPC connection
	conn, err := grpcurl.BlockingDial(ctx, "tcp", serverURL, creds, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial server %s: %w", serverURL, err)
	}
	g.Conn.ActiveConnection = conn
	g.logInfo("✅ connected to [blue]" + serverURL)

	// Server reflection
	if err := g.setupServerReflection(ctx, conn); err != nil {
		return err
	}

	return nil
}

func (g *GRPC) resetActiveConnection() error {
	if g.Conn.ActiveConnection == nil {
		return nil
	}

	// Attempt to close the connection with error handling
	err := g.Conn.ActiveConnection.Close()
	if err != nil {
		// Log the error but don't propagate it
		g.logWarning(fmt.Sprintf("Error closing active connection: %v", err))
	}

	// Explicitly set ActiveConnection to nil after closing
	g.Conn.ActiveConnection = nil

	return nil
}

func (g *GRPC) configureTLSCredentials() (credentials.TransportCredentials, error) {
	if g.Conn.SSLCert == nil {
		return insecure.NewCredentials(), nil
	}

	caCertPool := x509.NewCertPool()
	var certs tls.Certificate

	if g.Conn.SSLCert.ClientCert_Path != nil && g.Conn.SSLCert.ClientKey_Path != nil {
		var err error
		certs, err = tls.LoadX509KeyPair(*g.Conn.SSLCert.ClientCert_Path, *g.Conn.SSLCert.ClientKey_Path)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
	}

	if g.Conn.SSLCert.CA_Path != nil {
		caCert, err := os.ReadFile(*g.Conn.SSLCert.CA_Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate to pool")
		}
	}

	return credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: g.Conn.InsecureSkipVerify,
		Certificates:       []tls.Certificate{certs},
		RootCAs:            caCertPool,
	}), nil
}

func (g *GRPC) setupServerReflection(ctx context.Context, conn *grpc.ClientConn) error {
	refClient := grpcreflect.NewClientV1Alpha(ctx, reflectpb.NewServerReflectionClient(conn))
	reflSource := grpcurl.DescriptorSourceFromServer(ctx, refClient)

	svcs, err := grpcurl.ListServices(reflSource)
	if err != nil {
		g.logWarning("❗️ connected but failed to get services from server reflection")
		return nil
	}

	g.logInfo("✅ this server supports server reflection")
	g.Conn.DescriptorSource = reflSource
	g.Conn.AvailableMethods = make([]string, 0, len(svcs)*5) // Preallocate with reasonable estimate
	g.Conn.AvailableServices = make([]string, 0, len(svcs))

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

func (g *GRPC) logInfo(message string) {
	log := entity.Log{
		Content: message,
		Type:    entity.LOG_INFO,
	}
	log.DumpLogToChannel(g.LogChannel)
}

func (g *GRPC) logWarning(message string) {
	log := entity.Log{
		Content: message,
		Type:    entity.LOG_WARNING,
	}
	log.DumpLogToChannel(g.LogChannel)
}
