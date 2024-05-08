// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: user.proto

package userv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	UserManager_GetData_FullMethodName                 = "/user.UserManager/GetData"
	UserManager_UpdateAddressByUnauthId_FullMethodName = "/user.UserManager/UpdateAddressByUnauthId"
	UserManager_GetAddressByUnauthId_FullMethodName    = "/user.UserManager/GetAddressByUnauthId"
	UserManager_UpdateData_FullMethodName              = "/user.UserManager/UpdateData"
	UserManager_UpdateAddress_FullMethodName           = "/user.UserManager/UpdateAddress"
	UserManager_SetNewPassword_FullMethodName          = "/user.UserManager/SetNewPassword"
	UserManager_Create_FullMethodName                  = "/user.UserManager/Create"
)

// UserManagerClient is the client API for UserManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserManagerClient interface {
	GetData(ctx context.Context, in *Email, opts ...grpc.CallOption) (*User, error)
	UpdateAddressByUnauthId(ctx context.Context, in *AddressDataUnauth, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetAddressByUnauthId(ctx context.Context, in *UnauthId, opts ...grpc.CallOption) (*Address, error)
	UpdateData(ctx context.Context, in *UpdateUserData, opts ...grpc.CallOption) (*User, error)
	UpdateAddress(ctx context.Context, in *AddressData, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetNewPassword(ctx context.Context, in *PasswordsChange, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Create(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
}

type userManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewUserManagerClient(cc grpc.ClientConnInterface) UserManagerClient {
	return &userManagerClient{cc}
}

func (c *userManagerClient) GetData(ctx context.Context, in *Email, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserManager_GetData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userManagerClient) UpdateAddressByUnauthId(ctx context.Context, in *AddressDataUnauth, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, UserManager_UpdateAddressByUnauthId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userManagerClient) GetAddressByUnauthId(ctx context.Context, in *UnauthId, opts ...grpc.CallOption) (*Address, error) {
	out := new(Address)
	err := c.cc.Invoke(ctx, UserManager_GetAddressByUnauthId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userManagerClient) UpdateData(ctx context.Context, in *UpdateUserData, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserManager_UpdateData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userManagerClient) UpdateAddress(ctx context.Context, in *AddressData, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, UserManager_UpdateAddress_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userManagerClient) SetNewPassword(ctx context.Context, in *PasswordsChange, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, UserManager_SetNewPassword_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userManagerClient) Create(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserManager_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserManagerServer is the server API for UserManager service.
// All implementations must embed UnimplementedUserManagerServer
// for forward compatibility
type UserManagerServer interface {
	GetData(context.Context, *Email) (*User, error)
	UpdateAddressByUnauthId(context.Context, *AddressDataUnauth) (*emptypb.Empty, error)
	GetAddressByUnauthId(context.Context, *UnauthId) (*Address, error)
	UpdateData(context.Context, *UpdateUserData) (*User, error)
	UpdateAddress(context.Context, *AddressData) (*emptypb.Empty, error)
	SetNewPassword(context.Context, *PasswordsChange) (*emptypb.Empty, error)
	Create(context.Context, *User) (*User, error)
	mustEmbedUnimplementedUserManagerServer()
}

// UnimplementedUserManagerServer must be embedded to have forward compatible implementations.
type UnimplementedUserManagerServer struct {
}

func (UnimplementedUserManagerServer) GetData(context.Context, *Email) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetData not implemented")
}
func (UnimplementedUserManagerServer) UpdateAddressByUnauthId(context.Context, *AddressDataUnauth) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAddressByUnauthId not implemented")
}
func (UnimplementedUserManagerServer) GetAddressByUnauthId(context.Context, *UnauthId) (*Address, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddressByUnauthId not implemented")
}
func (UnimplementedUserManagerServer) UpdateData(context.Context, *UpdateUserData) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateData not implemented")
}
func (UnimplementedUserManagerServer) UpdateAddress(context.Context, *AddressData) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAddress not implemented")
}
func (UnimplementedUserManagerServer) SetNewPassword(context.Context, *PasswordsChange) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetNewPassword not implemented")
}
func (UnimplementedUserManagerServer) Create(context.Context, *User) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUserManagerServer) mustEmbedUnimplementedUserManagerServer() {}

// UnsafeUserManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserManagerServer will
// result in compilation errors.
type UnsafeUserManagerServer interface {
	mustEmbedUnimplementedUserManagerServer()
}

func RegisterUserManagerServer(s grpc.ServiceRegistrar, srv UserManagerServer) {
	s.RegisterService(&UserManager_ServiceDesc, srv)
}

func _UserManager_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagerServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManager_GetData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagerServer).GetData(ctx, req.(*Email))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserManager_UpdateAddressByUnauthId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddressDataUnauth)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagerServer).UpdateAddressByUnauthId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManager_UpdateAddressByUnauthId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagerServer).UpdateAddressByUnauthId(ctx, req.(*AddressDataUnauth))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserManager_GetAddressByUnauthId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnauthId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagerServer).GetAddressByUnauthId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManager_GetAddressByUnauthId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagerServer).GetAddressByUnauthId(ctx, req.(*UnauthId))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserManager_UpdateData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagerServer).UpdateData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManager_UpdateData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagerServer).UpdateData(ctx, req.(*UpdateUserData))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserManager_UpdateAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddressData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagerServer).UpdateAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManager_UpdateAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagerServer).UpdateAddress(ctx, req.(*AddressData))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserManager_SetNewPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordsChange)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagerServer).SetNewPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManager_SetNewPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagerServer).SetNewPassword(ctx, req.(*PasswordsChange))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserManager_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManager_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagerServer).Create(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

// UserManager_ServiceDesc is the grpc.ServiceDesc for UserManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserManager",
	HandlerType: (*UserManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetData",
			Handler:    _UserManager_GetData_Handler,
		},
		{
			MethodName: "UpdateAddressByUnauthId",
			Handler:    _UserManager_UpdateAddressByUnauthId_Handler,
		},
		{
			MethodName: "GetAddressByUnauthId",
			Handler:    _UserManager_GetAddressByUnauthId_Handler,
		},
		{
			MethodName: "UpdateData",
			Handler:    _UserManager_UpdateData_Handler,
		},
		{
			MethodName: "UpdateAddress",
			Handler:    _UserManager_UpdateAddress_Handler,
		},
		{
			MethodName: "SetNewPassword",
			Handler:    _UserManager_SetNewPassword_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _UserManager_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
