package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/felangga/chiko/internal/entity"
)

// Connect is used to connect to the server and doing check if the server support server reflection
func (g *GRPC) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()

	if g.Conn.MaxTimeOut > 0 {
		timeout := time.Duration(g.Conn.MaxTimeOut * float64(time.Second))
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	var creds credentials.TransportCredentials

	// Reset connection state
	if err := g.resetActiveConnection(); err != nil {
		return fmt.Errorf("failed to reset connection: %w", err)
	}

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

	if g.Conn.KeepAliveTime > 0 {
		timeout := time.Duration(g.Conn.KeepAliveTime * float64(time.Second))
		opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    timeout,
			Timeout: timeout,
		}))
	}

	if g.Conn.MaxMsgSz > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(g.Conn.MaxMsgSz)))
	}

	// Establish gRPC connection
	conn, err := grpcurl.BlockingDial(ctx, "tcp", g.Conn.ServerURL, creds, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial server %s: %w", g.Conn.ServerURL, err)
	}
	g.Conn.ActiveConnection = conn
	g.Logger.Info("✅ connected to [blue]" + g.Conn.ServerURL)

	// Server reflection
	if err := g.setupServerReflection(ctx, conn); err != nil {
		return err
	}

	// Check if selected service and methods is not available
	err = g.CheckSelectedMethod()
	if err != nil {
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
		g.Logger.Warning(fmt.Sprintf("Error closing active connection: %v", err))
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
		g.Logger.Warning("❗️ connected but failed to get services from server reflection")
		return nil
	}

	g.Logger.Info("✅ this server supports server reflection")
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

func (g *GRPC) CheckSelectedMethod() error {
	for _, methods := range g.Conn.AvailableMethods {
		if methods == *g.Conn.SelectedMethod {
			return nil
		}
	}

	return fmt.Errorf("❗ method %s not found", *g.Conn.SelectedMethod)
}
