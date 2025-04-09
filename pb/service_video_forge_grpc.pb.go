// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: service_video_forge.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	VideosForge_RenewAccessToken_FullMethodName = "/pb.VideosForge/RenewAccessToken"
)

// VideosForgeClient is the client API for VideosForge service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VideosForgeClient interface {
	RenewAccessToken(ctx context.Context, in *RenewAccessTokenRequest, opts ...grpc.CallOption) (*RenewAccessTokenResponse, error)
}

type videosForgeClient struct {
	cc grpc.ClientConnInterface
}

func NewVideosForgeClient(cc grpc.ClientConnInterface) VideosForgeClient {
	return &videosForgeClient{cc}
}

func (c *videosForgeClient) RenewAccessToken(ctx context.Context, in *RenewAccessTokenRequest, opts ...grpc.CallOption) (*RenewAccessTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RenewAccessTokenResponse)
	err := c.cc.Invoke(ctx, VideosForge_RenewAccessToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VideosForgeServer is the server API for VideosForge service.
// All implementations must embed UnimplementedVideosForgeServer
// for forward compatibility.
type VideosForgeServer interface {
	RenewAccessToken(context.Context, *RenewAccessTokenRequest) (*RenewAccessTokenResponse, error)
	mustEmbedUnimplementedVideosForgeServer()
}

// UnimplementedVideosForgeServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedVideosForgeServer struct{}

func (UnimplementedVideosForgeServer) RenewAccessToken(context.Context, *RenewAccessTokenRequest) (*RenewAccessTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenewAccessToken not implemented")
}
func (UnimplementedVideosForgeServer) mustEmbedUnimplementedVideosForgeServer() {}
func (UnimplementedVideosForgeServer) testEmbeddedByValue()                     {}

// UnsafeVideosForgeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VideosForgeServer will
// result in compilation errors.
type UnsafeVideosForgeServer interface {
	mustEmbedUnimplementedVideosForgeServer()
}

func RegisterVideosForgeServer(s grpc.ServiceRegistrar, srv VideosForgeServer) {
	// If the following call pancis, it indicates UnimplementedVideosForgeServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&VideosForge_ServiceDesc, srv)
}

func _VideosForge_RenewAccessToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenewAccessTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideosForgeServer).RenewAccessToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VideosForge_RenewAccessToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideosForgeServer).RenewAccessToken(ctx, req.(*RenewAccessTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VideosForge_ServiceDesc is the grpc.ServiceDesc for VideosForge service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VideosForge_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.VideosForge",
	HandlerType: (*VideosForgeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RenewAccessToken",
			Handler:    _VideosForge_RenewAccessToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_video_forge.proto",
}
