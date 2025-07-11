// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: im_protobuf.proto

package protobuf

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
	AccServer_QueryUsersOnline_FullMethodName = "/protobuf.AccServer/QueryUsersOnline"
	AccServer_SendMsg_FullMethodName          = "/protobuf.AccServer/SendMsg"
	AccServer_SendMsgAll_FullMethodName       = "/protobuf.AccServer/SendMsgAll"
	AccServer_GetUserList_FullMethodName      = "/protobuf.AccServer/GetUserList"
)

// AccServerClient is the client API for AccServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The AccServer service definition.
type AccServerClient interface {
	// 查询用户是否在线
	QueryUsersOnline(ctx context.Context, in *QueryUsersOnlineReq, opts ...grpc.CallOption) (*QueryUsersOnlineRsp, error)
	// 发送消息
	SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgRsp, error)
	// 给这台机器的房间内所有用户发送消息
	SendMsgAll(ctx context.Context, in *SendMsgAllReq, opts ...grpc.CallOption) (*SendMsgAllRsp, error)
	// 获取用户列表
	GetUserList(ctx context.Context, in *GetUserListReq, opts ...grpc.CallOption) (*GetUserListRsp, error)
}

type accServerClient struct {
	cc grpc.ClientConnInterface
}

func NewAccServerClient(cc grpc.ClientConnInterface) AccServerClient {
	return &accServerClient{cc}
}

func (c *accServerClient) QueryUsersOnline(ctx context.Context, in *QueryUsersOnlineReq, opts ...grpc.CallOption) (*QueryUsersOnlineRsp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryUsersOnlineRsp)
	err := c.cc.Invoke(ctx, AccServer_QueryUsersOnline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accServerClient) SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgRsp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendMsgRsp)
	err := c.cc.Invoke(ctx, AccServer_SendMsg_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accServerClient) SendMsgAll(ctx context.Context, in *SendMsgAllReq, opts ...grpc.CallOption) (*SendMsgAllRsp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendMsgAllRsp)
	err := c.cc.Invoke(ctx, AccServer_SendMsgAll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accServerClient) GetUserList(ctx context.Context, in *GetUserListReq, opts ...grpc.CallOption) (*GetUserListRsp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserListRsp)
	err := c.cc.Invoke(ctx, AccServer_GetUserList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccServerServer is the server API for AccServer service.
// All implementations must embed UnimplementedAccServerServer
// for forward compatibility.
//
// The AccServer service definition.
type AccServerServer interface {
	// 查询用户是否在线
	QueryUsersOnline(context.Context, *QueryUsersOnlineReq) (*QueryUsersOnlineRsp, error)
	// 发送消息
	SendMsg(context.Context, *SendMsgReq) (*SendMsgRsp, error)
	// 给这台机器的房间内所有用户发送消息
	SendMsgAll(context.Context, *SendMsgAllReq) (*SendMsgAllRsp, error)
	// 获取用户列表
	GetUserList(context.Context, *GetUserListReq) (*GetUserListRsp, error)
	mustEmbedUnimplementedAccServerServer()
}

// UnimplementedAccServerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAccServerServer struct{}

func (UnimplementedAccServerServer) QueryUsersOnline(context.Context, *QueryUsersOnlineReq) (*QueryUsersOnlineRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryUsersOnline not implemented")
}
func (UnimplementedAccServerServer) SendMsg(context.Context, *SendMsgReq) (*SendMsgRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsg not implemented")
}
func (UnimplementedAccServerServer) SendMsgAll(context.Context, *SendMsgAllReq) (*SendMsgAllRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsgAll not implemented")
}
func (UnimplementedAccServerServer) GetUserList(context.Context, *GetUserListReq) (*GetUserListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserList not implemented")
}
func (UnimplementedAccServerServer) mustEmbedUnimplementedAccServerServer() {}
func (UnimplementedAccServerServer) testEmbeddedByValue()                   {}

// UnsafeAccServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccServerServer will
// result in compilation errors.
type UnsafeAccServerServer interface {
	mustEmbedUnimplementedAccServerServer()
}

func RegisterAccServerServer(s grpc.ServiceRegistrar, srv AccServerServer) {
	// If the following call pancis, it indicates UnimplementedAccServerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AccServer_ServiceDesc, srv)
}

func _AccServer_QueryUsersOnline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryUsersOnlineReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccServerServer).QueryUsersOnline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccServer_QueryUsersOnline_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccServerServer).QueryUsersOnline(ctx, req.(*QueryUsersOnlineReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccServer_SendMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccServerServer).SendMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccServer_SendMsg_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccServerServer).SendMsg(ctx, req.(*SendMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccServer_SendMsgAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsgAllReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccServerServer).SendMsgAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccServer_SendMsgAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccServerServer).SendMsgAll(ctx, req.(*SendMsgAllReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccServer_GetUserList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccServerServer).GetUserList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccServer_GetUserList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccServerServer).GetUserList(ctx, req.(*GetUserListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AccServer_ServiceDesc is the grpc.ServiceDesc for AccServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.AccServer",
	HandlerType: (*AccServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryUsersOnline",
			Handler:    _AccServer_QueryUsersOnline_Handler,
		},
		{
			MethodName: "SendMsg",
			Handler:    _AccServer_SendMsg_Handler,
		},
		{
			MethodName: "SendMsgAll",
			Handler:    _AccServer_SendMsgAll_Handler,
		},
		{
			MethodName: "GetUserList",
			Handler:    _AccServer_GetUserList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "im_protobuf.proto",
}
