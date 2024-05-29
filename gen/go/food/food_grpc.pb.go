// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: food/food.proto

package food

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

// FoodManagerClient is the client API for FoodManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FoodManagerClient interface {
	GetByRestId(ctx context.Context, in *RestId, opts ...grpc.CallOption) (*RestCategories, error)
	GetById(ctx context.Context, in *FoodId, opts ...grpc.CallOption) (*Food, error)
}

type foodManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewFoodManagerClient(cc grpc.ClientConnInterface) FoodManagerClient {
	return &foodManagerClient{cc}
}

func (c *foodManagerClient) GetByRestId(ctx context.Context, in *RestId, opts ...grpc.CallOption) (*RestCategories, error) {
	out := new(RestCategories)
	err := c.cc.Invoke(ctx, "/food.FoodManager/GetByRestId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *foodManagerClient) GetById(ctx context.Context, in *FoodId, opts ...grpc.CallOption) (*Food, error) {
	out := new(Food)
	err := c.cc.Invoke(ctx, "/food.FoodManager/GetById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FoodManagerServer is the server API for FoodManager service.
// All implementations must embed UnimplementedFoodManagerServer
// for forward compatibility
type FoodManagerServer interface {
	GetByRestId(context.Context, *RestId) (*RestCategories, error)
	GetById(context.Context, *FoodId) (*Food, error)
	mustEmbedUnimplementedFoodManagerServer()
}

// UnimplementedFoodManagerServer must be embedded to have forward compatible implementations.
type UnimplementedFoodManagerServer struct {
}

func (UnimplementedFoodManagerServer) GetByRestId(context.Context, *RestId) (*RestCategories, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByRestId not implemented")
}
func (UnimplementedFoodManagerServer) GetById(context.Context, *FoodId) (*Food, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedFoodManagerServer) mustEmbedUnimplementedFoodManagerServer() {}

// UnsafeFoodManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FoodManagerServer will
// result in compilation errors.
type UnsafeFoodManagerServer interface {
	mustEmbedUnimplementedFoodManagerServer()
}

func RegisterFoodManagerServer(s grpc.ServiceRegistrar, srv FoodManagerServer) {
	s.RegisterService(&FoodManager_ServiceDesc, srv)
}

func _FoodManager_GetByRestId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FoodManagerServer).GetByRestId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/food.FoodManager/GetByRestId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FoodManagerServer).GetByRestId(ctx, req.(*RestId))
	}
	return interceptor(ctx, in, info, handler)
}

func _FoodManager_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FoodId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FoodManagerServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/food.FoodManager/GetById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FoodManagerServer).GetById(ctx, req.(*FoodId))
	}
	return interceptor(ctx, in, info, handler)
}

// FoodManager_ServiceDesc is the grpc.ServiceDesc for FoodManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FoodManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "food.FoodManager",
	HandlerType: (*FoodManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByRestId",
			Handler:    _FoodManager_GetByRestId_Handler,
		},
		{
			MethodName: "GetById",
			Handler:    _FoodManager_GetById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "food/food.proto",
}
