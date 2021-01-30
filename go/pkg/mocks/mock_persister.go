// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/remotehack/bottle/pkg/persister (interfaces: Persister)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPersister is a mock of Persister interface
type MockPersister struct {
	ctrl     *gomock.Controller
	recorder *MockPersisterMockRecorder
}

// MockPersisterMockRecorder is the mock recorder for MockPersister
type MockPersisterMockRecorder struct {
	mock *MockPersister
}

// NewMockPersister creates a new mock instance
func NewMockPersister(ctrl *gomock.Controller) *MockPersister {
	mock := &MockPersister{ctrl: ctrl}
	mock.recorder = &MockPersisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPersister) EXPECT() *MockPersisterMockRecorder {
	return m.recorder
}

// Write mocks base method
func (m *MockPersister) Write(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write
func (mr *MockPersisterMockRecorder) Write(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockPersister)(nil).Write), arg0, arg1)
}
