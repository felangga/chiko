package grpc

import (
	"fmt"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
)

// GenerateRPCPayloadSample is used to generate sample request payload that will be filled to txtPayload text field
func (g *GRPC) GenerateRPCPayloadSample() (string, error) {
	dsc, err := g.Conn.DescriptorSource.FindSymbol(*g.Conn.SelectedMethod)
	if err != nil {
		return "", err
	}

	txt, err := grpcurl.GetDescriptorText(dsc, g.Conn.DescriptorSource)
	if err != nil {
		return "", err
	}

	rr := g.parseRequestResponse(txt)
	if len(rr) < 2 {
		return "", fmt.Errorf("failed to parse request response: %s", txt)
	}

	requestMessage := strings.ReplaceAll(rr[0][1], "stream", "")
	requestMessage = strings.TrimSpace(requestMessage)
	if requestMessage[0:1] == "." {
		requestMessage = requestMessage[1:]
	}

	dsc, err = g.Conn.DescriptorSource.FindSymbol(requestMessage)
	if err != nil {
		return "", err
	}

	if dsc, ok := dsc.(*desc.MessageDescriptor); ok {
		tmpl := grpcurl.MakeTemplate(dsc)
		options := grpcurl.FormatOptions{EmitJSONDefaultFields: true}
		_, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), g.Conn.DescriptorSource, nil, options)
		if err != nil {
			return "", err
		}
		str, err := formatter(tmpl)
		if err != nil {
			return "", err
		}
		return str, nil
	}

	return "", fmt.Errorf("failed to generate sample")
}
