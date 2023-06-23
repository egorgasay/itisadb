// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: api/proto/balancer.proto

package balancer

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

// BalancerClient is the client API for Balancer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BalancerClient interface {
	Index(ctx context.Context, in *BalancerIndexRequest, opts ...grpc.CallOption) (*BalancerIndexResponse, error)
	Set(ctx context.Context, in *BalancerSetRequest, opts ...grpc.CallOption) (*BalancerSetResponse, error)
	SetToIndex(ctx context.Context, in *BalancerSetToIndexRequest, opts ...grpc.CallOption) (*BalancerSetToIndexResponse, error)
	AttachToIndex(ctx context.Context, in *BalancerAttachToIndexRequest, opts ...grpc.CallOption) (*BalancerAttachToIndexResponse, error)
	Get(ctx context.Context, in *BalancerGetRequest, opts ...grpc.CallOption) (*BalancerGetResponse, error)
	GetFromIndex(ctx context.Context, in *BalancerGetFromIndexRequest, opts ...grpc.CallOption) (*BalancerGetFromIndexResponse, error)
	Connect(ctx context.Context, in *BalancerConnectRequest, opts ...grpc.CallOption) (*BalancerConnectResponse, error)
	Disconnect(ctx context.Context, in *BalancerDisconnectRequest, opts ...grpc.CallOption) (*BalancerDisconnectResponse, error)
	Servers(ctx context.Context, in *BalancerServersRequest, opts ...grpc.CallOption) (*BalancerServersResponse, error)
	IndexToJSON(ctx context.Context, in *BalancerIndexToJSONRequest, opts ...grpc.CallOption) (*BalancerIndexToJSONResponse, error)
	IsIndex(ctx context.Context, in *BalancerIsIndexRequest, opts ...grpc.CallOption) (*BalancerIsIndexResponse, error)
	Size(ctx context.Context, in *BalancerIndexSizeRequest, opts ...grpc.CallOption) (*BalancerIndexSizeResponse, error)
	Delete(ctx context.Context, in *BalancerDeleteRequest, opts ...grpc.CallOption) (*BalancerDeleteResponse, error)
	DeleteIfExists(ctx context.Context, in *BalancerDeleteRequest, opts ...grpc.CallOption) (*BalancerDeleteResponse, error)
	DeleteIndex(ctx context.Context, in *BalancerDeleteIndexRequest, opts ...grpc.CallOption) (*BalancerDeleteIndexResponse, error)
	DeleteAttr(ctx context.Context, in *BalancerDeleteAttrRequest, opts ...grpc.CallOption) (*BalancerDeleteAttrResponse, error)
}

type balancerClient struct {
	cc grpc.ClientConnInterface
}

func NewBalancerClient(cc grpc.ClientConnInterface) BalancerClient {
	return &balancerClient{cc}
}

