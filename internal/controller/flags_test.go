package controller

import (
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Set up the flags
	os.Args = []string{"cmd", "--plaintext", "--connect-timeout", "10", "--cacert", "path/to/cacert", "grpcb.in:9000"}

	// Call the ParseFlags function
	session, err := ParseFlags()

	// Check for errors
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Validate the session object
	if session.EnableTLS != false {
		t.Errorf("expected EnableTLS to be true, got %v", session.EnableTLS)
	}
	if *session.SSLCert.CA_Path != "path/to/cacert" {
		t.Errorf("expected CACert to be 'path/to/cacert', got %v", session.SSLCert.CA_Path)
	}
	if session.ServerURL != "grpcb.in:9000" {
		t.Errorf("expected ServerURL to be 'grpcb.in:9000', got %v", session.ServerURL)
	}
}
