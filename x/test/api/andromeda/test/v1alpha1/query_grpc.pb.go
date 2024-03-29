// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: andromeda/test/v1alpha1/query.proto

package testv1alpha1

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

const (
	Query_Asset_FullMethodName  = "/andromeda.test.v1alpha1.Query/Asset"
	Query_Assets_FullMethodName = "/andromeda.test.v1alpha1.Query/Assets"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Asset queries an asset.
	Asset(ctx context.Context, in *QueryAssetRequest, opts ...grpc.CallOption) (*QueryAssetResponse, error)
	// Assets queries all the assets.
	Assets(ctx context.Context, in *QueryAssetsRequest, opts ...grpc.CallOption) (*QueryAssetsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Asset(ctx context.Context, in *QueryAssetRequest, opts ...grpc.CallOption) (*QueryAssetResponse, error) {
	out := new(QueryAssetResponse)
	err := c.cc.Invoke(ctx, Query_Asset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Assets(ctx context.Context, in *QueryAssetsRequest, opts ...grpc.CallOption) (*QueryAssetsResponse, error) {
	out := new(QueryAssetsResponse)
	err := c.cc.Invoke(ctx, Query_Assets_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Asset queries an asset.
	Asset(context.Context, *QueryAssetRequest) (*QueryAssetResponse, error)
	// Assets queries all the assets.
	Assets(context.Context, *QueryAssetsRequest) (*QueryAssetsResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Asset(context.Context, *QueryAssetRequest) (*QueryAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Asset not implemented")
}
func (UnimplementedQueryServer) Assets(context.Context, *QueryAssetsRequest) (*QueryAssetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Assets not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Asset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Asset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Asset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Asset(ctx, req.(*QueryAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Assets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAssetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Assets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Assets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Assets(ctx, req.(*QueryAssetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "andromeda.test.v1alpha1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Asset",
			Handler:    _Query_Asset_Handler,
		},
		{
			MethodName: "Assets",
			Handler:    _Query_Assets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "andromeda/test/v1alpha1/query.proto",
}
