package controller

import (
	"flag"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/felangga/chiko/internal/entity"
)

// Config holds all the configuration options for the application
type Flag struct {
	Plaintext          bool
	Insecure           bool
	CACert             string
	Cert               string
	Key                string
	Data               string
	AllowUnknownFields bool
	ConnectTimeout     float64
	KeepAliveTime      float64
	MaxTime            float64
	MaxMsgSz           int
}

func (f Flag) Validate() error {
	if f.Plaintext && f.Insecure {
		return fmt.Errorf("cannot use -plaintext and -insecure together")
	}
	if (f.Cert == "") != (f.Key == "") {
		return fmt.Errorf("must provide both -cert and -key, or neither")
	}
	if f.ConnectTimeout < 0 {
		return fmt.Errorf("The -connect-timeout argument must not be negative.")
	}
	if f.KeepAliveTime < 0 {
		return fmt.Errorf("The -keepalive-time argument must not be negative.")
	}
	if f.MaxTime < 0 {
		return fmt.Errorf("The -max-time argument must not be negative.")
	}
	if f.MaxMsgSz < 0 {
		return fmt.Errorf("The -max-msg-sz argument must not be negative.")
	}
	if f.Plaintext && f.Insecure {
		return fmt.Errorf("The -plaintext and -insecure arguments are mutually exclusive.")
	}
	if f.Insecure && !f.Plaintext {
		return fmt.Errorf("The -insecure argument can only be used with TLS.")
	}
	if f.Cert != "" && !f.Plaintext {
		return fmt.Errorf("The -cert argument can only be used with TLS.")
	}
	if f.Key != "" && !f.Plaintext {
		return fmt.Errorf("The -key argument can only be used with TLS.")
	}
	if (f.Key == "") != (f.Cert == "") {
		return fmt.Errorf("The -cert and -key arguments must be used together and both be present.")
	}

	return nil
}

// ParseFlags handles command-line flags and returns a Session configuration
func ParseFlags() (entity.Session, error) {
	f := Flag{}

	// Flags from grpcurl
	// Some flags are removed due not supported with Chiko
	flag.BoolVar(&f.Plaintext, "plaintext", false, "Use plain-text HTTP/2 when connecting to server (no TLS)")
	flag.BoolVar(&f.Insecure, "insecure", false, "Skip server certificate and domain verification")
	flag.StringVar(&f.CACert, "cacert", "", "File containing trusted root certificates for verifying the server")
	flag.StringVar(&f.Cert, "cert", "", "File containing client certificate (public key)")
	flag.StringVar(&f.Key, "key", "", "File containing client private key")
	flag.StringVar(&f.Data, "d", "", "Data for request contents")
	flag.BoolVar(&f.AllowUnknownFields, "allow-unknown-fields", false, "Allow unknown fields in JSON request")
	flag.Float64Var(&f.ConnectTimeout, "connect-timeout", 10, "Maximum time to wait for connection (seconds)")
	flag.Float64Var(&f.KeepAliveTime, "keepalive-time", 0, "Maximum idle time before keepalive probe")
	flag.Float64Var(&f.MaxTime, "max-time", 0, "Maximum total operation time")
	flag.IntVar(&f.MaxMsgSz, "max-msg-sz", 4*1024*1024, "Maximum encoded response message size")

	flag.Parse()

	if err := f.Validate(); err != nil {
		return entity.Session{}, err
	}

	var (
		method    *string
		serverURL string = "localhost:20010"
	)

	if flag.NArg() > 0 {
		args := flag.Args()

		// Skip "list" or "describe" if present
		for len(args) > 0 && (args[0] == "list" || args[0] == "describe") {
			args = args[1:]
		}

		if len(args) > 0 {
			serverURL = args[0]
			args = args[1:]

			// Skip "list" or "describe" again if present
			if len(args) > 0 && (args[0] == "list" || args[0] == "describe") {
				args = args[1:]
			}

			if len(args) > 0 {
				normalize := strings.ReplaceAll(args[0], "/", ".")
				method = &normalize
			}
		}
	}

	return entity.Session{
		ID:                 uuid.New(),
		ServerURL:          serverURL,
		EnableTLS:          !f.Plaintext,
		InsecureSkipVerify: f.Insecure,
		RequestPayload:     f.Data,
		AllowUnknownFields: f.AllowUnknownFields,
		MaxMsgSz:           f.MaxMsgSz,
		MaxTimeOut:         f.MaxTime,
		ConnectTimeout:     f.ConnectTimeout,
		KeepAliveTime:      f.KeepAliveTime,
		SelectedMethod:     method,
		SSLCert: &entity.Cert{
			CA_Path:         &f.CACert,
			ClientCert_Path: &f.Cert,
			ClientKey_Path:  &f.Key,
		},
	}, nil
}