func (c *balancerClient) Index(ctx context.Context, in *BalancerIndexRequest, opts ...grpc.CallOption) (*BalancerIndexResponse, error) {
	out := new(BalancerIndexResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/Index", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) Set(ctx context.Context, in *BalancerSetRequest, opts ...grpc.CallOption) (*BalancerSetResponse, error) {
	out := new(BalancerSetResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) SetToIndex(ctx context.Context, in *BalancerSetToIndexRequest, opts ...grpc.CallOption) (*BalancerSetToIndexResponse, error) {
	out := new(BalancerSetToIndexResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/SetToIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) AttachToIndex(ctx context.Context, in *BalancerAttachToIndexRequest, opts ...grpc.CallOption) (*BalancerAttachToIndexResponse, error) {
	out := new(BalancerAttachToIndexResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/AttachToIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) Get(ctx context.Context, in *BalancerGetRequest, opts ...grpc.CallOption) (*BalancerGetResponse, error) {
	out := new(BalancerGetResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) GetFromIndex(ctx context.Context, in *BalancerGetFromIndexRequest, opts ...grpc.CallOption) (*BalancerGetFromIndexResponse, error) {
	out := new(BalancerGetFromIndexResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/GetFromIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) Connect(ctx context.Context, in *BalancerConnectRequest, opts ...grpc.CallOption) (*BalancerConnectResponse, error) {
	out := new(BalancerConnectResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/Connect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) Disconnect(ctx context.Context, in *BalancerDisconnectRequest, opts ...grpc.CallOption) (*BalancerDisconnectResponse, error) {
	out := new(BalancerDisconnectResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/Disconnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) Servers(ctx context.Context, in *BalancerServersRequest, opts ...grpc.CallOption) (*BalancerServersResponse, error) {
	out := new(BalancerServersResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/Servers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) IndexToJSON(ctx context.Context, in *BalancerIndexToJSONRequest, opts ...grpc.CallOption) (*BalancerIndexToJSONResponse, error) {
	out := new(BalancerIndexToJSONResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/IndexToJSON", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) IsIndex(ctx context.Context, in *BalancerIsIndexRequest, opts ...grpc.CallOption) (*BalancerIsIndexResponse, error) {
	out := new(BalancerIsIndexResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/IsIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) Size(ctx context.Context, in *BalancerIndexSizeRequest, opts ...grpc.CallOption) (*BalancerIndexSizeResponse, error) {
	out := new(BalancerIndexSizeResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/Size", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) Delete(ctx context.Context, in *BalancerDeleteRequest, opts ...grpc.CallOption) (*BalancerDeleteResponse, error) {
	out := new(BalancerDeleteResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) DeleteIfExists(ctx context.Context, in *BalancerDeleteRequest, opts ...grpc.CallOption) (*BalancerDeleteResponse, error) {
	out := new(BalancerDeleteResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/DeleteIfExists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) DeleteIndex(ctx context.Context, in *BalancerDeleteIndexRequest, opts ...grpc.CallOption) (*BalancerDeleteIndexResponse, error) {
	out := new(BalancerDeleteIndexResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/DeleteIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balancerClient) DeleteAttr(ctx context.Context, in *BalancerDeleteAttrRequest, opts ...grpc.CallOption) (*BalancerDeleteAttrResponse, error) {
	out := new(BalancerDeleteAttrResponse)
	err := c.cc.Invoke(ctx, "/api.Balancer/DeleteAttr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BalancerServer is the server API for Balancer service.
// All implementations must embed UnimplementedBalancerServer
// for forward compatibility
type BalancerServer interface {
	Index(context.Context, *BalancerIndexRequest) (*BalancerIndexResponse, error)
	Set(context.Context, *BalancerSetRequest) (*BalancerSetResponse, error)
	SetToIndex(context.Context, *BalancerSetToIndexRequest) (*BalancerSetToIndexResponse, error)
	AttachToIndex(context.Context, *BalancerAttachToIndexRequest) (*BalancerAttachToIndexResponse, error)
	Get(context.Context, *BalancerGetRequest) (*BalancerGetResponse, error)
	GetFromIndex(context.Context, *BalancerGetFromIndexRequest) (*BalancerGetFromIndexResponse, error)
	Connect(context.Context, *BalancerConnectRequest) (*BalancerConnectResponse, error)
	Disconnect(context.Context, *BalancerDisconnectRequest) (*BalancerDisconnectResponse, error)
	Servers(context.Context, *BalancerServersRequest) (*BalancerServersResponse, error)
	IndexToJSON(context.Context, *BalancerIndexToJSONRequest) (*BalancerIndexToJSONResponse, error)
	IsIndex(context.Context, *BalancerIsIndexRequest) (*BalancerIsIndexResponse, error)
	Size(context.Context, *BalancerIndexSizeRequest) (*BalancerIndexSizeResponse, error)
	Delete(context.Context, *BalancerDeleteRequest) (*BalancerDeleteResponse, error)
	DeleteIfExists(context.Context, *BalancerDeleteRequest) (*BalancerDeleteResponse, error)
	DeleteIndex(context.Context, *BalancerDeleteIndexRequest) (*BalancerDeleteIndexResponse, error)
	DeleteAttr(context.Context, *BalancerDeleteAttrRequest) (*BalancerDeleteAttrResponse, error)
	mustEmbedUnimplementedBalancerServer()
}

// UnimplementedBalancerServer must be embedded to have forward compatible implementations.
type UnimplementedBalancerServer struct {
}

func (UnimplementedBalancerServer) Index(context.Context, *BalancerIndexRequest) (*BalancerIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}
func (UnimplementedBalancerServer) Set(context.Context, *BalancerSetRequest) (*BalancerSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedBalancerServer) SetToIndex(context.Context, *BalancerSetToIndexRequest) (*BalancerSetToIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetToIndex not implemented")
}
func (UnimplementedBalancerServer) AttachToIndex(context.Context, *BalancerAttachToIndexRequest) (*BalancerAttachToIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AttachToIndex not implemented")
}
func (UnimplementedBalancerServer) Get(context.Context, *BalancerGetRequest) (*BalancerGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedBalancerServer) GetFromIndex(context.Context, *BalancerGetFromIndexRequest) (*BalancerGetFromIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFromIndex not implemented")
}
func (UnimplementedBalancerServer) Connect(context.Context, *BalancerConnectRequest) (*BalancerConnectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedBalancerServer) Disconnect(context.Context, *BalancerDisconnectRequest) (*BalancerDisconnectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disconnect not implemented")
}
func (UnimplementedBalancerServer) Servers(context.Context, *BalancerServersRequest) (*BalancerServersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Servers not implemented")
}
func (UnimplementedBalancerServer) IndexToJSON(context.Context, *BalancerIndexToJSONRequest) (*BalancerIndexToJSONResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexToJSON not implemented")
}
func (UnimplementedBalancerServer) IsIndex(context.Context, *BalancerIsIndexRequest) (*BalancerIsIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsIndex not implemented")
}
func (UnimplementedBalancerServer) Size(context.Context, *BalancerIndexSizeRequest) (*BalancerIndexSizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Size not implemented")
}
func (UnimplementedBalancerServer) Delete(context.Context, *BalancerDeleteRequest) (*BalancerDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedBalancerServer) DeleteIfExists(context.Context, *BalancerDeleteRequest) (*BalancerDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteIfExists not implemented")
}
func (UnimplementedBalancerServer) DeleteIndex(context.Context, *BalancerDeleteIndexRequest) (*BalancerDeleteIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteIndex not implemented")
}
func (UnimplementedBalancerServer) DeleteAttr(context.Context, *BalancerDeleteAttrRequest) (*BalancerDeleteAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAttr not implemented")
}
func (UnimplementedBalancerServer) mustEmbedUnimplementedBalancerServer() {}

// UnsafeBalancerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BalancerServer will
// result in compilation errors.
type UnsafeBalancerServer interface {
	mustEmbedUnimplementedBalancerServer()
}

func RegisterBalancerServer(s grpc.ServiceRegistrar, srv BalancerServer) {
	s.RegisterService(&Balancer_ServiceDesc, srv)
}

func _Balancer_Index_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).Index(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/Index",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).Index(ctx, req.(*BalancerIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).Set(ctx, req.(*BalancerSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_SetToIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerSetToIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).SetToIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/SetToIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).SetToIndex(ctx, req.(*BalancerSetToIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_AttachToIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerAttachToIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).AttachToIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/AttachToIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).AttachToIndex(ctx, req.(*BalancerAttachToIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).Get(ctx, req.(*BalancerGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_GetFromIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerGetFromIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).GetFromIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/GetFromIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).GetFromIndex(ctx, req.(*BalancerGetFromIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/Connect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).Connect(ctx, req.(*BalancerConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_Disconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerDisconnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).Disconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/Disconnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).Disconnect(ctx, req.(*BalancerDisconnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_Servers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerServersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).Servers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/Servers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).Servers(ctx, req.(*BalancerServersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_IndexToJSON_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerIndexToJSONRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).IndexToJSON(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/IndexToJSON",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).IndexToJSON(ctx, req.(*BalancerIndexToJSONRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_IsIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerIsIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).IsIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/IsIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).IsIndex(ctx, req.(*BalancerIsIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_Size_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerIndexSizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).Size(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/Size",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).Size(ctx, req.(*BalancerIndexSizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).Delete(ctx, req.(*BalancerDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_DeleteIfExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).DeleteIfExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/DeleteIfExists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).DeleteIfExists(ctx, req.(*BalancerDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_DeleteIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerDeleteIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).DeleteIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/DeleteIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).DeleteIndex(ctx, req.(*BalancerDeleteIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Balancer_DeleteAttr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancerDeleteAttrRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalancerServer).DeleteAttr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Balancer/DeleteAttr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalancerServer).DeleteAttr(ctx, req.(*BalancerDeleteAttrRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Balancer_ServiceDesc is the grpc.ServiceDesc for Balancer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Balancer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Balancer",
	HandlerType: (*BalancerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Index",
			Handler:    _Balancer_Index_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _Balancer_Set_Handler,
		},
		{
			MethodName: "SetToIndex",
			Handler:    _Balancer_SetToIndex_Handler,
		},
		{
			MethodName: "AttachToIndex",
			Handler:    _Balancer_AttachToIndex_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Balancer_Get_Handler,
		},
		{
			MethodName: "GetFromIndex",
			Handler:    _Balancer_GetFromIndex_Handler,
		},
		{
			MethodName: "Connect",
			Handler:    _Balancer_Connect_Handler,
		},
		{
			MethodName: "Disconnect",
			Handler:    _Balancer_Disconnect_Handler,
		},
		{
			MethodName: "Servers",
			Handler:    _Balancer_Servers_Handler,
		},
		{
			MethodName: "IndexToJSON",
			Handler:    _Balancer_IndexToJSON_Handler,
		},
		{
			MethodName: "IsIndex",
			Handler:    _Balancer_IsIndex_Handler,
		},
		{
			MethodName: "Size",
			Handler:    _Balancer_Size_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Balancer_Delete_Handler,
		},
		{
			MethodName: "DeleteIfExists",
			Handler:    _Balancer_DeleteIfExists_Handler,
		},
		{
			MethodName: "DeleteIndex",
			Handler:    _Balancer_DeleteIndex_Handler,
		},
		{
			MethodName: "DeleteAttr",
			Handler:    _Balancer_DeleteAttr_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/balancer.proto",
}
