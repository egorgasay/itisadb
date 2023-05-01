// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: api/proto/storage.proto

package storage

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

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageClient interface {
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	SetToIndex(ctx context.Context, in *SetToIndexRequest, opts ...grpc.CallOption) (*SetResponse, error)
	GetFromIndex(ctx context.Context, in *GetFromIndexRequest, opts ...grpc.CallOption) (*GetResponse, error)
	GetIndex(ctx context.Context, in *GetIndexRequest, opts ...grpc.CallOption) (*GetIndexResponse, error)
	IsIndex(ctx context.Context, in *IsIndexRequest, opts ...grpc.CallOption) (*IsIndexResponse, error)
	NewIndex(ctx context.Context, in *NewIndexRequest, opts ...grpc.CallOption) (*NewIndexResponse, error)
	Size(ctx context.Context, in *IndexSizeRequest, opts ...grpc.CallOption) (*IndexSizeResponse, error)
}

type storageClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageClient(cc grpc.ClientConnInterface) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) SetToIndex(ctx context.Context, in *SetToIndexRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/SetToIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) GetFromIndex(ctx context.Context, in *GetFromIndexRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/GetFromIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) GetIndex(ctx context.Context, in *GetIndexRequest, opts ...grpc.CallOption) (*GetIndexResponse, error) {
	out := new(GetIndexResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/GetIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) IsIndex(ctx context.Context, in *IsIndexRequest, opts ...grpc.CallOption) (*IsIndexResponse, error) {
	out := new(IsIndexResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/IsIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) NewIndex(ctx context.Context, in *NewIndexRequest, opts ...grpc.CallOption) (*NewIndexResponse, error) {
	out := new(NewIndexResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/NewIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) Size(ctx context.Context, in *IndexSizeRequest, opts ...grpc.CallOption) (*IndexSizeResponse, error) {
	out := new(IndexSizeResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/Size", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageServer is the server API for Storage service.
// All implementations must embed UnimplementedStorageServer
// for forward compatibility
type StorageServer interface {
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	SetToIndex(context.Context, *SetToIndexRequest) (*SetResponse, error)
	GetFromIndex(context.Context, *GetFromIndexRequest) (*GetResponse, error)
	GetIndex(context.Context, *GetIndexRequest) (*GetIndexResponse, error)
	IsIndex(context.Context, *IsIndexRequest) (*IsIndexResponse, error)
	NewIndex(context.Context, *NewIndexRequest) (*NewIndexResponse, error)
	Size(context.Context, *IndexSizeRequest) (*IndexSizeResponse, error)
	mustEmbedUnimplementedStorageServer()
}

// UnimplementedStorageServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServer struct {
}

func (UnimplementedStorageServer) Set(context.Context, *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedStorageServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedStorageServer) SetToIndex(context.Context, *SetToIndexRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetToIndex not implemented")
}
func (UnimplementedStorageServer) GetFromIndex(context.Context, *GetFromIndexRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFromIndex not implemented")
}
func (UnimplementedStorageServer) GetIndex(context.Context, *GetIndexRequest) (*GetIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIndex not implemented")
}
func (UnimplementedStorageServer) IsIndex(context.Context, *IsIndexRequest) (*IsIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsIndex not implemented")
}
func (UnimplementedStorageServer) NewIndex(context.Context, *NewIndexRequest) (*NewIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewIndex not implemented")
}
func (UnimplementedStorageServer) Size(context.Context, *IndexSizeRequest) (*IndexSizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Size not implemented")
}
func (UnimplementedStorageServer) mustEmbedUnimplementedStorageServer() {}

// UnsafeStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServer will
// result in compilation errors.
type UnsafeStorageServer interface {
	mustEmbedUnimplementedStorageServer()
}

func RegisterStorageServer(s grpc.ServiceRegistrar, srv StorageServer) {
	s.RegisterService(&Storage_ServiceDesc, srv)
}

func _Storage_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_SetToIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetToIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).SetToIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/SetToIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).SetToIndex(ctx, req.(*SetToIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_GetFromIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFromIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).GetFromIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/GetFromIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).GetFromIndex(ctx, req.(*GetFromIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_GetIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).GetIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/GetIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).GetIndex(ctx, req.(*GetIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_IsIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).IsIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/IsIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).IsIndex(ctx, req.(*IsIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_NewIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).NewIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/NewIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).NewIndex(ctx, req.(*NewIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_Size_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexSizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Size(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/Size",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Size(ctx, req.(*IndexSizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Storage_ServiceDesc is the grpc.ServiceDesc for Storage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Storage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _Storage_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Storage_Get_Handler,
		},
		{
			MethodName: "SetToIndex",
			Handler:    _Storage_SetToIndex_Handler,
		},
		{
			MethodName: "GetFromIndex",
			Handler:    _Storage_GetFromIndex_Handler,
		},
		{
			MethodName: "GetIndex",
			Handler:    _Storage_GetIndex_Handler,
		},
		{
			MethodName: "IsIndex",
			Handler:    _Storage_IsIndex_Handler,
		},
		{
			MethodName: "NewIndex",
			Handler:    _Storage_NewIndex_Handler,
		},
		{
			MethodName: "Size",
			Handler:    _Storage_Size_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/storage.proto",
}
