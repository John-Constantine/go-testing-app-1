// Code generated by MockGen. DO NOT EDIT.
// Source: ./db/repository/field_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "flow/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// FetchFieldFromFieldVersion mocks base method
func (m *MockRepository) FetchFieldFromFieldVersion(completeFieldVersionNumberList map[int]bool) []model.FieldVersion {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchFieldFromFieldVersion", completeFieldVersionNumberList)
	ret0, _ := ret[0].([]model.FieldVersion)
	return ret0
}

// FetchFieldFromFieldVersion indicates an expected call of FetchFieldFromFieldVersion
func (mr *MockRepositoryMockRecorder) FetchFieldFromFieldVersion(completeFieldVersionNumberList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchFieldFromFieldVersion", reflect.TypeOf((*MockRepository)(nil).FetchFieldFromFieldVersion), completeFieldVersionNumberList)
}

// FindByExternalId mocks base method
func (m *MockRepository) FindByExternalId(flowExternalId string) model.Flow {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByExternalId", flowExternalId)
	ret0, _ := ret[0].(model.Flow)
	return ret0
}

// FindByExternalId indicates an expected call of FindByExternalId
func (mr *MockRepositoryMockRecorder) FindByExternalId(flowExternalId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByExternalId", reflect.TypeOf((*MockRepository)(nil).FindByExternalId), flowExternalId)
}
