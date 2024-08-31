// Code generated by MockGen. DO NOT EDIT.
// Source: ./webook/internal/repository/cache/user.go
//
// Generated by this command:
//
//	mockgen -source=./webook/internal/repository/cache/user.go -package=cachemocks -destination=./webook/internal/repository/cache/mocks/user.mock.go
//

// Package cachemocks is a generated GoMock package.
package cachemocks

import (
	context "context"
	reflect "reflect"

	domain "gitee.com/geekbang/basic-go/webook/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockUserCache is a mock of UserCache interface.
type MockUserCache struct {
	ctrl     *gomock.Controller
	recorder *MockUserCacheMockRecorder
}

// MockUserCacheMockRecorder is the mock recorder for MockUserCache.
type MockUserCacheMockRecorder struct {
	mock *MockUserCache
}

// NewMockUserCache creates a new mock instance.
func NewMockUserCache(ctrl *gomock.Controller) *MockUserCache {
	mock := &MockUserCache{ctrl: ctrl}
	mock.recorder = &MockUserCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserCache) EXPECT() *MockUserCacheMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, uid)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserCacheMockRecorder) Get(ctx, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserCache)(nil).Get), ctx, uid)
}

// Key mocks base method.
func (m *MockUserCache) Key(uid int64) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Key", uid)
	ret0, _ := ret[0].(string)
	return ret0
}

// Key indicates an expected call of Key.
func (mr *MockUserCacheMockRecorder) Key(uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Key", reflect.TypeOf((*MockUserCache)(nil).Key), uid)
}

// Set mocks base method.
func (m *MockUserCache) Set(ctx context.Context, du domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, du)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockUserCacheMockRecorder) Set(ctx, du any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockUserCache)(nil).Set), ctx, du)
}
