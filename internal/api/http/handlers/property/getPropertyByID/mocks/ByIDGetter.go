// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/alimzhanoff/property-finder/internal/models"
)

// ByIDGetter is an autogenerated mock type for the ByIDGetter type
type ByIDGetter struct {
	mock.Mock
}

// GetPropertyByID provides a mock function with given fields: ctx, id
func (_m *ByIDGetter) GetPropertyByID(ctx context.Context, id int) (models.Property, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetPropertyByID")
	}

	var r0 models.Property
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (models.Property, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) models.Property); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.Property)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewByIDGetter creates a new instance of ByIDGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewByIDGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *ByIDGetter {
	mock := &ByIDGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
