// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-scalingo (interfaces: SubresourceService)

// Package scalingo is a generated GoMock package.
package scalingo

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSubresourceService is a mock of SubresourceService interface
type MockSubresourceService struct {
	ctrl     *gomock.Controller
	recorder *MockSubresourceServiceMockRecorder
}

// MockSubresourceServiceMockRecorder is the mock recorder for MockSubresourceService
type MockSubresourceServiceMockRecorder struct {
	mock *MockSubresourceService
}

// NewMockSubresourceService creates a new mock instance
func NewMockSubresourceService(ctrl *gomock.Controller) *MockSubresourceService {
	mock := &MockSubresourceService{ctrl: ctrl}
	mock.recorder = &MockSubresourceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSubresourceService) EXPECT() *MockSubresourceServiceMockRecorder {
	return m.recorder
}

// subresourceAdd mocks base method
func (m *MockSubresourceService) subresourceAdd(arg0, arg1 string, arg2, arg3 interface{}) error {
	ret := m.ctrl.Call(m, "subresourceAdd", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// subresourceAdd indicates an expected call of subresourceAdd
func (mr *MockSubresourceServiceMockRecorder) subresourceAdd(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "subresourceAdd", reflect.TypeOf((*MockSubresourceService)(nil).subresourceAdd), arg0, arg1, arg2, arg3)
}

// subresourceDelete mocks base method
func (m *MockSubresourceService) subresourceDelete(arg0, arg1, arg2 string) error {
	ret := m.ctrl.Call(m, "subresourceDelete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// subresourceDelete indicates an expected call of subresourceDelete
func (mr *MockSubresourceServiceMockRecorder) subresourceDelete(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "subresourceDelete", reflect.TypeOf((*MockSubresourceService)(nil).subresourceDelete), arg0, arg1, arg2)
}

// subresourceGet mocks base method
func (m *MockSubresourceService) subresourceGet(arg0, arg1, arg2 string, arg3, arg4 interface{}) error {
	ret := m.ctrl.Call(m, "subresourceGet", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// subresourceGet indicates an expected call of subresourceGet
func (mr *MockSubresourceServiceMockRecorder) subresourceGet(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "subresourceGet", reflect.TypeOf((*MockSubresourceService)(nil).subresourceGet), arg0, arg1, arg2, arg3, arg4)
}

// subresourceList mocks base method
func (m *MockSubresourceService) subresourceList(arg0, arg1 string, arg2, arg3 interface{}) error {
	ret := m.ctrl.Call(m, "subresourceList", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// subresourceList indicates an expected call of subresourceList
func (mr *MockSubresourceServiceMockRecorder) subresourceList(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "subresourceList", reflect.TypeOf((*MockSubresourceService)(nil).subresourceList), arg0, arg1, arg2, arg3)
}

// subresourceUpdate mocks base method
func (m *MockSubresourceService) subresourceUpdate(arg0, arg1, arg2 string, arg3, arg4 interface{}) error {
	ret := m.ctrl.Call(m, "subresourceUpdate", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// subresourceUpdate indicates an expected call of subresourceUpdate
func (mr *MockSubresourceServiceMockRecorder) subresourceUpdate(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "subresourceUpdate", reflect.TypeOf((*MockSubresourceService)(nil).subresourceUpdate), arg0, arg1, arg2, arg3, arg4)
}