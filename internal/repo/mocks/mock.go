// Code generated by MockGen. DO NOT EDIT.
// Source: repo.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	models "creatly-task/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUsers is a mock of Users interface.
type MockUsers struct {
	ctrl     *gomock.Controller
	recorder *MockUsersMockRecorder
}

// MockUsersMockRecorder is the mock recorder for MockUsers.
type MockUsersMockRecorder struct {
	mock *MockUsers
}

// NewMockUsers creates a new mock instance.
func NewMockUsers(ctrl *gomock.Controller) *MockUsers {
	mock := &MockUsers{ctrl: ctrl}
	mock.recorder = &MockUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsers) EXPECT() *MockUsersMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUsers) CreateUser(arg0 *models.UserSignUpInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUsersMockRecorder) CreateUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUsers)(nil).CreateUser), arg0)
}

// GetUserByCreds mocks base method.
func (m *MockUsers) GetUserByCreds(email string) (*models.UserSignInOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByCreds", email)
	ret0, _ := ret[0].(*models.UserSignInOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByCreds indicates an expected call of GetUserByCreds.
func (mr *MockUsersMockRecorder) GetUserByCreds(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByCreds", reflect.TypeOf((*MockUsers)(nil).GetUserByCreds), email)
}

// MockTokens is a mock of Tokens interface.
type MockTokens struct {
	ctrl     *gomock.Controller
	recorder *MockTokensMockRecorder
}

// MockTokensMockRecorder is the mock recorder for MockTokens.
type MockTokensMockRecorder struct {
	mock *MockTokens
}

// NewMockTokens creates a new mock instance.
func NewMockTokens(ctrl *gomock.Controller) *MockTokens {
	mock := &MockTokens{ctrl: ctrl}
	mock.recorder = &MockTokensMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokens) EXPECT() *MockTokensMockRecorder {
	return m.recorder
}

// GetUserIDByToken mocks base method.
func (m *MockTokens) GetUserIDByToken(token string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDByToken", token)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDByToken indicates an expected call of GetUserIDByToken.
func (mr *MockTokensMockRecorder) GetUserIDByToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDByToken", reflect.TypeOf((*MockTokens)(nil).GetUserIDByToken), token)
}

// MockFiles is a mock of Files interface.
type MockFiles struct {
	ctrl     *gomock.Controller
	recorder *MockFilesMockRecorder
}

// MockFilesMockRecorder is the mock recorder for MockFiles.
type MockFilesMockRecorder struct {
	mock *MockFiles
}

// NewMockFiles creates a new mock instance.
func NewMockFiles(ctrl *gomock.Controller) *MockFiles {
	mock := &MockFiles{ctrl: ctrl}
	mock.recorder = &MockFilesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFiles) EXPECT() *MockFilesMockRecorder {
	return m.recorder
}

// AddLog mocks base method.
func (m *MockFiles) AddLog(log *models.FileUploadLogInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLog", log)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLog indicates an expected call of AddLog.
func (mr *MockFilesMockRecorder) AddLog(log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLog", reflect.TypeOf((*MockFiles)(nil).AddLog), log)
}

// All mocks base method.
func (m *MockFiles) All() ([]models.FileOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All")
	ret0, _ := ret[0].([]models.FileOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockFilesMockRecorder) All() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockFiles)(nil).All))
}
