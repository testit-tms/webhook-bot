// Code generated by MockGen. DO NOT EDIT.
// Source: sendmessage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entities "github.com/testit-tms/webhook-bot/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockchatGeter is a mock of chatGeter interface.
type MockchatGeter struct {
	ctrl     *gomock.Controller
	recorder *MockchatGeterMockRecorder
}

// MockchatGeterMockRecorder is the mock recorder for MockchatGeter.
type MockchatGeterMockRecorder struct {
	mock *MockchatGeter
}

// NewMockchatGeter creates a new mock instance.
func NewMockchatGeter(ctrl *gomock.Controller) *MockchatGeter {
	mock := &MockchatGeter{ctrl: ctrl}
	mock.recorder = &MockchatGeterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockchatGeter) EXPECT() *MockchatGeterMockRecorder {
	return m.recorder
}

// GetChatsByCompanyToken mocks base method.
func (m *MockchatGeter) GetChatsByCompanyToken(ctx context.Context, t string) ([]entities.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChatsByCompanyToken", ctx, t)
	ret0, _ := ret[0].([]entities.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChatsByCompanyToken indicates an expected call of GetChatsByCompanyToken.
func (mr *MockchatGeterMockRecorder) GetChatsByCompanyToken(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChatsByCompanyToken", reflect.TypeOf((*MockchatGeter)(nil).GetChatsByCompanyToken), ctx, t)
}

// MockbotSender is a mock of botSender interface.
type MockbotSender struct {
	ctrl     *gomock.Controller
	recorder *MockbotSenderMockRecorder
}

// MockbotSenderMockRecorder is the mock recorder for MockbotSender.
type MockbotSenderMockRecorder struct {
	mock *MockbotSender
}

// NewMockbotSender creates a new mock instance.
func NewMockbotSender(ctrl *gomock.Controller) *MockbotSender {
	mock := &MockbotSender{ctrl: ctrl}
	mock.recorder = &MockbotSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockbotSender) EXPECT() *MockbotSenderMockRecorder {
	return m.recorder
}

// SendMessage mocks base method.
func (m *MockbotSender) SendMessage(ctx context.Context, msg entities.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", ctx, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockbotSenderMockRecorder) SendMessage(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockbotSender)(nil).SendMessage), ctx, msg)
}
