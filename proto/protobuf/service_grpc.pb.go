// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// BannerRotatorServiceClient is the client API for BannerRotatorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BannerRotatorServiceClient interface {
	AddBanner(ctx context.Context, in *AddBannerRequest, opts ...grpc.CallOption) (*AddBannerResponse, error)
	RemoveBanner(ctx context.Context, in *RemoveBannerRequest, opts ...grpc.CallOption) (*RemoveBannerResponse, error)
	HitBanner(ctx context.Context, in *HitBannerRequest, opts ...grpc.CallOption) (*HitBannerResponse, error)
	GetBanner(ctx context.Context, in *GetBannerRequest, opts ...grpc.CallOption) (*GetBannerResponse, error)
}

type bannerRotatorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBannerRotatorServiceClient(cc grpc.ClientConnInterface) BannerRotatorServiceClient {
	return &bannerRotatorServiceClient{cc}
}

func (c *bannerRotatorServiceClient) AddBanner(ctx context.Context, in *AddBannerRequest, opts ...grpc.CallOption) (*AddBannerResponse, error) {
	out := new(AddBannerResponse)
	err := c.cc.Invoke(ctx, "/banner.bannerRotatorService/AddBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bannerRotatorServiceClient) RemoveBanner(ctx context.Context, in *RemoveBannerRequest, opts ...grpc.CallOption) (*RemoveBannerResponse, error) {
	out := new(RemoveBannerResponse)
	err := c.cc.Invoke(ctx, "/banner.bannerRotatorService/RemoveBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bannerRotatorServiceClient) HitBanner(ctx context.Context, in *HitBannerRequest, opts ...grpc.CallOption) (*HitBannerResponse, error) {
	out := new(HitBannerResponse)
	err := c.cc.Invoke(ctx, "/banner.bannerRotatorService/HitBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bannerRotatorServiceClient) GetBanner(ctx context.Context, in *GetBannerRequest, opts ...grpc.CallOption) (*GetBannerResponse, error) {
	out := new(GetBannerResponse)
	err := c.cc.Invoke(ctx, "/banner.bannerRotatorService/GetBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BannerRotatorServiceServer is the server API for BannerRotatorService service.
// All implementations must embed UnimplementedBannerRotatorServiceServer
// for forward compatibility
type BannerRotatorServiceServer interface {
	AddBanner(context.Context, *AddBannerRequest) (*AddBannerResponse, error)
	RemoveBanner(context.Context, *RemoveBannerRequest) (*RemoveBannerResponse, error)
	HitBanner(context.Context, *HitBannerRequest) (*HitBannerResponse, error)
	GetBanner(context.Context, *GetBannerRequest) (*GetBannerResponse, error)
	mustEmbedUnimplementedBannerRotatorServiceServer()
}

// UnimplementedBannerRotatorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBannerRotatorServiceServer struct {
}

func (UnimplementedBannerRotatorServiceServer) AddBanner(context.Context, *AddBannerRequest) (*AddBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBanner not implemented")
}
func (UnimplementedBannerRotatorServiceServer) RemoveBanner(context.Context, *RemoveBannerRequest) (*RemoveBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveBanner not implemented")
}
func (UnimplementedBannerRotatorServiceServer) HitBanner(context.Context, *HitBannerRequest) (*HitBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HitBanner not implemented")
}
func (UnimplementedBannerRotatorServiceServer) GetBanner(context.Context, *GetBannerRequest) (*GetBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBanner not implemented")
}
func (UnimplementedBannerRotatorServiceServer) mustEmbedUnimplementedBannerRotatorServiceServer() {}

// UnsafeBannerRotatorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BannerRotatorServiceServer will
// result in compilation errors.
type UnsafeBannerRotatorServiceServer interface {
	mustEmbedUnimplementedBannerRotatorServiceServer()
}

func RegisterBannerRotatorServiceServer(s grpc.ServiceRegistrar, srv BannerRotatorServiceServer) {
	s.RegisterService(&BannerRotatorService_ServiceDesc, srv)
}

func _BannerRotatorService_AddBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BannerRotatorServiceServer).AddBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/banner.bannerRotatorService/AddBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BannerRotatorServiceServer).AddBanner(ctx, req.(*AddBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BannerRotatorService_RemoveBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BannerRotatorServiceServer).RemoveBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/banner.bannerRotatorService/RemoveBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BannerRotatorServiceServer).RemoveBanner(ctx, req.(*RemoveBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BannerRotatorService_HitBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HitBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BannerRotatorServiceServer).HitBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/banner.bannerRotatorService/HitBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BannerRotatorServiceServer).HitBanner(ctx, req.(*HitBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BannerRotatorService_GetBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BannerRotatorServiceServer).GetBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/banner.bannerRotatorService/GetBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BannerRotatorServiceServer).GetBanner(ctx, req.(*GetBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BannerRotatorService_ServiceDesc is the grpc.ServiceDesc for BannerRotatorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BannerRotatorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "banner.bannerRotatorService",
	HandlerType: (*BannerRotatorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddBanner",
			Handler:    _BannerRotatorService_AddBanner_Handler,
		},
		{
			MethodName: "RemoveBanner",
			Handler:    _BannerRotatorService_RemoveBanner_Handler,
		},
		{
			MethodName: "HitBanner",
			Handler:    _BannerRotatorService_HitBanner_Handler,
		},
		{
			MethodName: "GetBanner",
			Handler:    _BannerRotatorService_GetBanner_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
