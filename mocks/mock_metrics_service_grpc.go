// Code generated by MockGen. DO NOT EDIT.
// Source: ./proto/gen/metrics/metricsapi/v1/metrics_service_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=./proto/gen/metrics/metricsapi/v1/metrics_service_grpc.pb.go -destination=./mocks/mock_metrics_service_grpc.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	v1 "github.com/Stern-Ritter/metrics-and-alerting-service/proto/gen/metrics/metricsapi/v1"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockMetricsV1ServiceClient is a mock of MetricsV1ServiceClient interface.
type MockMetricsV1ServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsV1ServiceClientMockRecorder
}

// MockMetricsV1ServiceClientMockRecorder is the mock recorder for MockMetricsV1ServiceClient.
type MockMetricsV1ServiceClientMockRecorder struct {
	mock *MockMetricsV1ServiceClient
}

// NewMockMetricsV1ServiceClient creates a new mock instance.
func NewMockMetricsV1ServiceClient(ctrl *gomock.Controller) *MockMetricsV1ServiceClient {
	mock := &MockMetricsV1ServiceClient{ctrl: ctrl}
	mock.recorder = &MockMetricsV1ServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsV1ServiceClient) EXPECT() *MockMetricsV1ServiceClientMockRecorder {
	return m.recorder
}

// GetMetric mocks base method.
func (m *MockMetricsV1ServiceClient) GetMetric(ctx context.Context, in *v1.MetricsV1ServiceGetMetricRequest, opts ...grpc.CallOption) (*v1.MetricsV1ServiceGetMetricResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMetric", varargs...)
	ret0, _ := ret[0].(*v1.MetricsV1ServiceGetMetricResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockMetricsV1ServiceClientMockRecorder) GetMetric(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockMetricsV1ServiceClient)(nil).GetMetric), varargs...)
}

// GetMetrics mocks base method.
func (m *MockMetricsV1ServiceClient) GetMetrics(ctx context.Context, in *v1.MetricsV1ServiceGetMetricsRequest, opts ...grpc.CallOption) (*v1.MetricsV1ServiceGetMetricsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMetrics", varargs...)
	ret0, _ := ret[0].(*v1.MetricsV1ServiceGetMetricsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockMetricsV1ServiceClientMockRecorder) GetMetrics(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockMetricsV1ServiceClient)(nil).GetMetrics), varargs...)
}

// Ping mocks base method.
func (m *MockMetricsV1ServiceClient) Ping(ctx context.Context, in *v1.MetricsV1ServicePingRequest, opts ...grpc.CallOption) (*v1.MetricsV1ServicePingResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Ping", varargs...)
	ret0, _ := ret[0].(*v1.MetricsV1ServicePingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ping indicates an expected call of Ping.
func (mr *MockMetricsV1ServiceClientMockRecorder) Ping(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockMetricsV1ServiceClient)(nil).Ping), varargs...)
}

// UpdateMetric mocks base method.
func (m *MockMetricsV1ServiceClient) UpdateMetric(ctx context.Context, in *v1.MetricsV1ServiceUpdateMetricRequest, opts ...grpc.CallOption) (*v1.MetricsV1ServiceUpdateMetricResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateMetric", varargs...)
	ret0, _ := ret[0].(*v1.MetricsV1ServiceUpdateMetricResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockMetricsV1ServiceClientMockRecorder) UpdateMetric(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockMetricsV1ServiceClient)(nil).UpdateMetric), varargs...)
}

// UpdateMetricsBatch mocks base method.
func (m *MockMetricsV1ServiceClient) UpdateMetricsBatch(ctx context.Context, in *v1.MetricsV1ServiceUpdateMetricsBatchRequest, opts ...grpc.CallOption) (*v1.MetricsV1ServiceUpdateMetricsBatchResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateMetricsBatch", varargs...)
	ret0, _ := ret[0].(*v1.MetricsV1ServiceUpdateMetricsBatchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMetricsBatch indicates an expected call of UpdateMetricsBatch.
func (mr *MockMetricsV1ServiceClientMockRecorder) UpdateMetricsBatch(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricsBatch", reflect.TypeOf((*MockMetricsV1ServiceClient)(nil).UpdateMetricsBatch), varargs...)
}

// MockMetricsV1ServiceServer is a mock of MetricsV1ServiceServer interface.
type MockMetricsV1ServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsV1ServiceServerMockRecorder
}

// MockMetricsV1ServiceServerMockRecorder is the mock recorder for MockMetricsV1ServiceServer.
type MockMetricsV1ServiceServerMockRecorder struct {
	mock *MockMetricsV1ServiceServer
}

// NewMockMetricsV1ServiceServer creates a new mock instance.
func NewMockMetricsV1ServiceServer(ctrl *gomock.Controller) *MockMetricsV1ServiceServer {
	mock := &MockMetricsV1ServiceServer{ctrl: ctrl}
	mock.recorder = &MockMetricsV1ServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsV1ServiceServer) EXPECT() *MockMetricsV1ServiceServerMockRecorder {
	return m.recorder
}

// GetMetric mocks base method.
func (m *MockMetricsV1ServiceServer) GetMetric(arg0 context.Context, arg1 *v1.MetricsV1ServiceGetMetricRequest) (*v1.MetricsV1ServiceGetMetricResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", arg0, arg1)
	ret0, _ := ret[0].(*v1.MetricsV1ServiceGetMetricResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockMetricsV1ServiceServerMockRecorder) GetMetric(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockMetricsV1ServiceServer)(nil).GetMetric), arg0, arg1)
}

