// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package ticketing

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// CustomerControllerClient is the client API for CustomerController service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CustomerControllerClient interface {
	GetCustomer(ctx context.Context, in *GetCustomerRequest, opts ...grpc.CallOption) (*Customer, error)
	CreateCustomer(ctx context.Context, in *CreateCustomerRequest, opts ...grpc.CallOption) (*Customer, error)
	ListCustomers(ctx context.Context, in *ListCustomersRequest, opts ...grpc.CallOption) (*CustomerList, error)
}

type customerControllerClient struct {
	cc grpc.ClientConnInterface
}

func NewCustomerControllerClient(cc grpc.ClientConnInterface) CustomerControllerClient {
	return &customerControllerClient{cc}
}

func (c *customerControllerClient) GetCustomer(ctx context.Context, in *GetCustomerRequest, opts ...grpc.CallOption) (*Customer, error) {
	out := new(Customer)
	err := c.cc.Invoke(ctx, "/ticketing.CustomerController/GetCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerControllerClient) CreateCustomer(ctx context.Context, in *CreateCustomerRequest, opts ...grpc.CallOption) (*Customer, error) {
	out := new(Customer)
	err := c.cc.Invoke(ctx, "/ticketing.CustomerController/CreateCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerControllerClient) ListCustomers(ctx context.Context, in *ListCustomersRequest, opts ...grpc.CallOption) (*CustomerList, error) {
	out := new(CustomerList)
	err := c.cc.Invoke(ctx, "/ticketing.CustomerController/ListCustomers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CustomerControllerServer is the server API for CustomerController service.
// All implementations must embed UnimplementedCustomerControllerServer
// for forward compatibility
type CustomerControllerServer interface {
	GetCustomer(context.Context, *GetCustomerRequest) (*Customer, error)
	CreateCustomer(context.Context, *CreateCustomerRequest) (*Customer, error)
	ListCustomers(context.Context, *ListCustomersRequest) (*CustomerList, error)
	mustEmbedUnimplementedCustomerControllerServer()
}

// UnimplementedCustomerControllerServer must be embedded to have forward compatible implementations.
type UnimplementedCustomerControllerServer struct {
}

func (UnimplementedCustomerControllerServer) GetCustomer(context.Context, *GetCustomerRequest) (*Customer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCustomer not implemented")
}
func (UnimplementedCustomerControllerServer) CreateCustomer(context.Context, *CreateCustomerRequest) (*Customer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCustomer not implemented")
}
func (UnimplementedCustomerControllerServer) ListCustomers(context.Context, *ListCustomersRequest) (*CustomerList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCustomers not implemented")
}
func (UnimplementedCustomerControllerServer) mustEmbedUnimplementedCustomerControllerServer() {}

// UnsafeCustomerControllerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CustomerControllerServer will
// result in compilation errors.
type UnsafeCustomerControllerServer interface {
	mustEmbedUnimplementedCustomerControllerServer()
}

func RegisterCustomerControllerServer(s grpc.ServiceRegistrar, srv CustomerControllerServer) {
	s.RegisterService(&_CustomerController_serviceDesc, srv)
}

func _CustomerController_GetCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerControllerServer).GetCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ticketing.CustomerController/GetCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerControllerServer).GetCustomer(ctx, req.(*GetCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomerController_CreateCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerControllerServer).CreateCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ticketing.CustomerController/CreateCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerControllerServer).CreateCustomer(ctx, req.(*CreateCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomerController_ListCustomers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCustomersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerControllerServer).ListCustomers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ticketing.CustomerController/ListCustomers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerControllerServer).ListCustomers(ctx, req.(*ListCustomersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CustomerController_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ticketing.CustomerController",
	HandlerType: (*CustomerControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCustomer",
			Handler:    _CustomerController_GetCustomer_Handler,
		},
		{
			MethodName: "CreateCustomer",
			Handler:    _CustomerController_CreateCustomer_Handler,
		},
		{
			MethodName: "ListCustomers",
			Handler:    _CustomerController_ListCustomers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ticketing/ticketing.proto",
}