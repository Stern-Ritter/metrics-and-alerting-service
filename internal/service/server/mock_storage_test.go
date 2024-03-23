// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/storage/server.go
//
// Generated by this command:
//
//	mockgen -source=./internal/storage/server.go -destination ./internal/service/server/mock_storage_test.go -package server
//

// Package server is a generated GoMock package.
package server

import (
	reflect "reflect"

	metrics "github.com/Stern-Ritter/metrics-and-alerting-service/internal/model/metrics"
	gomock "go.uber.org/mock/gomock"
)

// MockServerStorage is a mock of ServerStorage interface.
type MockServerStorage struct {
	ctrl     *gomock.Controller
	recorder *MockServerStorageMockRecorder
}

// MockServerStorageMockRecorder is the mock recorder for MockServerStorage.
type MockServerStorageMockRecorder struct {
	mock *MockServerStorage
}

// NewMockServerStorage creates a new mock instance.
func NewMockServerStorage(ctrl *gomock.Controller) *MockServerStorage {
	mock := &MockServerStorage{ctrl: ctrl}
	mock.recorder = &MockServerStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerStorage) EXPECT() *MockServerStorageMockRecorder {
	return m.recorder
}

// GetCounterMetric mocks base method.
func (m *MockServerStorage) GetCounterMetric(metricName string) (metrics.CounterMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounterMetric", metricName)
	ret0, _ := ret[0].(metrics.CounterMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounterMetric indicates an expected call of GetCounterMetric.
func (mr *MockServerStorageMockRecorder) GetCounterMetric(metricName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounterMetric", reflect.TypeOf((*MockServerStorage)(nil).GetCounterMetric), metricName)
}

// GetGaugeMetric mocks base method.
func (m *MockServerStorage) GetGaugeMetric(metricName string) (metrics.GaugeMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGaugeMetric", metricName)
	ret0, _ := ret[0].(metrics.GaugeMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGaugeMetric indicates an expected call of GetGaugeMetric.
func (mr *MockServerStorageMockRecorder) GetGaugeMetric(metricName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGaugeMetric", reflect.TypeOf((*MockServerStorage)(nil).GetGaugeMetric), metricName)
}

// GetMetricValueByTypeAndName mocks base method.
func (m *MockServerStorage) GetMetricValueByTypeAndName(metricType, metricName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricValueByTypeAndName", metricType, metricName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricValueByTypeAndName indicates an expected call of GetMetricValueByTypeAndName.
func (mr *MockServerStorageMockRecorder) GetMetricValueByTypeAndName(metricType, metricName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricValueByTypeAndName", reflect.TypeOf((*MockServerStorage)(nil).GetMetricValueByTypeAndName), metricType, metricName)
}

// GetMetrics mocks base method.
func (m *MockServerStorage) GetMetrics() (map[string]metrics.GaugeMetric, map[string]metrics.CounterMetric) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics")
	ret0, _ := ret[0].(map[string]metrics.GaugeMetric)
	ret1, _ := ret[1].(map[string]metrics.CounterMetric)
	return ret0, ret1
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockServerStorageMockRecorder) GetMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockServerStorage)(nil).GetMetrics))
}

// ResetMetricValue mocks base method.
func (m *MockServerStorage) ResetMetricValue(metricType, metricName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetMetricValue", metricType, metricName)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetMetricValue indicates an expected call of ResetMetricValue.
func (mr *MockServerStorageMockRecorder) ResetMetricValue(metricType, metricName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetMetricValue", reflect.TypeOf((*MockServerStorage)(nil).ResetMetricValue), metricType, metricName)
}

// UpdateCounterMetric mocks base method.
func (m *MockServerStorage) UpdateCounterMetric(metric metrics.CounterMetric) (metrics.CounterMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounterMetric", metric)
	ret0, _ := ret[0].(metrics.CounterMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCounterMetric indicates an expected call of UpdateCounterMetric.
func (mr *MockServerStorageMockRecorder) UpdateCounterMetric(metric any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounterMetric", reflect.TypeOf((*MockServerStorage)(nil).UpdateCounterMetric), metric)
}

// UpdateGaugeMetric mocks base method.
func (m *MockServerStorage) UpdateGaugeMetric(metric metrics.GaugeMetric) (metrics.GaugeMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGaugeMetric", metric)
	ret0, _ := ret[0].(metrics.GaugeMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGaugeMetric indicates an expected call of UpdateGaugeMetric.
func (mr *MockServerStorageMockRecorder) UpdateGaugeMetric(metric any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGaugeMetric", reflect.TypeOf((*MockServerStorage)(nil).UpdateGaugeMetric), metric)
}

// UpdateMetric mocks base method.
func (m *MockServerStorage) UpdateMetric(metricType, metricName, metricValue string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", metricType, metricName, metricValue)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockServerStorageMockRecorder) UpdateMetric(metricType, metricName, metricValue any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockServerStorage)(nil).UpdateMetric), metricType, metricName, metricValue)
}
