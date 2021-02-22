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

// EventsClient is the client API for Events service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventsClient interface {
	GetEventByID(ctx context.Context, in *Event, opts ...grpc.CallOption) (*GetEventRS, error)
	GetEvents(ctx context.Context, in *Event, opts ...grpc.CallOption) (*GetEventsRS, error)
	UpdateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*UpdateEventRS, error)
	CreateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*CreateEventRS, error)
	DeleteEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*DeleteEventRS, error)
}

type eventsClient struct {
	cc grpc.ClientConnInterface
}

func NewEventsClient(cc grpc.ClientConnInterface) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) GetEventByID(ctx context.Context, in *Event, opts ...grpc.CallOption) (*GetEventRS, error) {
	out := new(GetEventRS)
	err := c.cc.Invoke(ctx, "/event.Events/GetEventByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) GetEvents(ctx context.Context, in *Event, opts ...grpc.CallOption) (*GetEventsRS, error) {
	out := new(GetEventsRS)
	err := c.cc.Invoke(ctx, "/event.Events/GetEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) UpdateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*UpdateEventRS, error) {
	out := new(UpdateEventRS)
	err := c.cc.Invoke(ctx, "/event.Events/UpdateEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) CreateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*CreateEventRS, error) {
	out := new(CreateEventRS)
	err := c.cc.Invoke(ctx, "/event.Events/CreateEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) DeleteEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*DeleteEventRS, error) {
	out := new(DeleteEventRS)
	err := c.cc.Invoke(ctx, "/event.Events/DeleteEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventsServer is the server API for Events service.
// All implementations must embed UnimplementedEventsServer
// for forward compatibility
type EventsServer interface {
	GetEventByID(context.Context, *Event) (*GetEventRS, error)
	GetEvents(context.Context, *Event) (*GetEventsRS, error)
	UpdateEvent(context.Context, *Event) (*UpdateEventRS, error)
	CreateEvent(context.Context, *Event) (*CreateEventRS, error)
	DeleteEvent(context.Context, *Event) (*DeleteEventRS, error)
	mustEmbedUnimplementedEventsServer()
}

// UnimplementedEventsServer must be embedded to have forward compatible implementations.
type UnimplementedEventsServer struct {
}

func (UnimplementedEventsServer) GetEventByID(context.Context, *Event) (*GetEventRS, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventByID not implemented")
}
func (UnimplementedEventsServer) GetEvents(context.Context, *Event) (*GetEventsRS, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvents not implemented")
}
func (UnimplementedEventsServer) UpdateEvent(context.Context, *Event) (*UpdateEventRS, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (UnimplementedEventsServer) CreateEvent(context.Context, *Event) (*CreateEventRS, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (UnimplementedEventsServer) DeleteEvent(context.Context, *Event) (*DeleteEventRS, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (UnimplementedEventsServer) mustEmbedUnimplementedEventsServer() {}

// UnsafeEventsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventsServer will
// result in compilation errors.
type UnsafeEventsServer interface {
	mustEmbedUnimplementedEventsServer()
}

func RegisterEventsServer(s grpc.ServiceRegistrar, srv EventsServer) {
	s.RegisterService(&Events_ServiceDesc, srv)
}

func _Events_GetEventByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).GetEventByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/GetEventByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).GetEventByID(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_GetEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).GetEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/GetEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).GetEvents(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/UpdateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).UpdateEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/CreateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).CreateEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/DeleteEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).DeleteEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

// Events_ServiceDesc is the grpc.ServiceDesc for Events service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Events_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.Events",
	HandlerType: (*EventsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEventByID",
			Handler:    _Events_GetEventByID_Handler,
		},
		{
			MethodName: "GetEvents",
			Handler:    _Events_GetEvents_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _Events_UpdateEvent_Handler,
		},
		{
			MethodName: "CreateEvent",
			Handler:    _Events_CreateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _Events_DeleteEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "event.proto",
}
