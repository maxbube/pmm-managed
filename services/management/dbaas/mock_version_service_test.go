// Code generated by mockery v1.0.0. DO NOT EDIT.

package dbaas

import (
	context "context"

	version "github.com/hashicorp/go-version"
	mock "github.com/stretchr/testify/mock"
)

// mockVersionService is an autogenerated mock type for the versionService type
type mockVersionService struct {
	mock.Mock
}

// GetLatestOperatorVersion provides a mock function with given fields: _a0, _a1
func (_m *mockVersionService) GetLatestOperatorVersion(_a0 context.Context, _a1 string) (*version.Version, *version.Version, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *version.Version
	if rf, ok := ret.Get(0).(func(context.Context, string) *version.Version); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*version.Version)
		}
	}

	var r1 *version.Version
	if rf, ok := ret.Get(1).(func(context.Context, string) *version.Version); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*version.Version)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// IsOperatorVersionSupported provides a mock function with given fields: ctx, operatorType, pmmVersion, operatorVersion
func (_m *mockVersionService) IsOperatorVersionSupported(ctx context.Context, operatorType string, pmmVersion string, operatorVersion string) (bool, error) {
	ret := _m.Called(ctx, operatorType, pmmVersion, operatorVersion)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) bool); ok {
		r0 = rf(ctx, operatorType, pmmVersion, operatorVersion)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, operatorType, pmmVersion, operatorVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Matrix provides a mock function with given fields: ctx, params
func (_m *mockVersionService) Matrix(ctx context.Context, params componentsParams) (*VersionServiceResponse, error) {
	ret := _m.Called(ctx, params)

	var r0 *VersionServiceResponse
	if rf, ok := ret.Get(0).(func(context.Context, componentsParams) *VersionServiceResponse); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*VersionServiceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, componentsParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
