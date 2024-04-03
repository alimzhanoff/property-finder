// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/api/http/handlers/property/saveProperty/handler.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/alimzhanoff/property-finder/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockPropertySaver is a mock of PropertySaver interface.
type MockPropertySaver struct {
	ctrl     *gomock.Controller
	recorder *MockPropertySaverMockRecorder
}

// MockPropertySaverMockRecorder is the mock recorder for MockPropertySaver.
type MockPropertySaverMockRecorder struct {
	mock *MockPropertySaver
}

// NewMockPropertySaver creates a new mock instance.
func NewMockPropertySaver(ctrl *gomock.Controller) *MockPropertySaver {
	mock := &MockPropertySaver{ctrl: ctrl}
	mock.recorder = &MockPropertySaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPropertySaver) EXPECT() *MockPropertySaverMockRecorder {
	return m.recorder
}

// SavePropertyWithPropertyTypeWithTx mocks base method.
func (m *MockPropertySaver) SavePropertyWithPropertyTypeWithTx(ctx context.Context, property models.Property) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePropertyWithPropertyTypeWithTx", ctx, property)
	ret0, _ := ret[0].(error)
	return ret0
}

// SavePropertyWithPropertyTypeWithTx indicates an expected call of SavePropertyWithPropertyTypeWithTx.
func (mr *MockPropertySaverMockRecorder) SavePropertyWithPropertyTypeWithTx(ctx, property interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePropertyWithPropertyTypeWithTx", reflect.TypeOf((*MockPropertySaver)(nil).SavePropertyWithPropertyTypeWithTx), ctx, property)
}