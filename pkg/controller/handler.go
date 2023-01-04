package controller

import (
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type handler struct {
	controller        Controller
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
	jsm := jsonpb.Marshaler{Indent: "  "}
	respStr, err := jsm.MarshalToString(msg)
	if err != nil {
		panic(fmt.Errorf("failed to generate JSON form of response message: %v", err))
	}
	h.respMessages = append(h.respMessages, respStr)
	h.controller.PrintLog("\nResponse Payload:\n"+respStr, LOG_INFO)
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
