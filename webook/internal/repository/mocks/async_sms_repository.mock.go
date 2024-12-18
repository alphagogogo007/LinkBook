// Code generated by MockGen. DO NOT EDIT.
// Source: ./async_sms_repository.go
//
// Generated by this command:
//
//	mockgen -source=./async_sms_repository.go -package=repomocks -destination=mocks/async_sms_repository.mock.go AsyncSmsRepository
//

// Package repomocks is a generated GoMock package.
package repomocks

import (
	context "context"
	reflect "reflect"

	domain "gitee.com/geekbang/basic-go/webook/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockAsyncSmsRepository is a mock of AsyncSmsRepository interface.
type MockAsyncSmsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAsyncSmsRepositoryMockRecorder
}

// MockAsyncSmsRepositoryMockRecorder is the mock recorder for MockAsyncSmsRepository.
type MockAsyncSmsRepositoryMockRecorder struct {
	mock *MockAsyncSmsRepository
}

// NewMockAsyncSmsRepository creates a new mock instance.
func NewMockAsyncSmsRepository(ctrl *gomock.Controller) *MockAsyncSmsRepository {
	mock := &MockAsyncSmsRepository{ctrl: ctrl}
	mock.recorder = &MockAsyncSmsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAsyncSmsRepository) EXPECT() *MockAsyncSmsRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockAsyncSmsRepository) Add(ctx context.Context, s domain.AsyncSms) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, s)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockAsyncSmsRepositoryMockRecorder) Add(ctx, s any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockAsyncSmsRepository)(nil).Add), ctx, s)
}

// PreemptWaitingSMS mocks base method.
func (m *MockAsyncSmsRepository) PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreemptWaitingSMS", ctx)
	ret0, _ := ret[0].(domain.AsyncSms)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PreemptWaitingSMS indicates an expected call of PreemptWaitingSMS.
func (mr *MockAsyncSmsRepositoryMockRecorder) PreemptWaitingSMS(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreemptWaitingSMS", reflect.TypeOf((*MockAsyncSmsRepository)(nil).PreemptWaitingSMS), ctx)
}

// ReportScheduleResult mocks base method.
func (m *MockAsyncSmsRepository) ReportScheduleResult(ctx context.Context, id int64, success bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReportScheduleResult", ctx, id, success)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportScheduleResult indicates an expected call of ReportScheduleResult.
func (mr *MockAsyncSmsRepositoryMockRecorder) ReportScheduleResult(ctx, id, success any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportScheduleResult", reflect.TypeOf((*MockAsyncSmsRepository)(nil).ReportScheduleResult), ctx, id, success)
}
