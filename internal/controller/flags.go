package controller

import (
	"flag"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/felangga/chiko/internal/entity"
)

// This section of code is adapted from https://github.com/fullstorydev/grpcurl/blob/master/cmd/grpcurl/grpcurl.go#L892
// with some modification to make it work with the flags package

var (
	flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	plaintext = flags.Bool("plaintext", false, prettify(`
		Use plain-text HTTP/2 when connecting to server (no TLS).`))
	insecure = flags.Bool("insecure", false, prettify(`
		Skip server certificate and domain verification. (NOT SECURE!) Not
		valid with -plaintext option.`))

	// TLS Options
	cacert = flags.String("cacert", "", prettify(`
		File containing trusted root certificates for verifying the server.
		Ignored if -insecure is specified.`))
	cert = flags.String("cert", "", prettify(`
		File containing client certificate (public key), to present to the
		server. Not valid with -plaintext option. Must also provide -key option.`))
	key = flags.String("key", "", prettify(`
		File containing client private key, to present to the server. Not valid
		with -plaintext option. Must also provide -cert option.`))

	data = flags.String("d", "", prettify(`
		Data for request contents. If the value is '@' then the request contents
		are read from stdin. For calls that accept a stream of requests, the
		contents should include all such request messages concatenated together
		(possibly delimited; see -format).`))
	allowUnknownFields = flags.Bool("allow-unknown-fields", false, prettify(`
		When true, the request contents, if 'json' format is used, allows
		unknown fields to be present. They will be ignored when parsing
		the request.`))
	connectTimeout = flags.Float64("connect-timeout", 0, prettify(`
		The maximum time, in seconds, to wait for connection to be established.
		Defaults to 10 seconds.`))
	formatError = flags.Bool("format-error", false, prettify(`
		When a non-zero status is returned, format the response using the
		value set by the -format flag .`))
	keepaliveTime = flags.Float64("keepalive-time", 0, prettify(`
		If present, the maximum idle time in seconds, after which a keepalive
		probe is sent. If the connection remains idle and no keepalive response
		is received for this same period then the connection is closed and the
		operation fails.`))
	maxTime = flags.Float64("max-time", 0, prettify(`
		The maximum total time the operation can take, in seconds. This sets a
                timeout on the gRPC context, allowing both client and server to give up
		after the deadline has past. This is useful for preventing batch jobs
                that use grpcurl from hanging due to slow or bad network links or due
		to incorrect stream method usage.`))
	maxMsgSz = flags.Int("max-msg-sz", 0, prettify(`
		The maximum encoded size of a response message, in bytes, that grpcurl
		will accept. If not specified, defaults to 4,194,304 (4 megabytes).`))
	emitDefaults = flags.Bool("emit-defaults", false, prettify(`
		Emit default values for JSON-encoded responses.`))
	protosetOut = flags.String("protoset-out", "", prettify(`
		The name of a file to be written that will contain a FileDescriptorSet
		proto. With the list and describe verbs, the listed or described
		elements and their transitive dependencies will be written to the named
		file if this option is given. When invoking an RPC and this option is
		given, the method being invoked and its transitive dependencies will be
		included in the output file.`))
	protoOut = flags.String("proto-out-dir", "", prettify(`
		The name of a directory where the generated .proto files will be written.
		With the list and describe verbs, the listed or described elements and
		their transitive dependencies will be written as .proto files in the
		specified directory if this option is given. When invoking an RPC and
		this option is given, the method being invoked and its transitive
		dependencies will be included in the generated .proto files in the
		output directory.`))
	msgTemplate = flags.Bool("msg-template", false, prettify(`
		When describing messages, show a template of input data.`))
	verbose = flags.Bool("v", false, prettify(`
		Enable verbose output.`))
	veryVerbose = flags.Bool("vv", false, prettify(`
		Enable very verbose output (includes timing data).`))
	serverName = flags.String("servername", "", prettify(`
		Override server name when validating TLS certificate. This flag is
		ignored if -plaintext or -insecure is used.
		NOTE: Prefer -authority. This flag may be removed in the future. It is
		an error to use both -authority and -servername (though this will be
		permitted if they are both set to the same value, to increase backwards
		compatibility with earlier releases that allowed both to be set).`))
)

// ParseSession handles command-line flags and returns a Session configuration
func ParseFlags() entity.Session {
	plainText := flag.Bool("plaintext", false, "use plaintext mode")
	insecure := flag.Bool("insecure", false, "allow insecure server connections")
	payload := flag.String("d", "", "payload to send")
	// describe := flag.String("describe", "", "get reflection from the server")

	flag.Parse()

	return entity.Session{
		ID:                 uuid.New(),
		ServerURL:          "localhost:20010",
		EnableTLS:          !*plainText,
		InsecureSkipVerify: *insecure,
		RequestPayload:     *payload,
	}
}

func prettify(docString string) string {
	parts := strings.Split(docString, "\n")

	// cull empty lines and also remove trailing and leading spaces
	// from each line in the doc string
	j := 0
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		parts[j] = part
		j++
	}

	return strings.Join(parts[:j], "\n")
}
