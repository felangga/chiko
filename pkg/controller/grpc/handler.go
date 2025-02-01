package grpc

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/felangga/chiko/pkg/entity"
)

type handler struct {
	grpc GRPC

	method            *desc.MethodDescriptor
	methodCount       int
	reqHeaders        metadata.MD
	reqHeadersCount   int
	respHeaders       metadata.MD
	respHeadersCount  int
	respMessages      []string
	respTrailers      metadata.MD
	respStatus        *status.Status
	respTrailersCount int
}

func (h *handler) OnReceiveResponse(msg proto.Message) {
	// Print headers to log window
	headerResp := "Headers:"
	for key, values := range h.respHeaders {
		headerResp += fmt.Sprintf("\n  - %s: %s", key, strings.Join(values, ","))
	}

	statusCode := fmt.Sprintf("Status code: %d %s", h.respStatus.Code(), h.respStatus.Message())

	// Print the gRPC response
	jsm := jsonpb.Marshaler{Indent: "  "}
	respStr, err := jsm.MarshalToString(msg)
	if err != nil {
		log := entity.Log{
			Content: fmt.Sprintf("failed to generate JSON form of response message: %v", err),
			Type:    entity.LOG_ERROR,
		}
		log.DumpLogToChannel(h.grpc.LogChannel)
		return
	}

	h.respMessages = append(h.respMessages, respStr)
	output := fmt.Sprintf("\n%s\n\n%s\n%s", headerResp, statusCode, respStr)
	out := entity.Output{
		Content:        output,
		ShowTimeHeader: false,
		WithHeader:     true,
	}

	out.DumpLogToChannel(h.grpc.OutputChannel)

}

func (h *handler) OnResolveMethod(md *desc.MethodDescriptor) {
	h.methodCount++
	h.method = md
}

func (h *handler) OnSendHeaders(md metadata.MD) {
	h.reqHeadersCount++
	h.reqHeaders = md
}

func (h *handler) OnReceiveHeaders(md metadata.MD) {
	h.respHeadersCount++
	h.respHeaders = md
}

func (h *handler) OnReceiveTrailers(stat *status.Status, md metadata.MD) {
	h.respTrailersCount++
	h.respTrailers = md
	h.respStatus = stat
}
