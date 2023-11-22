// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/article.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/article.go -destination=internal/service/mocks/article.mock.go -package=svcmock
//
// Package svcmock is a generated GoMock package.
package svcmock

import (
	context "context"
	reflect "reflect"

	domain "github.com/jasonzhao47/cuddle/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockArticleService is a mock of ArticleService interface.
type MockArticleService struct {
	ctrl     *gomock.Controller
	recorder *MockArticleServiceMockRecorder
}

// MockArticleServiceMockRecorder is the mock recorder for MockArticleService.
type MockArticleServiceMockRecorder struct {
	mock *MockArticleService
}

// NewMockArticleService creates a new mock instance.
func NewMockArticleService(ctrl *gomock.Controller) *MockArticleService {
	mock := &MockArticleService{ctrl: ctrl}
	mock.recorder = &MockArticleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleService) EXPECT() *MockArticleServiceMockRecorder {
	return m.recorder
}

// GetById mocks base method.
func (m *MockArticleService) GetById(arg0 context.Context, arg1 int64) (*domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", arg0, arg1)
	ret0, _ := ret[0].(*domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockArticleServiceMockRecorder) GetById(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockArticleService)(nil).GetById), arg0, arg1)
}

// List mocks base method.
func (m *MockArticleService) List(ctx context.Context, authorId int64, page, pageSize int) ([]*domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, authorId, page, pageSize)
	ret0, _ := ret[0].([]*domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockArticleServiceMockRecorder) List(ctx, authorId, page, pageSize any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockArticleService)(nil).List), ctx, authorId, page, pageSize)
}

// Publish mocks base method.
func (m *MockArticleService) Publish(arg0 context.Context, arg1 *domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockArticleServiceMockRecorder) Publish(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockArticleService)(nil).Publish), arg0, arg1)
}

// Save mocks base method.
func (m *MockArticleService) Save(arg0 context.Context, arg1 *domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockArticleServiceMockRecorder) Save(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockArticleService)(nil).Save), arg0, arg1)
}
