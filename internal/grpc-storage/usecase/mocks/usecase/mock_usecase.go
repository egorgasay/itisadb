// Code generated by MockGen. DO NOT EDIT.
// Source: itisadb/internal/grpc-storage/usecase (interfaces: IUseCase)

// Package mocks is a generated GoMock package.
package usecase_mocks

import (
	usecase "itisadb/internal/grpc-storage/usecase"
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
func (m *MockIUseCase) AttachToIndex(arg0, arg1 string) (usecase.RAM, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachToIndex", arg0, arg1)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AttachToIndex indicates an expected call of AttachToIndex.
func (mr *MockIUseCaseMockRecorder) AttachToIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachToIndex", reflect.TypeOf((*MockIUseCase)(nil).AttachToIndex), arg0, arg1)
}

// Delete mocks base method.
func (m *MockIUseCase) Delete(arg0 string) (usecase.RAM, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockIUseCaseMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIUseCase)(nil).Delete), arg0)
}

// DeleteAttr mocks base method.
func (m *MockIUseCase) DeleteAttr(arg0, arg1 string) (usecase.RAM, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAttr", arg0, arg1)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAttr indicates an expected call of DeleteAttr.
func (mr *MockIUseCaseMockRecorder) DeleteAttr(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAttr", reflect.TypeOf((*MockIUseCase)(nil).DeleteAttr), arg0, arg1)
}

// DeleteIfExists mocks base method.
func (m *MockIUseCase) DeleteIfExists(arg0 string) usecase.RAM {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIfExists", arg0)
	ret0, _ := ret[0].(usecase.RAM)
	return ret0
}

// DeleteIfExists indicates an expected call of DeleteIfExists.
func (mr *MockIUseCaseMockRecorder) DeleteIfExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIfExists", reflect.TypeOf((*MockIUseCase)(nil).DeleteIfExists), arg0)
}

// DeleteIndex mocks base method.
func (m *MockIUseCase) DeleteIndex(arg0 string) (usecase.RAM, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIndex", arg0)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteIndex indicates an expected call of DeleteIndex.
func (mr *MockIUseCaseMockRecorder) DeleteIndex(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIndex", reflect.TypeOf((*MockIUseCase)(nil).DeleteIndex), arg0)
}

// Get mocks base method.
func (m *MockIUseCase) Get(arg0 string) (usecase.RAM, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockIUseCaseMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIUseCase)(nil).Get), arg0)
}

// GetFromIndex mocks base method.
func (m *MockIUseCase) GetFromIndex(arg0, arg1 string) (usecase.RAM, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFromIndex", arg0, arg1)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFromIndex indicates an expected call of GetFromIndex.
func (mr *MockIUseCaseMockRecorder) GetFromIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFromIndex", reflect.TypeOf((*MockIUseCase)(nil).GetFromIndex), arg0, arg1)
}

// GetIndex mocks base method.
func (m *MockIUseCase) GetIndex(arg0 string) (usecase.RAM, map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIndex", arg0)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(map[string]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetIndex indicates an expected call of GetIndex.
func (mr *MockIUseCaseMockRecorder) GetIndex(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIndex", reflect.TypeOf((*MockIUseCase)(nil).GetIndex), arg0)
}

// NewIndex mocks base method.
func (m *MockIUseCase) NewIndex(arg0 string) (usecase.RAM, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIndex", arg0)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewIndex indicates an expected call of NewIndex.
func (mr *MockIUseCaseMockRecorder) NewIndex(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIndex", reflect.TypeOf((*MockIUseCase)(nil).NewIndex), arg0)
}

// Save mocks base method.
func (m *MockIUseCase) Save() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Save")
}

// Save indicates an expected call of Save.
func (mr *MockIUseCaseMockRecorder) Save() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIUseCase)(nil).Save))
}

// Set mocks base method.
func (m *MockIUseCase) Set(arg0, arg1 string, arg2 bool) (usecase.RAM, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockIUseCaseMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockIUseCase)(nil).Set), arg0, arg1, arg2)
}

// SetToIndex mocks base method.
func (m *MockIUseCase) SetToIndex(arg0, arg1, arg2 string, arg3 bool) (usecase.RAM, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetToIndex", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetToIndex indicates an expected call of SetToIndex.
func (mr *MockIUseCaseMockRecorder) SetToIndex(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetToIndex", reflect.TypeOf((*MockIUseCase)(nil).SetToIndex), arg0, arg1, arg2, arg3)
}

// Size mocks base method.
func (m *MockIUseCase) Size(arg0 string) (usecase.RAM, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size", arg0)
	ret0, _ := ret[0].(usecase.RAM)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Size indicates an expected call of Size.
func (mr *MockIUseCaseMockRecorder) Size(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockIUseCase)(nil).Size), arg0)
}