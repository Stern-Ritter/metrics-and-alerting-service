// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: metrics/metricsapi/v1/metrics_service.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	MetricsV1Service_UpdateMetric_FullMethodName       = "/metrics.metricsapi.v1.MetricsV1Service/UpdateMetric"
	MetricsV1Service_UpdateMetricsBatch_FullMethodName = "/metrics.metricsapi.v1.MetricsV1Service/UpdateMetricsBatch"
	MetricsV1Service_GetMetric_FullMethodName          = "/metrics.metricsapi.v1.MetricsV1Service/GetMetric"
	MetricsV1Service_GetMetrics_FullMethodName         = "/metrics.metricsapi.v1.MetricsV1Service/GetMetrics"
	MetricsV1Service_Ping_FullMethodName               = "/metrics.metricsapi.v1.MetricsV1Service/Ping"
)

// MetricsV1ServiceClient is the client API for MetricsV1Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsV1ServiceClient interface {
	UpdateMetric(ctx context.Context, in *MetricsV1ServiceUpdateMetricRequest, opts ...grpc.CallOption) (*MetricsV1ServiceUpdateMetricResponse, error)
	UpdateMetricsBatch(ctx context.Context, in *MetricsV1ServiceUpdateMetricsBatchRequest, opts ...grpc.CallOption) (*MetricsV1ServiceUpdateMetricsBatchResponse, error)
	GetMetric(ctx context.Context, in *MetricsV1ServiceGetMetricRequest, opts ...grpc.CallOption) (*MetricsV1ServiceGetMetricResponse, error)
	GetMetrics(ctx context.Context, in *MetricsV1ServiceGetMetricsRequest, opts ...grpc.CallOption) (*MetricsV1ServiceGetMetricsResponse, error)
	Ping(ctx context.Context, in *MetricsV1ServicePingRequest, opts ...grpc.CallOption) (*MetricsV1ServicePingResponse, error)
}

type metricsV1ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsV1ServiceClient(cc grpc.ClientConnInterface) MetricsV1ServiceClient {
	return &metricsV1ServiceClient{cc}
}

func (c *metricsV1ServiceClient) UpdateMetric(ctx context.Context, in *MetricsV1ServiceUpdateMetricRequest, opts ...grpc.CallOption) (*MetricsV1ServiceUpdateMetricResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MetricsV1ServiceUpdateMetricResponse)
	err := c.cc.Invoke(ctx, MetricsV1Service_UpdateMetric_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsV1ServiceClient) UpdateMetricsBatch(ctx context.Context, in *MetricsV1ServiceUpdateMetricsBatchRequest, opts ...grpc.CallOption) (*MetricsV1ServiceUpdateMetricsBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MetricsV1ServiceUpdateMetricsBatchResponse)
	err := c.cc.Invoke(ctx, MetricsV1Service_UpdateMetricsBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsV1ServiceClient) GetMetric(ctx context.Context, in *MetricsV1ServiceGetMetricRequest, opts ...grpc.CallOption) (*MetricsV1ServiceGetMetricResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MetricsV1ServiceGetMetricResponse)
	err := c.cc.Invoke(ctx, MetricsV1Service_GetMetric_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsV1ServiceClient) GetMetrics(ctx context.Context, in *MetricsV1ServiceGetMetricsRequest, opts ...grpc.CallOption) (*MetricsV1ServiceGetMetricsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MetricsV1ServiceGetMetricsResponse)
	err := c.cc.Invoke(ctx, MetricsV1Service_GetMetrics_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsV1ServiceClient) Ping(ctx context.Context, in *MetricsV1ServicePingRequest, opts ...grpc.CallOption) (*MetricsV1ServicePingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MetricsV1ServicePingResponse)
	err := c.cc.Invoke(ctx, MetricsV1Service_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsV1ServiceServer is the server API for MetricsV1Service service.
// All implementations must embed UnimplementedMetricsV1ServiceServer
// for forward compatibility
type MetricsV1ServiceServer interface {
	UpdateMetric(context.Context, *MetricsV1ServiceUpdateMetricRequest) (*MetricsV1ServiceUpdateMetricResponse, error)
	UpdateMetricsBatch(context.Context, *MetricsV1ServiceUpdateMetricsBatchRequest) (*MetricsV1ServiceUpdateMetricsBatchResponse, error)
	GetMetric(context.Context, *MetricsV1ServiceGetMetricRequest) (*MetricsV1ServiceGetMetricResponse, error)
	GetMetrics(context.Context, *MetricsV1ServiceGetMetricsRequest) (*MetricsV1ServiceGetMetricsResponse, error)
	Ping(context.Context, *MetricsV1ServicePingRequest) (*MetricsV1ServicePingResponse, error)
	mustEmbedUnimplementedMetricsV1ServiceServer()
}

// UnimplementedMetricsV1ServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsV1ServiceServer struct {
}

