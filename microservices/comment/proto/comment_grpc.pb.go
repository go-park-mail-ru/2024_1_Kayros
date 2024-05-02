// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package comment

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

// CommentWorkerClient is the client API for CommentWorker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommentWorkerClient interface {
	CreateComment(ctx context.Context, in *Comment, opts ...grpc.CallOption) (*Comment, error)
	DeleteComment(ctx context.Context, in *CommentId, opts ...grpc.CallOption) (*Empty, error)
	GetCommentsByRest(ctx context.Context, in *RestId, opts ...grpc.CallOption) (*CommentList, error)
}

type commentWorkerClient struct {
	cc grpc.ClientConnInterface
}

func NewCommentWorkerClient(cc grpc.ClientConnInterface) CommentWorkerClient {
	return &commentWorkerClient{cc}
}

func (c *commentWorkerClient) CreateComment(ctx context.Context, in *Comment, opts ...grpc.CallOption) (*Comment, error) {
	out := new(Comment)
	err := c.cc.Invoke(ctx, "/comment.CommentWorker/CreateComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentWorkerClient) DeleteComment(ctx context.Context, in *CommentId, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/comment.CommentWorker/DeleteComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentWorkerClient) GetCommentsByRest(ctx context.Context, in *RestId, opts ...grpc.CallOption) (*CommentList, error) {
	out := new(CommentList)
	err := c.cc.Invoke(ctx, "/comment.CommentWorker/GetCommentsByRest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommentWorkerServer is the server API for CommentWorker service.
// All implementations must embed UnimplementedCommentWorkerServer
// for forward compatibility
type CommentWorkerServer interface {
	CreateComment(context.Context, *Comment) (*Comment, error)
	DeleteComment(context.Context, *CommentId) (*Empty, error)
	GetCommentsByRest(context.Context, *RestId) (*CommentList, error)
	mustEmbedUnimplementedCommentWorkerServer()
}

// UnimplementedCommentWorkerServer must be embedded to have forward compatible implementations.
type UnimplementedCommentWorkerServer struct {
}

func (UnimplementedCommentWorkerServer) CreateComment(context.Context, *Comment) (*Comment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateComment not implemented")
}
func (UnimplementedCommentWorkerServer) DeleteComment(context.Context, *CommentId) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (UnimplementedCommentWorkerServer) GetCommentsByRest(context.Context, *RestId) (*CommentList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommentsByRest not implemented")
}
func (UnimplementedCommentWorkerServer) mustEmbedUnimplementedCommentWorkerServer() {}

// UnsafeCommentWorkerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommentWorkerServer will
// result in compilation errors.
type UnsafeCommentWorkerServer interface {
	mustEmbedUnimplementedCommentWorkerServer()
}

func RegisterCommentWorkerServer(s grpc.ServiceRegistrar, srv CommentWorkerServer) {
	s.RegisterService(&CommentWorker_ServiceDesc, srv)
}

func _CommentWorker_CreateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Comment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentWorkerServer).CreateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comment.CommentWorker/CreateComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentWorkerServer).CreateComment(ctx, req.(*Comment))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentWorker_DeleteComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommentId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentWorkerServer).DeleteComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comment.CommentWorker/DeleteComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentWorkerServer).DeleteComment(ctx, req.(*CommentId))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentWorker_GetCommentsByRest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentWorkerServer).GetCommentsByRest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comment.CommentWorker/GetCommentsByRest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentWorkerServer).GetCommentsByRest(ctx, req.(*RestId))
	}
	return interceptor(ctx, in, info, handler)
}

// CommentWorker_ServiceDesc is the grpc.ServiceDesc for CommentWorker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommentWorker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "comment.CommentWorker",
	HandlerType: (*CommentWorkerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateComment",
			Handler:    _CommentWorker_CreateComment_Handler,
		},
		{
			MethodName: "DeleteComment",
			Handler:    _CommentWorker_DeleteComment_Handler,
		},
		{
			MethodName: "GetCommentsByRest",
			Handler:    _CommentWorker_GetCommentsByRest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "comment.proto",
}
