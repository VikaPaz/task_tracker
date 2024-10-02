// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: service.proto

package task

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
	TaskService_GetTasks_FullMethodName         = "/task.TaskService/GetTasks"
	TaskService_UpdateTaskStatus_FullMethodName = "/task.TaskService/UpdateTaskStatus"
)

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskServiceClient interface {
	GetTasks(ctx context.Context, in *GetTasksRequest, opts ...grpc.CallOption) (*GetTasksResponse, error)
	UpdateTaskStatus(ctx context.Context, in *UpdateTaskStatusRequest, opts ...grpc.CallOption) (*UpdateTaskStatusResponse, error)
}

type taskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskServiceClient(cc grpc.ClientConnInterface) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) GetTasks(ctx context.Context, in *GetTasksRequest, opts ...grpc.CallOption) (*GetTasksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTasksResponse)
	err := c.cc.Invoke(ctx, TaskService_GetTasks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) UpdateTaskStatus(ctx context.Context, in *UpdateTaskStatusRequest, opts ...grpc.CallOption) (*UpdateTaskStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateTaskStatusResponse)
	err := c.cc.Invoke(ctx, TaskService_UpdateTaskStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskServiceServer is the server API for TaskService service.
// All implementations must embed UnimplementedTaskServiceServer
// for forward compatibility.
type TaskServiceServer interface {
	GetTasks(context.Context, *GetTasksRequest) (*GetTasksResponse, error)
	UpdateTaskStatus(context.Context, *UpdateTaskStatusRequest) (*UpdateTaskStatusResponse, error)
	mustEmbedUnimplementedTaskServiceServer()
}

// UnimplementedTaskServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTaskServiceServer struct{}

func (UnimplementedTaskServiceServer) GetTasks(context.Context, *GetTasksRequest) (*GetTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTasks not implemented")
}
func (UnimplementedTaskServiceServer) UpdateTaskStatus(context.Context, *UpdateTaskStatusRequest) (*UpdateTaskStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTaskStatus not implemented")
}
func (UnimplementedTaskServiceServer) mustEmbedUnimplementedTaskServiceServer() {}
func (UnimplementedTaskServiceServer) testEmbeddedByValue()                     {}

// UnsafeTaskServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskServiceServer will
// result in compilation errors.
type UnsafeTaskServiceServer interface {
	mustEmbedUnimplementedTaskServiceServer()
}

func RegisterTaskServiceServer(s grpc.ServiceRegistrar, srv TaskServiceServer) {
	// If the following call pancis, it indicates UnimplementedTaskServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TaskService_ServiceDesc, srv)
}

func _TaskService_GetTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTasksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_GetTasks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetTasks(ctx, req.(*GetTasksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_UpdateTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTaskStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).UpdateTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_UpdateTaskStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).UpdateTaskStatus(ctx, req.(*UpdateTaskStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TaskService_ServiceDesc is the grpc.ServiceDesc for TaskService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "task.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTasks",
			Handler:    _TaskService_GetTasks_Handler,
		},
		{
			MethodName: "UpdateTaskStatus",
			Handler:    _TaskService_UpdateTaskStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