func (UnimplementedMetricsV1ServiceServer) UpdateMetric(context.Context, *MetricsV1ServiceUpdateMetricRequest) (*MetricsV1ServiceUpdateMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetric not implemented")
}
func (UnimplementedMetricsV1ServiceServer) UpdateMetricsBatch(context.Context, *MetricsV1ServiceUpdateMetricsBatchRequest) (*MetricsV1ServiceUpdateMetricsBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetricsBatch not implemented")
}
func (UnimplementedMetricsV1ServiceServer) GetMetric(context.Context, *MetricsV1ServiceGetMetricRequest) (*MetricsV1ServiceGetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedMetricsV1ServiceServer) GetMetrics(context.Context, *MetricsV1ServiceGetMetricsRequest) (*MetricsV1ServiceGetMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetrics not implemented")
}
func (UnimplementedMetricsV1ServiceServer) Ping(context.Context, *MetricsV1ServicePingRequest) (*MetricsV1ServicePingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedMetricsV1ServiceServer) mustEmbedUnimplementedMetricsV1ServiceServer() {}

// UnsafeMetricsV1ServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsV1ServiceServer will
// result in compilation errors.
type UnsafeMetricsV1ServiceServer interface {
	mustEmbedUnimplementedMetricsV1ServiceServer()
}

func RegisterMetricsV1ServiceServer(s grpc.ServiceRegistrar, srv MetricsV1ServiceServer) {
	s.RegisterService(&MetricsV1Service_ServiceDesc, srv)
}

func _MetricsV1Service_UpdateMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsV1ServiceUpdateMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsV1ServiceServer).UpdateMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsV1Service_UpdateMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsV1ServiceServer).UpdateMetric(ctx, req.(*MetricsV1ServiceUpdateMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsV1Service_UpdateMetricsBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsV1ServiceUpdateMetricsBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsV1ServiceServer).UpdateMetricsBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsV1Service_UpdateMetricsBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsV1ServiceServer).UpdateMetricsBatch(ctx, req.(*MetricsV1ServiceUpdateMetricsBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsV1Service_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsV1ServiceGetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsV1ServiceServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsV1Service_GetMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsV1ServiceServer).GetMetric(ctx, req.(*MetricsV1ServiceGetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsV1Service_GetMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsV1ServiceGetMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsV1ServiceServer).GetMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsV1Service_GetMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsV1ServiceServer).GetMetrics(ctx, req.(*MetricsV1ServiceGetMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsV1Service_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsV1ServicePingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsV1ServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsV1Service_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsV1ServiceServer).Ping(ctx, req.(*MetricsV1ServicePingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricsV1Service_ServiceDesc is the grpc.ServiceDesc for MetricsV1Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsV1Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metrics.metricsapi.v1.MetricsV1Service",
	HandlerType: (*MetricsV1ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateMetric",
			Handler:    _MetricsV1Service_UpdateMetric_Handler,
		},
		{
			MethodName: "UpdateMetricsBatch",
			Handler:    _MetricsV1Service_UpdateMetricsBatch_Handler,
		},
		{
			MethodName: "GetMetric",
			Handler:    _MetricsV1Service_GetMetric_Handler,
		},
		{
			MethodName: "GetMetrics",
			Handler:    _MetricsV1Service_GetMetrics_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _MetricsV1Service_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metrics/metricsapi/v1/metrics_service.proto",
}
