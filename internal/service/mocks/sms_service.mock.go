// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/sms/sms_service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/sms/sms_service.go -destination=internal/service/mocks/sms_service.mock.go -package=svcmock
//
// Package svcmock is a generated GoMock package.
package svcmock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockService) Send(ctx context.Context, tplId string, args, phoneNums []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ctx, tplId, args, phoneNums)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockServiceMockRecorder) Send(ctx, tplId, args, phoneNums any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockService)(nil).Send), ctx, tplId, args, phoneNums)
}
