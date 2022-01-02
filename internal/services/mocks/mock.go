// Code generated by MockGen. DO NOT EDIT.
// Source: services.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTokener is a mock of Tokener interface.
type MockTokener struct {
	ctrl     *gomock.Controller
	recorder *MockTokenerMockRecorder
}

// MockTokenerMockRecorder is the mock recorder for MockTokener.
type MockTokenerMockRecorder struct {
	mock *MockTokener
}

// NewMockTokener creates a new mock instance.
func NewMockTokener(ctrl *gomock.Controller) *MockTokener {
	mock := &MockTokener{ctrl: ctrl}
	mock.recorder = &MockTokenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokener) EXPECT() *MockTokenerMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockTokener) GenerateToken(userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockTokenerMockRecorder) GenerateToken(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockTokener)(nil).GenerateToken), userId)
}

// ParseToken mocks base method.
func (m *MockTokener) ParseToken(token string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockTokenerMockRecorder) ParseToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockTokener)(nil).ParseToken), token)
}

// MockCloudStorage is a mock of CloudStorage interface.
type MockCloudStorage struct {
	ctrl     *gomock.Controller
	recorder *MockCloudStorageMockRecorder
}

// MockCloudStorageMockRecorder is the mock recorder for MockCloudStorage.
type MockCloudStorageMockRecorder struct {
	mock *MockCloudStorage
}

// NewMockCloudStorage creates a new mock instance.
func NewMockCloudStorage(ctrl *gomock.Controller) *MockCloudStorage {
	mock := &MockCloudStorage{ctrl: ctrl}
	mock.recorder = &MockCloudStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCloudStorage) EXPECT() *MockCloudStorageMockRecorder {
	return m.recorder
}

// UploadFile mocks base method.
func (m *MockCloudStorage) UploadFile(file []byte, filesize int64, filename string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", file, filesize, filename)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockCloudStorageMockRecorder) UploadFile(file, filesize, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockCloudStorage)(nil).UploadFile), file, filesize, filename)
}