// GetMetrics mocks base method.
func (m *MockMetricsV1ServiceServer) GetMetrics(arg0 context.Context, arg1 *v1.MetricsV1ServiceGetMetricsRequest) (*v1.MetricsV1ServiceGetMetricsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics", arg0, arg1)
	ret0, _ := ret[0].(*v1.MetricsV1ServiceGetMetricsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockMetricsV1ServiceServerMockRecorder) GetMetrics(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockMetricsV1ServiceServer)(nil).GetMetrics), arg0, arg1)
}

// Ping mocks base method.
func (m *MockMetricsV1ServiceServer) Ping(arg0 context.Context, arg1 *v1.MetricsV1ServicePingRequest) (*v1.MetricsV1ServicePingResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0, arg1)
	ret0, _ := ret[0].(*v1.MetricsV1ServicePingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ping indicates an expected call of Ping.
func (mr *MockMetricsV1ServiceServerMockRecorder) Ping(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockMetricsV1ServiceServer)(nil).Ping), arg0, arg1)
}

// UpdateMetric mocks base method.
func (m *MockMetricsV1ServiceServer) UpdateMetric(arg0 context.Context, arg1 *v1.MetricsV1ServiceUpdateMetricRequest) (*v1.MetricsV1ServiceUpdateMetricResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", arg0, arg1)
	ret0, _ := ret[0].(*v1.MetricsV1ServiceUpdateMetricResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockMetricsV1ServiceServerMockRecorder) UpdateMetric(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockMetricsV1ServiceServer)(nil).UpdateMetric), arg0, arg1)
}

// UpdateMetricsBatch mocks base method.
func (m *MockMetricsV1ServiceServer) UpdateMetricsBatch(arg0 context.Context, arg1 *v1.MetricsV1ServiceUpdateMetricsBatchRequest) (*v1.MetricsV1ServiceUpdateMetricsBatchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetricsBatch", arg0, arg1)
	ret0, _ := ret[0].(*v1.MetricsV1ServiceUpdateMetricsBatchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMetricsBatch indicates an expected call of UpdateMetricsBatch.
func (mr *MockMetricsV1ServiceServerMockRecorder) UpdateMetricsBatch(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricsBatch", reflect.TypeOf((*MockMetricsV1ServiceServer)(nil).UpdateMetricsBatch), arg0, arg1)
}

// mustEmbedUnimplementedMetricsV1ServiceServer mocks base method.
func (m *MockMetricsV1ServiceServer) mustEmbedUnimplementedMetricsV1ServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedMetricsV1ServiceServer")
}

// mustEmbedUnimplementedMetricsV1ServiceServer indicates an expected call of mustEmbedUnimplementedMetricsV1ServiceServer.
func (mr *MockMetricsV1ServiceServerMockRecorder) mustEmbedUnimplementedMetricsV1ServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedMetricsV1ServiceServer", reflect.TypeOf((*MockMetricsV1ServiceServer)(nil).mustEmbedUnimplementedMetricsV1ServiceServer))
}

// MockUnsafeMetricsV1ServiceServer is a mock of UnsafeMetricsV1ServiceServer interface.
type MockUnsafeMetricsV1ServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeMetricsV1ServiceServerMockRecorder
}

// MockUnsafeMetricsV1ServiceServerMockRecorder is the mock recorder for MockUnsafeMetricsV1ServiceServer.
type MockUnsafeMetricsV1ServiceServerMockRecorder struct {
	mock *MockUnsafeMetricsV1ServiceServer
}

// NewMockUnsafeMetricsV1ServiceServer creates a new mock instance.
func NewMockUnsafeMetricsV1ServiceServer(ctrl *gomock.Controller) *MockUnsafeMetricsV1ServiceServer {
	mock := &MockUnsafeMetricsV1ServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeMetricsV1ServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeMetricsV1ServiceServer) EXPECT() *MockUnsafeMetricsV1ServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedMetricsV1ServiceServer mocks base method.
func (m *MockUnsafeMetricsV1ServiceServer) mustEmbedUnimplementedMetricsV1ServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedMetricsV1ServiceServer")
}

// mustEmbedUnimplementedMetricsV1ServiceServer indicates an expected call of mustEmbedUnimplementedMetricsV1ServiceServer.
func (mr *MockUnsafeMetricsV1ServiceServerMockRecorder) mustEmbedUnimplementedMetricsV1ServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedMetricsV1ServiceServer", reflect.TypeOf((*MockUnsafeMetricsV1ServiceServer)(nil).mustEmbedUnimplementedMetricsV1ServiceServer))
}
