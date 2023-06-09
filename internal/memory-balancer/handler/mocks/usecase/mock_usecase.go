// Code generated by MockGen. DO NOT EDIT.
// Source: itisadb/internal/memory-balancer/handler/grpc (interfaces: IUseCase)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIUseCase is a mock of IUseCase interface.
type MockIUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockIUseCaseMockRecorder
}

// MockIUseCaseMockRecorder is the mock recorder for MockIUseCase.
type MockIUseCaseMockRecorder struct {
	mock *MockIUseCase
}

// NewMockIUseCase creates a new mock instance.
func NewMockIUseCase(ctrl *gomock.Controller) *MockIUseCase {
	mock := &MockIUseCase{ctrl: ctrl}
	mock.recorder = &MockIUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUseCase) EXPECT() *MockIUseCaseMockRecorder {
	return m.recorder
}

// AttachToIndex mocks base method.
func (m *MockIUseCase) AttachToIndex(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachToIndex", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AttachToIndex indicates an expected call of AttachToIndex.
func (mr *MockIUseCaseMockRecorder) AttachToIndex(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachToIndex", reflect.TypeOf((*MockIUseCase)(nil).AttachToIndex), arg0, arg1, arg2)
}

// Connect mocks base method.
func (m *MockIUseCase) Connect(arg0 string, arg1, arg2 uint64, arg3 int32) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connect indicates an expected call of Connect.
func (mr *MockIUseCaseMockRecorder) Connect(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockIUseCase)(nil).Connect), arg0, arg1, arg2, arg3)
}

// Delete mocks base method.
func (m *MockIUseCase) Delete(arg0 context.Context, arg1 string, arg2 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIUseCaseMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIUseCase)(nil).Delete), arg0, arg1, arg2)
}

// DeleteAttr mocks base method.
func (m *MockIUseCase) DeleteAttr(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAttr", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAttr indicates an expected call of DeleteAttr.
func (mr *MockIUseCaseMockRecorder) DeleteAttr(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAttr", reflect.TypeOf((*MockIUseCase)(nil).DeleteAttr), arg0, arg1, arg2)
}

// DeleteIndex mocks base method.
func (m *MockIUseCase) DeleteIndex(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIndex", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteIndex indicates an expected call of DeleteIndex.
func (mr *MockIUseCaseMockRecorder) DeleteIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIndex", reflect.TypeOf((*MockIUseCase)(nil).DeleteIndex), arg0, arg1)
}

// Disconnect mocks base method.
func (m *MockIUseCase) Disconnect(arg0 context.Context, arg1 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockIUseCaseMockRecorder) Disconnect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockIUseCase)(nil).Disconnect), arg0, arg1)
}

// Get mocks base method.
func (m *MockIUseCase) Get(arg0 context.Context, arg1 string, arg2 int32) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIUseCaseMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIUseCase)(nil).Get), arg0, arg1, arg2)
}

// GetFromIndex mocks base method.
func (m *MockIUseCase) GetFromIndex(arg0 context.Context, arg1, arg2 string, arg3 int32) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFromIndex", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFromIndex indicates an expected call of GetFromIndex.
func (mr *MockIUseCaseMockRecorder) GetFromIndex(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFromIndex", reflect.TypeOf((*MockIUseCase)(nil).GetFromIndex), arg0, arg1, arg2, arg3)
}

// IndexToJSON mocks base method.
func (m *MockIUseCase) IndexToJSON(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexToJSON", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexToJSON indicates an expected call of IndexToJSON.
func (mr *MockIUseCaseMockRecorder) IndexToJSON(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexToJSON", reflect.TypeOf((*MockIUseCase)(nil).IndexToJSON), arg0, arg1)
}

// Index mocks base method.
func (m *MockIUseCase) Index(arg0 context.Context, arg1 string) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index", arg0, arg1)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Index indicates an expected call of Index.
func (mr *MockIUseCaseMockRecorder) Index(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockIUseCase)(nil).Index), arg0, arg1)
}

// IsIndex mocks base method.
func (m *MockIUseCase) IsIndex(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsIndex", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsIndex indicates an expected call of IsIndex.
func (mr *MockIUseCaseMockRecorder) IsIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsIndex", reflect.TypeOf((*MockIUseCase)(nil).IsIndex), arg0, arg1)
}

// Servers mocks base method.
func (m *MockIUseCase) Servers() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Servers")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Servers indicates an expected call of Servers.
func (mr *MockIUseCaseMockRecorder) Servers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Servers", reflect.TypeOf((*MockIUseCase)(nil).Servers))
}

// Set mocks base method.
func (m *MockIUseCase) Set(arg0 context.Context, arg1, arg2 string, arg3 int32, arg4 bool) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockIUseCaseMockRecorder) Set(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockIUseCase)(nil).Set), arg0, arg1, arg2, arg3, arg4)
}

// SetToIndex mocks base method.
func (m *MockIUseCase) SetToIndex(arg0 context.Context, arg1, arg2, arg3 string, arg4 bool) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetToIndex", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetToIndex indicates an expected call of SetToIndex.
func (mr *MockIUseCaseMockRecorder) SetToIndex(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetToIndex", reflect.TypeOf((*MockIUseCase)(nil).SetToIndex), arg0, arg1, arg2, arg3, arg4)
}

// Size mocks base method.
func (m *MockIUseCase) Size(arg0 context.Context, arg1 string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Size indicates an expected call of Size.
func (mr *MockIUseCaseMockRecorder) Size(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockIUseCase)(nil).Size), arg0, arg1)
}
