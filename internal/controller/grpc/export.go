package grpc

import (
	"errors"
	"fmt"
	"strings"
)

func (g *GRPC) validateExport() error {
	if g.Conn.ServerURL == "" {
		return errors.New("host is empty")
	}

	if g.Conn.SelectedMethod == nil {
		return errors.New("no method selected")
	}

	if g.Conn.RequestPayload == "" {
		return errors.New("no request payload")
	}

	return nil
}

// Export is used to export the current grpc session into grpcurl commands
func (g *GRPC) ExportGrpcurlCommand() (string, error) {
	if err := g.validateExport(); err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString("grpcurl")

	// TLS options
	if !g.Conn.EnableTLS {
		builder.WriteString(" -plaintext")
	}
	if g.Conn.InsecureSkipVerify {
		builder.WriteString(" -insecure")
	}

	if g.Conn.SSLCert != nil {
		if g.Conn.SSLCert.CA_Path != nil {
			builder.WriteString(" -cacert " + *g.Conn.SSLCert.CA_Path)
		}
		if g.Conn.SSLCert.ClientCert_Path != nil {
			builder.WriteString(" -cert " + *g.Conn.SSLCert.ClientCert_Path)
		}
		if g.Conn.SSLCert.ClientKey_Path != nil {
			builder.WriteString(" -key " + *g.Conn.SSLCert.ClientKey_Path)
		}
	}

	// Request options
	if g.Conn.RequestPayload != "" {
		escapedPayload := strings.Replace(string(g.Conn.RequestPayload), "'", "'\\''", -1)
		builder.WriteString(fmt.Sprintf(" -d '%s'", escapedPayload))
	}
	if g.Conn.AllowUnknownFields {
		builder.WriteString(" -allow-unknown-fields")
	}

	// Timeout options
	if g.Conn.ConnectTimeout != 10 {
		builder.WriteString(fmt.Sprintf(" -connect-timeout %f", g.Conn.ConnectTimeout))
	}
	if g.Conn.KeepAliveTime != 0 {
		builder.WriteString(fmt.Sprintf(" -keepalive-time %f", g.Conn.KeepAliveTime))
	}
	if g.Conn.MaxMsgSz != GRPC_MAX_MESSAGE_SIZE {
		builder.WriteString(fmt.Sprintf(" -max-msg-sz %d", g.Conn.MaxMsgSz))
	}
	if g.Conn.MaxTimeOut != 0 {
		builder.WriteString(fmt.Sprintf(" -max-time %f", g.Conn.MaxTimeOut))
	}

	// Server URL
	builder.WriteString(" " + g.Conn.ServerURL)

	// Selected Method
	builder.WriteString(" " + *g.Conn.SelectedMethod)

	return builder.String(), nil
}
