// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.2
// source: proto/processor.proto

package processing_orchestrator

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ProcessorClient is the client API for Processor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProcessorClient interface {
	Start(ctx context.Context, in *ProcessRequest, opts ...grpc.CallOption) (*ProcessResponse, error)
	Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (Processor_WatchClient, error)
	UpdateInfo(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
}

type processorClient struct {
	cc grpc.ClientConnInterface
}

func NewProcessorClient(cc grpc.ClientConnInterface) ProcessorClient {
	return &processorClient{cc}
}

func (c *processorClient) Start(ctx context.Context, in *ProcessRequest, opts ...grpc.CallOption) (*ProcessResponse, error) {
	out := new(ProcessResponse)
	err := c.cc.Invoke(ctx, "/processing_orchestrator.Processor/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processorClient) Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (Processor_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Processor_ServiceDesc.Streams[0], "/processing_orchestrator.Processor/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &processorWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Processor_WatchClient interface {
	Recv() (*ProcessingStatus, error)
	grpc.ClientStream
}

type processorWatchClient struct {
	grpc.ClientStream
}

func (x *processorWatchClient) Recv() (*ProcessingStatus, error) {
	m := new(ProcessingStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *processorClient) UpdateInfo(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/processing_orchestrator.Processor/UpdateInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessorServer is the server API for Processor service.
// All implementations must embed UnimplementedProcessorServer
// for forward compatibility
type ProcessorServer interface {
	Start(context.Context, *ProcessRequest) (*ProcessResponse, error)
	Watch(*WatchRequest, Processor_WatchServer) error
	UpdateInfo(context.Context, *UpdateRequest) (*UpdateResponse, error)
	mustEmbedUnimplementedProcessorServer()
}

// UnimplementedProcessorServer must be embedded to have forward compatible implementations.
type UnimplementedProcessorServer struct {
}

func (UnimplementedProcessorServer) Start(context.Context, *ProcessRequest) (*ProcessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (UnimplementedProcessorServer) Watch(*WatchRequest, Processor_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedProcessorServer) UpdateInfo(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateInfo not implemented")
}
func (UnimplementedProcessorServer) mustEmbedUnimplementedProcessorServer() {}

// UnsafeProcessorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProcessorServer will
// result in compilation errors.
type UnsafeProcessorServer interface {
	mustEmbedUnimplementedProcessorServer()
}

func RegisterProcessorServer(s grpc.ServiceRegistrar, srv ProcessorServer) {
	s.RegisterService(&Processor_ServiceDesc, srv)
}

func _Processor_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessorServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/processing_orchestrator.Processor/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessorServer).Start(ctx, req.(*ProcessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Processor_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ProcessorServer).Watch(m, &processorWatchServer{stream})
}

type Processor_WatchServer interface {
	Send(*ProcessingStatus) error
	grpc.ServerStream
}

type processorWatchServer struct {
	grpc.ServerStream
}

func (x *processorWatchServer) Send(m *ProcessingStatus) error {
	return x.ServerStream.SendMsg(m)
}

func _Processor_UpdateInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessorServer).UpdateInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/processing_orchestrator.Processor/UpdateInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessorServer).UpdateInfo(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Processor_ServiceDesc is the grpc.ServiceDesc for Processor service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Processor_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "processing_orchestrator.Processor",
	HandlerType: (*ProcessorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Start",
			Handler:    _Processor_Start_Handler,
		},
		{
			MethodName: "UpdateInfo",
			Handler:    _Processor_UpdateInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _Processor_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/processor.proto",
}
