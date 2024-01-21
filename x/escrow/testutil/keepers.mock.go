// Code generated by MockGen. DO NOT EDIT.
// Source: keeper/expected/keepers.go

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	expected "github.com/0tech/andromeda/x/escrow/keeper/expected"
	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

// MockMessageRouter is a mock of MessageRouter interface.
type MockMessageRouter struct {
	ctrl     *gomock.Controller
	recorder *MockMessageRouterMockRecorder
}

// MockMessageRouterMockRecorder is the mock recorder for MockMessageRouter.
type MockMessageRouterMockRecorder struct {
	mock *MockMessageRouter
}

// NewMockMessageRouter creates a new mock instance.
func NewMockMessageRouter(ctrl *gomock.Controller) *MockMessageRouter {
	mock := &MockMessageRouter{ctrl: ctrl}
	mock.recorder = &MockMessageRouterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageRouter) EXPECT() *MockMessageRouterMockRecorder {
	return m.recorder
}

// Handler mocks base method.
func (m *MockMessageRouter) Handler(msg types.Msg) expected.MsgServiceHandler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handler", msg)
	ret0, _ := ret[0].(expected.MsgServiceHandler)
	return ret0
}

// Handler indicates an expected call of Handler.
func (mr *MockMessageRouterMockRecorder) Handler(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handler", reflect.TypeOf((*MockMessageRouter)(nil).Handler), msg)
}

// MockAuthKeeper is a mock of AuthKeeper interface.
type MockAuthKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAuthKeeperMockRecorder
}

// MockAuthKeeperMockRecorder is the mock recorder for MockAuthKeeper.
type MockAuthKeeperMockRecorder struct {
	mock *MockAuthKeeper
}

// NewMockAuthKeeper creates a new mock instance.
func NewMockAuthKeeper(ctrl *gomock.Controller) *MockAuthKeeper {
	mock := &MockAuthKeeper{ctrl: ctrl}
	mock.recorder = &MockAuthKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthKeeper) EXPECT() *MockAuthKeeperMockRecorder {
	return m.recorder
}

// HasAccount mocks base method.
func (m *MockAuthKeeper) HasAccount(arg0 context.Context, arg1 types.AccAddress) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasAccount", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasAccount indicates an expected call of HasAccount.
func (mr *MockAuthKeeperMockRecorder) HasAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasAccount", reflect.TypeOf((*MockAuthKeeper)(nil).HasAccount), arg0, arg1)
}

// NewAccount mocks base method.
func (m *MockAuthKeeper) NewAccount(arg0 context.Context, arg1 types.AccountI) types.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAccount", arg0, arg1)
	ret0, _ := ret[0].(types.AccountI)
	return ret0
}

// NewAccount indicates an expected call of NewAccount.
func (mr *MockAuthKeeperMockRecorder) NewAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAccount", reflect.TypeOf((*MockAuthKeeper)(nil).NewAccount), arg0, arg1)
}

// SetAccount mocks base method.
func (m *MockAuthKeeper) SetAccount(arg0 context.Context, arg1 types.AccountI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAccount", arg0, arg1)
}

// SetAccount indicates an expected call of SetAccount.
func (mr *MockAuthKeeperMockRecorder) SetAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAccount", reflect.TypeOf((*MockAuthKeeper)(nil).SetAccount), arg0, arg1)
}