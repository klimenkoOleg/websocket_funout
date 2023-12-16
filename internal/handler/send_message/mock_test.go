// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package send_message is a generated GoMock package.
package send_message

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	device "github.com/klimenkoOleg/websocket_funout/internal/device"
)

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogger) Debug(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), args...)
}

// Error mocks base method.
func (m *MockLogger) Error(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), args...)
}

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockStorage) Count() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	return ret0
}

// Count indicates an expected call of Count.
func (mr *MockStorageMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockStorage)(nil).Count))
}

// Delete mocks base method.
func (m *MockStorage) Delete(id uuid.UUID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", id)
}

// Delete indicates an expected call of Delete.
func (mr *MockStorageMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorage)(nil).Delete), id)
}

// IsDeviceRegistered mocks base method.
func (m *MockStorage) IsDeviceRegistered(id uuid.UUID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDeviceRegistered", id)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDeviceRegistered indicates an expected call of IsDeviceRegistered.
func (mr *MockStorageMockRecorder) IsDeviceRegistered(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDeviceRegistered", reflect.TypeOf((*MockStorage)(nil).IsDeviceRegistered), id)
}

// Store mocks base method.
func (m *MockStorage) Store(d *device.Device) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Store", d)
}

// Store indicates an expected call of Store.
func (mr *MockStorageMockRecorder) Store(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockStorage)(nil).Store), d)
}
