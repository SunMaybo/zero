// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: greeter.proto

package greeter

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

// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreeterClient interface {
	SayHello(ctx context.Context, in *HelloRequest2, opts ...grpc.CallOption) (*HelloReply, error)
	SayHelloWord(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
	SayStream(ctx context.Context, opts ...grpc.CallOption) (Greeter_SayStreamClient, error)
	SayStream1(ctx context.Context, opts ...grpc.CallOption) (Greeter_SayStream1Client, error)
	SayStream2(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (Greeter_SayStream2Client, error)
}

type greeterClient struct {
	cc grpc.ClientConnInterface
}

func NewGreeterClient(cc grpc.ClientConnInterface) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest2, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/greeter.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) SayHelloWord(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/greeter.Greeter/SayHelloWord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) SayStream(ctx context.Context, opts ...grpc.CallOption) (Greeter_SayStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Greeter_ServiceDesc.Streams[0], "/greeter.Greeter/SayStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &greeterSayStreamClient{stream}
	return x, nil
}

type Greeter_SayStreamClient interface {
	Send(*HelloRequest) error
	Recv() (*HelloReply, error)
	grpc.ClientStream
}

type greeterSayStreamClient struct {
	grpc.ClientStream
}

func (x *greeterSayStreamClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *greeterSayStreamClient) Recv() (*HelloReply, error) {
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greeterClient) SayStream1(ctx context.Context, opts ...grpc.CallOption) (Greeter_SayStream1Client, error) {
	stream, err := c.cc.NewStream(ctx, &Greeter_ServiceDesc.Streams[1], "/greeter.Greeter/SayStream1", opts...)
	if err != nil {
		return nil, err
	}
	x := &greeterSayStream1Client{stream}
	return x, nil
}

type Greeter_SayStream1Client interface {
	Send(*HelloRequest) error
	CloseAndRecv() (*HelloReply, error)
	grpc.ClientStream
}

type greeterSayStream1Client struct {
	grpc.ClientStream
}

func (x *greeterSayStream1Client) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *greeterSayStream1Client) CloseAndRecv() (*HelloReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greeterClient) SayStream2(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (Greeter_SayStream2Client, error) {
	stream, err := c.cc.NewStream(ctx, &Greeter_ServiceDesc.Streams[2], "/greeter.Greeter/SayStream2", opts...)
	if err != nil {
		return nil, err
	}
	x := &greeterSayStream2Client{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Greeter_SayStream2Client interface {
	Recv() (*HelloReply, error)
	grpc.ClientStream
}

type greeterSayStream2Client struct {
	grpc.ClientStream
}

func (x *greeterSayStream2Client) Recv() (*HelloReply, error) {
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GreeterServer is the server API for Greeter service.
// All implementations must embed UnimplementedGreeterServer
// for forward compatibility
type GreeterServer interface {
	SayHello(context.Context, *HelloRequest2) (*HelloReply, error)
	SayHelloWord(context.Context, *HelloRequest) (*HelloReply, error)
	SayStream(Greeter_SayStreamServer) error
	SayStream1(Greeter_SayStream1Server) error
	SayStream2(*HelloRequest, Greeter_SayStream2Server) error
	mustEmbedUnimplementedGreeterServer()
}

// UnimplementedGreeterServer must be embedded to have forward compatible implementations.
type UnimplementedGreeterServer struct {
}

func (UnimplementedGreeterServer) SayHello(context.Context, *HelloRequest2) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedGreeterServer) SayHelloWord(context.Context, *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHelloWord not implemented")
}
func (UnimplementedGreeterServer) SayStream(Greeter_SayStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SayStream not implemented")
}
func (UnimplementedGreeterServer) SayStream1(Greeter_SayStream1Server) error {
	return status.Errorf(codes.Unimplemented, "method SayStream1 not implemented")
}
func (UnimplementedGreeterServer) SayStream2(*HelloRequest, Greeter_SayStream2Server) error {
	return status.Errorf(codes.Unimplemented, "method SayStream2 not implemented")
}
func (UnimplementedGreeterServer) mustEmbedUnimplementedGreeterServer() {}

// UnsafeGreeterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GreeterServer will
// result in compilation errors.
type UnsafeGreeterServer interface {
	mustEmbedUnimplementedGreeterServer()
}

func RegisterGreeterServer(s grpc.ServiceRegistrar, srv GreeterServer) {
	s.RegisterService(&Greeter_ServiceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/greeter.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest2))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_SayHelloWord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHelloWord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/greeter.Greeter/SayHelloWord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHelloWord(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_SayStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GreeterServer).SayStream(&greeterSayStreamServer{stream})
}

type Greeter_SayStreamServer interface {
	Send(*HelloReply) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type greeterSayStreamServer struct {
	grpc.ServerStream
}

func (x *greeterSayStreamServer) Send(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *greeterSayStreamServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Greeter_SayStream1_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GreeterServer).SayStream1(&greeterSayStream1Server{stream})
}

type Greeter_SayStream1Server interface {
	SendAndClose(*HelloReply) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type greeterSayStream1Server struct {
	grpc.ServerStream
}

func (x *greeterSayStream1Server) SendAndClose(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *greeterSayStream1Server) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Greeter_SayStream2_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(HelloRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GreeterServer).SayStream2(m, &greeterSayStream2Server{stream})
}

type Greeter_SayStream2Server interface {
	Send(*HelloReply) error
	grpc.ServerStream
}

type greeterSayStream2Server struct {
	grpc.ServerStream
}

func (x *greeterSayStream2Server) Send(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

// Greeter_ServiceDesc is the grpc.ServiceDesc for Greeter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Greeter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "greeter.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
		{
			MethodName: "SayHelloWord",
			Handler:    _Greeter_SayHelloWord_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SayStream",
			Handler:       _Greeter_SayStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "SayStream1",
			Handler:       _Greeter_SayStream1_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SayStream2",
			Handler:       _Greeter_SayStream2_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "greeter.proto",
}