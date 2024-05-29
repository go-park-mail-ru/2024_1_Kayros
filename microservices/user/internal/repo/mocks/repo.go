// Code generated by MockGen. DO NOT EDIT.
// Source: repo.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	user "2024_1_kayros/gen/go/user"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepo) Create(ctx context.Context, u *user.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepoMockRecorder) Create(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepo)(nil).Create), ctx, u)
}

// CreateAddressByUnauthId mocks base method.
func (m *MockRepo) CreateAddressByUnauthId(ctx context.Context, data *user.AddressDataUnauth) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAddressByUnauthId", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAddressByUnauthId indicates an expected call of CreateAddressByUnauthId.
func (mr *MockRepoMockRecorder) CreateAddressByUnauthId(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAddressByUnauthId", reflect.TypeOf((*MockRepo)(nil).CreateAddressByUnauthId), ctx, data)
}

// DeleteByEmail mocks base method.
func (m *MockRepo) DeleteByEmail(ctx context.Context, email *user.Email) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByEmail", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByEmail indicates an expected call of DeleteByEmail.
func (mr *MockRepoMockRecorder) DeleteByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByEmail", reflect.TypeOf((*MockRepo)(nil).DeleteByEmail), ctx, email)
}

// GetAddressByUnauthId mocks base method.
func (m *MockRepo) GetAddressByUnauthId(ctx context.Context, id *user.UnauthId) (*user.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddressByUnauthId", ctx, id)
	ret0, _ := ret[0].(*user.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAddressByUnauthId indicates an expected call of GetAddressByUnauthId.
func (mr *MockRepoMockRecorder) GetAddressByUnauthId(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddressByUnauthId", reflect.TypeOf((*MockRepo)(nil).GetAddressByUnauthId), ctx, id)
}

// GetByEmail mocks base method.
func (m *MockRepo) GetByEmail(ctx context.Context, email *user.Email) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", ctx, email)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockRepoMockRecorder) GetByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockRepo)(nil).GetByEmail), ctx, email)
}

// Update mocks base method.
func (m *MockRepo) Update(ctx context.Context, data *user.UpdateUserData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepoMockRecorder) Update(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepo)(nil).Update), ctx, data)
}

// UpdateAddressByUnauthId mocks base method.
func (m *MockRepo) UpdateAddressByUnauthId(ctx context.Context, data *user.AddressDataUnauth) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAddressByUnauthId", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAddressByUnauthId indicates an expected call of UpdateAddressByUnauthId.
func (mr *MockRepoMockRecorder) UpdateAddressByUnauthId(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAddressByUnauthId", reflect.TypeOf((*MockRepo)(nil).UpdateAddressByUnauthId), ctx, data)
}
